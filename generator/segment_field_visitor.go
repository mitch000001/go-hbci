package generator

import (
	"fmt"
	"go/ast"
	"go/token"
)

const segmentElementsMethodIdentifier = "elements"
const segmentElementsMethodReturnType = "[]element.DataElement"

type elementVisitor struct {
	receiverName string
	fileSet      *token.FileSet
	object       *ast.Object
	fields       []field
	err          error
}

func (e *elementVisitor) Visit(node ast.Node) ast.Visitor {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		if e.isElementsMethod(funcDecl) {
			bodyStatements := funcDecl.Body.List
			if len(bodyStatements) == 1 {
				if ret, ok := bodyStatements[0].(*ast.ReturnStmt); ok {
					if len(ret.Results) == 1 {
						if res, ok := ret.Results[0].(*ast.CompositeLit); ok {
							resType := nodeToString(res.Type, e.fileSet)
							if resType == segmentElementsMethodReturnType {
								for _, element := range res.Elts {
									if sel, ok := element.(*ast.SelectorExpr); ok {
										pos := sel.Pos()
										position := e.fileSet.Position(pos)
										elemField := field{Name: sel.Sel.Name, Line: position.Line}
										if _, ok := sel.X.(*ast.Ident); ok {
											fieldVisitor := &structFieldVisitor{fileSet: e.fileSet, fieldName: sel.Sel.Name}
											ast.Walk(fieldVisitor, e.object.Decl.(*ast.TypeSpec))
											if fieldVisitor.err != nil {
												e.err = fieldVisitor.err
												return nil
											}
											elemField.TypeDecl = fieldVisitor.fieldTypeDecl
										} else {
											e.err = fmt.Errorf("Unexpected Selector: %q", nodeToString(sel, e.fileSet))
											return nil
										}
										if elemField.TypeDecl == "" {
											e.err = fmt.Errorf("No type declaration found for element %q", sel.Sel.Name)
											return nil
										}
										e.fields = append(e.fields, elemField)
									} else {
										e.err = fmt.Errorf("Unsupported element in elements method: %q", nodeToString(element, e.fileSet))
										return nil
									}
								}
								return nil
							}
						}
					}
				}
			}
		}
	}
	return e
}

func (e *elementVisitor) isElementsMethod(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Name.Name == segmentElementsMethodIdentifier {
		if funcDecl.Recv != nil {
			if fields := funcDecl.Recv.List; fields != nil {
				if len(fields) == 1 {
					if typ, ok := fields[0].Type.(*ast.StarExpr); ok {
						if ident, ok := typ.X.(*ast.Ident); ok {
							if ident.Name == e.receiverName {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

type structFieldVisitor struct {
	fileSet       *token.FileSet
	fieldName     string
	fieldTypeDecl string
	err           error
}

func (s *structFieldVisitor) Visit(node ast.Node) ast.Visitor {
	if structType, ok := node.(*ast.StructType); ok {
		if fields := structType.Fields; fields != nil {
			for _, f := range fields.List {
				var fieldName string
				if names := f.Names; names != nil {
					fieldName = nodeToString(names[0], s.fileSet)
					if fieldName != s.fieldName {
						continue // not the field we're looking for
					}
				} else {
					continue // anonymous field
				}
				if starExpr, ok := f.Type.(*ast.StarExpr); ok {
					s.fieldTypeDecl = nodeToString(starExpr.X, s.fileSet)
				} else {
					s.err = fmt.Errorf("Unexpected type found: %q", nodeToString(f.Type, s.fileSet))
					return nil
				}
			}
			return nil
		}
	}
	return s
}

package crypto

import (
	"fmt"
	"math/big"
)

type RSAKeyParameters interface {
	Modulus() *big.Int
	Exponent() *big.Int
	Private() bool
}

func NewRSAKeyParameters(isPrivate bool, modulus, exponent *big.Int) RSAKeyParameters {
	return &rsaKeyParameters{
		privateKey: isPrivate,
		modulus:    modulus,
		exponent:   exponent,
	}
}

type rsaKeyParameters struct {
	privateKey bool
	modulus    *big.Int
	exponent   *big.Int
}

func (r *rsaKeyParameters) Modulus() *big.Int  { return r.modulus }
func (r *rsaKeyParameters) Exponent() *big.Int { return r.exponent }
func (r *rsaKeyParameters) Private() bool      { return r.privateKey }

func NewRSAPrivateCrtKeyParameters(
	modulus,
	publicExponent,
	privateExponent,
	p,
	q,
	dP,
	dQ,
	qInv *big.Int) *RSAPrivateCrtKeyParameters {
	return &RSAPrivateCrtKeyParameters{
		RSAKeyParameters: NewRSAKeyParameters(true, modulus, privateExponent),
		e:                publicExponent,
		p:                p,
		q:                q,
		dP:               dP,
		dQ:               dQ,
		qInv:             qInv,
	}
}

type RSAPrivateCrtKeyParameters struct {
	RSAKeyParameters
	e    *big.Int
	p    *big.Int
	q    *big.Int
	dP   *big.Int
	dQ   *big.Int
	qInv *big.Int
}

func (r *RSAPrivateCrtKeyParameters) PublicExponent() *big.Int { return r.e }
func (r *RSAPrivateCrtKeyParameters) P() *big.Int              { return r.p }
func (r *RSAPrivateCrtKeyParameters) Q() *big.Int              { return r.q }
func (r *RSAPrivateCrtKeyParameters) DP() *big.Int             { return r.dP }
func (r *RSAPrivateCrtKeyParameters) DQ() *big.Int             { return r.dQ }
func (r *RSAPrivateCrtKeyParameters) QInv() *big.Int           { return r.qInv }

type RSAEngine struct {
	*RSACoreEngine
}

func (r *RSAEngine) Init(forEncryption bool, key RSAKeyParameters) {
	if r.RSACoreEngine == nil {
		r.RSACoreEngine = newRsaCoreEngine()
	}
	r.RSACoreEngine.Init(forEncryption, key)
}

func (r *RSAEngine) ProcessBlock(in []byte, inOff, inLen int) []byte {
	if r.RSACoreEngine == nil {
		panic(fmt.Errorf("RAS engine not initialized"))
	}
	return r.RSACoreEngine.ConvertOutput(r.RSACoreEngine.ProcessBlock(r.RSACoreEngine.ConvertInput(in, inOff, inLen)))
}

func newRsaCoreEngine() *RSACoreEngine {
	return &RSACoreEngine{}
}

type RSACoreEngine struct {
	key           RSAKeyParameters
	forEncryption bool
}

func (r *RSACoreEngine) Init(forEncryption bool, key RSAKeyParameters) {
	r.key = key
	r.forEncryption = forEncryption
}

func (r *RSACoreEngine) InputBlockSize() int {
	bitSize := r.key.Modulus().BitLen()
	if r.forEncryption {
		return (bitSize+7)/8 - 1
	} else {
		return (bitSize + 7) / 8
	}
}

func (r *RSACoreEngine) OutputBlockSize() int {
	bitSize := r.key.Modulus().BitLen()
	if r.forEncryption {
		return (bitSize + 7) / 8
	} else {
		return (bitSize+7)/8 - 1
	}
}

func (r *RSACoreEngine) ConvertInput(in []byte, inOff, inLen int) *big.Int {
	if inLen > (r.InputBlockSize() + 1) {
		panic(fmt.Errorf("input too large for RSA cipher."))
	} else if inLen == (r.InputBlockSize()+1) && !r.forEncryption {
		panic(fmt.Errorf("input too large for RSA cipher."))
	}

	var block []byte

	if inOff != 0 || inLen != len(in) {
		block = make([]byte, inLen)
		// System.arraycopy(in, inOff, block, 0, inLen);
		arrayCopy(in, inOff, block, 0, inLen)
	} else {
		block = in
	}
	res := new(big.Int)
	res = res.Abs(res.SetBytes(block))
	if res.Cmp(r.key.Modulus()) >= 0 {
		panic(fmt.Errorf("input too large for RSA cipher."))
	}
	return res
}

func (r *RSACoreEngine) ConvertOutput(result *big.Int) []byte {
	output := result.Bytes()
	if r.forEncryption {
		if output[0] == 0 && len(output) > r.OutputBlockSize() {

			// have ended up with an extra zero byte, copy down.
			tmp := make([]byte, len(output)-1)

			// System.arraycopy(output, 1, tmp, 0, tmp.length);
			arrayCopy(output, 1, tmp, 0, len(tmp))

			return tmp
		}

		if len(output) < r.OutputBlockSize() { // have ended up with less bytes than normal, lengthen
			tmp := make([]byte, r.OutputBlockSize())

			// System.arraycopy(output, 0, tmp, tmp.length - output.length, output.length);
			arrayCopy(output, 0, tmp, len(tmp)-len(output), len(output))

			return tmp
		}

	} else {
		if output[0] == 0 { // have ended up with an extra zero byte, copy down.
			tmp := make([]byte, len(output)-1)

			// System.arraycopy(output, 1, tmp, 0, tmp.length);
			arrayCopy(output, 1, tmp, 0, len(tmp))

			return tmp
		}
	}
	return output
}

func (r *RSACoreEngine) ProcessBlock(input *big.Int) *big.Int {
	if crtKey, ok := r.key.(*RSAPrivateCrtKeyParameters); ok {
		//
		// we have the extra factors, use the Chinese Remainder Theorem - the author
		// wishes to express his thanks to Dirk Bonekaemper at rtsffm.com for
		// advice regarding the expression of this.
		//

		p := crtKey.P()
		q := crtKey.Q()
		dP := crtKey.DP()
		dQ := crtKey.DQ()
		qInv := crtKey.QInv()

		var mP, mQ, h, m *big.Int

		// mP = ((input mod p) ^ dP)) mod p
		mP = mP.Exp(new(big.Int).Rem(input, p), dP, p)

		// mQ = ((input mod q) ^ dQ)) mod q
		mQ = mQ.Exp(new(big.Int).Rem(input, q), dQ, q)

		// h = qInv * (mP - mQ) mod p
		h = h.Sub(mP, mQ)
		h = h.Mul(h, qInv)
		h = h.Mod(h, p) // mod (in Java) returns the positive residual

		// m = h * q + mQ
		m = h.Mul(h, q)
		m = m.Add(m, mQ)

		return m
	} else {
		return new(big.Int).Exp(input, r.key.Exponent(), r.key.Modulus())
	}
}

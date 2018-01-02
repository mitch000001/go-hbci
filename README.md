# go-hbci
[![Build Status](https://travis-ci.org/mitch000001/go-hbci.svg)](https://travis-ci.org/mitch000001/go-hbci)
[![License: Apache v2.0](https://badge.luzifer.io/v1/badge?color=5d79b5&title=license&text=Apache+v2.0)](http://www.apache.org/licenses/LICENSE-2.0)
[![GoDoc](https://godoc.org/github.com/mitch000001/go-hbci?status.svg)](http://godoc.org/github.com/mitch000001/go-hbci)

A client library to use the [Home Banking Computer Interface](https://de.wikipedia.org/wiki/Homebanking_Computer_Interface) (german only)

For an exhausted reference of HBCI visit the website of [The German Banking Industry](https://www.hbci-zka.de/)

## Status of the project
Due to the massive amount of the standard this library is only at the beginning of being useful to use.
Also, there is no client interface yet in terms of entry point for the library or management of pin/tan or any other data.

The implemented standard conforms to HBCI 2.2 and FINTS 3.0.

## Roadmap
- [x] Parsing Accounts
- [x] Listing transactions
- [ ] Some other read only action
- [ ] Write access

# go-hbci
[![Build Status](https://travis-ci.org/mitch000001/go-hbci.svg)](https://travis-ci.org/mitch000001/go-hbci)
[![Coverage Status](https://coveralls.io/repos/mitch000001/go-hbci/badge.svg?branch=master&service=github)](https://coveralls.io/github/mitch000001/go-hbci?branch=master)

A client library to use the [Home Banking Computer Interface](https://de.wikipedia.org/wiki/Homebanking_Computer_Interface) (german only)

For an exhausted reference of HBCI visit the website of [The German Banking Industry](https://www.hbci-zka.de/)

## Status of the project
Due to the massive amount of the standard this library is only at the beginning of being useful to use.
Also, there is no client interface yet in terms of entry point for the library or management of pin/tan or any other data.

The implemented standard conforms to HBCI 2.2.

## Roadmap
- [x] Parsing Accounts
- [x] Listing transactions
- [ ] Some other read only action
- [ ] Write access

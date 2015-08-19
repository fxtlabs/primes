# primes

[![Build Status](https://api.travis-ci.org/fxtlabs/primes.svg?branch=master)](https://travis-ci.org/fxtlabs/primes)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/fxtlabs/primes)

Package `primes` provides simple functionality for working with prime numbers.

Call `Sieve(n)` to generate all prime numbers less than or equal to n,
`IsPrime(n)` to test for primality, `Coprime(a,b)` to test for coprimality,
and `Pi(n)` to count (or estimate) the number of primes less than or equal to n.

The algorithms used to implement the functions above are fairly simple;
they work well with relatively small primes, but they are definitely not
intended for work in cryptography or any application requiring really
large primes.

See [package documentation](https://godoc.org/github.com/fxtlabs/primes) for
full documentation and examples.

## Installation

    go get -u github.com/fxtlabs/primes

## License

The MIT License (MIT). See the LICENSE files in this directory.

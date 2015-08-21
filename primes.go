// The MIT License (MIT)
//
// Copyright (c) 2015 Filippo Tampieri
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package primes provides simple functionality for working with prime numbers.
//
// Call Sieve(n) to generate all prime numbers less than or equal to n,
// IsPrime(n) to test for primality, Coprime(a,b) to test for coprimality,
// and Pi(n) to count (or estimate) the number of primes less than or equal to n.
//
// The algorithms used to implement the functions above are fairly simple;
// they work well with relatively small primes, but they are definitely not
// intended for work in cryptography or any application requiring really
// large primes.  Run the benchmarks to check their performance against
// simpler baseline implementations.
//
package primes

import (
	"math"
	"sort"
)

// primes is a cache of the first few prime numbers
var primes []int

func init() {
	// Cache the first 1,229 prime numbers (i.e. all primes <= 10,000)
	primes = Sieve(10000)
}

// Pi returns the number of primes less than or equal to n.
// If ok is true, the result is correct; otherwise it is an estimate
// computed as n/(log(n)-1) with an error generally below 1% (see tests).
// This estimate is based on the prime number theorem which states that the
// number of primes not exceeding n is asymptotic to n/log(n).
// Better approximations are known, but they are more complex.
// See https://primes.utm.edu/howmany.html#1,
// https://en.wikipedia.org/wiki/Prime_number_theorem, and
// https://en.wikipedia.org/wiki/Prime-counting_function for details.
func Pi(n int) (pi int, ok bool) {
	// If n is smaller than or equal to the largest cached prime,
	// we have an exact count
	if i := sort.SearchInts(primes, n); i < len(primes) {
		if n == primes[i] {
			// n is the prime at index i
			pi = i + 1
		} else {
			// n is not prime and primes[j] < n for all j in [0,i)
			pi = i
		}
		ok = true
	} else {
		// n is larger than the largest prime in the cache;
		// use the estimate
		pi = int(float64(n) / (math.Log(float64(n)) - 1))
	}
	return
}

// IsPrime is a primality test: it returns true if n is prime.
// It uses trial division to check whether n can be divided exactly by any
// number that is less than n.
// The algorithm can be very slow on large n despite a number of optimizations:
//
// * If n is relatively small, it is compared against a cached table of primes.
//
// * Only numbers up to sqrt(n) are used to check for primality.
//
// * n is first checked for divisibility by the primes in the cache and only if the test is inconclusive, n is checked against more numbers.
//
// * Only numbers of the form 6*k+|-1 that are greater than the last prime in the cache are checked after that.
//
// See https://en.wikipedia.org/wiki/Primality_test and
// https://en.wikipedia.org/wiki/Trial_division for details.
func IsPrime(n int) bool {
	pMax := primes[len(primes)-1]
	if n <= pMax {
		// If n is prime, it must be in the cache
		i := sort.SearchInts(primes, n)
		return n == primes[i]
	}
	max := int(math.Ceil(math.Sqrt(float64(n))))
	// Check if n is divisible by any of the cached primes
	for _, p := range primes {
		if p > max {
			return true
		}
		if n%p == 0 {
			return false
		}
	}
	// When you run out of cached primes, check if n is divisible by
	// any number 6*k+|-1 larger than the largest prime in the cache.
	for d := (pMax/6+1)*6 - 1; d <= max; d += 6 {
		if n%d == 0 || n%(d+2) == 0 {
			return false
		}
	}
	return true
}

// Coprime is a coprimality test: it returns true if the only positive integer
// that divides evenly both a and b is 1.
// This function implements the division-based version of the Euclidean algorithm.
// See https://en.wikipedia.org/wiki/Coprime_integers and
// https://en.wikipedia.org/wiki/Euclidean_algorithm for details.
func Coprime(a, b int) bool {
	// Set a to gcd(a,b)
	var t int
	for b != 0 {
		t = b
		b = a % b
		a = t
	}
	// By definition, a and b are coprime if gcd(a,b) == 1
	return a == 1
}

// Sieve returns a list of the prime numbers less than or equal to n.
// If n is less than 2, it returns an empty list.
// The function uses the sieve of Eratosthenes algorithm
// (see https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes)
// with the following optimizations:
//
// * The initial list of candidate primes includes odd numbers only.
//
// * Given a prime p, only multiples of p greater than or equal to p*p need to be marked off since smaller multiples of p have already been marked off by then.
//
// * The above also implies that the algorithm can terminate as soon as it finds  a prime p such that p*p is greater than n.
//
// Sieve takes O(n) memory and runs in O(n log log n) time.
func Sieve(n int) []int {
	switch {
	case n < 2:
		return []int{}
	case n == 2:
		return []int{2}
	}
	// a[i] == false ==> p=2*i+3 is a candidate prime
	// p in [3,n] ==> i in [0,(n-3)/2]
	length := 1 + (n-3)/2
	a := make([]bool, length, length)
	// Start with number 3 and consider only odd numbers
	sqrtn := int(math.Sqrt(float64(n)))
	for i, p := 0, 3; p <= sqrtn; p += 2 {
		if !a[i] {
			// 2*i+1 is a prime number; mark off its multiples
			for j := (p*p - 3) / 2; j < length; j += p {
				a[j] = true
			}
		}
		i++
	}
	// ps will store the computed primes; its initial capacity is based
	// an estimate of the prime-counting function pi(n)
	pi, _ := Pi(n)
	ps := make([]int, 1, pi)
	ps[0] = 2
	for i := 0; i < length; i++ {
		if !a[i] {
			ps = append(ps, 2*i+3)
		}
	}
	return ps
}

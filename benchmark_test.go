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

package primes_test

import (
	"math"
	"testing"
	"testing/quick"

	"github.com/fxtlabs/primes"
)

// baselineIsPrime returns true if n is prime.
// It uses trial division against all i in [2,sqrt(n)] and will be extremely
// slow for large n.
// Used for testing only.
func baselineIsPrime(n int) bool {
	switch {
	case n < 2:
		return false
	case n == 2:
		return true
	}
	max := int(math.Ceil(math.Sqrt(float64(n))))
	for d := 2; d <= max; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

// baselineSieve returns a list of the prime numbers less than or equal to n.
// If n is less than 2, it returns an empty list.
// The function uses the sieve of Eratosthenes algorithm with the following
// simple optimizations:
//
// * Given a prime p, only multiples of p greater than or equal to p*p need to be marked off since smaller multiples of p have already been marked off by then.
//
// * The above also implies that the algorithm can terminate as soon as it finds  a prime p such that p*p is greater than n.
//
// See https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes for details.
func baselineSieve(n int) []int {
	if n < 2 {
		return []int{}
	}
	// a[i] == false ==> i is a candidate prime
	a := make([]bool, n+1, n+1)
	sqrtn := int(math.Sqrt(float64(n)))
	for i := 2; i <= sqrtn; i++ {
		if !a[i] {
			for j := i * i; j <= n; j += i {
				a[j] = true
			}
		}
	}
	ps := make([]int, 0, 1000)
	for i := 2; i <= n; i++ {
		if !a[i] {
			ps = append(ps, i)
		}
	}
	return ps
}

func TestIsPrimeAgainstBaseline(t *testing.T) {
	for n := -1; n < 1000; n++ {
		p := primes.IsPrime(n)
		q := baselineIsPrime(n)
		if p != q {
			t.Errorf("ISPrimeAgainstBaseline(%d) == %v, want %v", n, p, q)
		}
	}
	if err := quick.CheckEqual(primes.IsPrime, baselineIsPrime, &quick.Config{MaxCount: 50}); err != nil {
		t.Error(err)
	}
}

func TestSieveAgainstBaseline(t *testing.T) {
	ns := []int{0, 1, 2, 3, 10000000}
	for _, n := range ns {
		ps := primes.Sieve(n)
		qs := baselineSieve(n)
		if len(ps) != len(qs) {
			t.Errorf("SieveAgainstBaseline(%d): len == %d, want %d", n, len(ps), len(qs))
			break
		}
		for i, p := range ps {
			if p != qs[i] {
				t.Errorf("SieveAgainstBaseline(%d): [%d] == %d, want %d", n, i, p, qs[i])
				break
			}
		}
	}
}

// Something to make sure the benchmarks are not optimized down to nothing
var nprimes int

func benchmarkSieve(b *testing.B, sieve func(n int) []int) int {
	ps := sieve(b.N)
	return len(ps)
}

// We expect primes.Sieve to be twice as fast as the baseline since it
// improves upon it by considering only odd numbers
func BenchmarkSieve(b *testing.B) {
	nprimes += benchmarkSieve(b, primes.Sieve)
}

func BenchmarkBaselineSieve(b *testing.B) {
	nprimes -= benchmarkSieve(b, baselineSieve)
}

func benchmarkIsPrime(b *testing.B, isPrime func(n int) bool) int {
	n := 0
	for i := 0; i < b.N; i++ {
		if isPrime(n) {
			n++
		}
	}
	return n
}

func BenchmarkIsPrime(b *testing.B) {
	nprimes += benchmarkIsPrime(b, primes.IsPrime)
}

func BenchmarkBaselineIsPrime(b *testing.B) {
	nprimes -= benchmarkIsPrime(b, baselineIsPrime)
}

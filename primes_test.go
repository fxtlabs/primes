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

	"github.com/fxtlabs/primes"
)

func TestPi(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{10, 4},
		{100, 25},
		{1000, 168},
		{10000, 1229},
		{100000, 9592},
		{1000000, 78498},
		{10000000, 664579},
		{100000000, 5761455},
		{1000000000, 50847534},
		{104730, 10000},
	}
	// Maximum relative error when Pi(n) returns an estimate
	const epsMax = 0.01
	for _, c := range cases {
		pi, ok := primes.Pi(c.n)
		if ok {
			if pi != c.want {
				t.Errorf("Pi(%d) == (%d,true), want %d", c.n, pi, c.want)
			}
		} else {
			eps := math.Abs(float64(pi-c.want) / float64(c.want))
			if eps >= epsMax {
				t.Errorf("Pi(%d) == (%d,false), want %d; eps=%f", c.n, pi, c.want, eps)
			}
		}

	}
}

func TestIsPrime(t *testing.T) {
	// Check a few simple cases
	cases := []struct {
		n    int
		want bool
	}{
		{-1, false},
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{6, false},
		{1000003, true},
		{1000037, true},
	}
	for _, c := range cases {
		p := primes.IsPrime(c.n)
		if p != c.want {
			t.Errorf("IsPrime(%d) == %v, want %v", c.n, p, c.want)
		}
	}

	// Check against each continguous sequence of primes that the primes
	// are classified as primes and the numbers in between as not.
	contiguousPrimes := [][]int{
		{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47},
		{127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181},
		{877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953},
		{2089, 2099, 2111, 2113, 2129, 2131, 2137, 2141, 2143, 2153},
		{9857, 9859, 9871, 9883, 9887, 9901, 9907, 9923, 9929, 9931},
		{1000003, 1000033, 1000037},
	}
	for _, ps := range contiguousPrimes {
		for _, p := range ps {
			if !primes.IsPrime(p) {
				t.Errorf("IsPrime(%d) == false, want true", p)
			}
		}
		for i := 1; i < len(ps); i++ {
			for n := ps[i-1] + 1; n < ps[i]; n++ {
				if primes.IsPrime(n) {
					t.Errorf("IsPrime(%d) == true, want false", n)
				}
			}
		}
	}

	// Check that the numbers obtain by multiplying any two of the following
	// prime factors are classified as not primes
	factors := []int{2, 3, 41, 157, 953, 2141, 9929}
	for i, p := range factors {
		for j := 0; j <= i; j++ {
			q := factors[j]
			n := p * q
			if primes.IsPrime(n) {
				t.Errorf("IsPrime(%d) == true, want false", n)
			}
		}
	}
}

func TestCoprime(t *testing.T) {
	ps := []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
		127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181,
		877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953,
		2089, 2099, 2111, 2113, 2129, 2131, 2137, 2141, 2143, 2153,
		9857, 9859, 9871, 9883, 9887, 9901, 9907, 9923, 9929, 9931,
		1000003, 1000033, 1000037,
	}

	// Distinct primes are always relatively coprime
	for i, p := range ps {
		for j := 0; j <= i; j++ {
			q := ps[j]
			got := primes.Coprime(p, q)
			want := p != q
			if got != want {
				t.Errorf("Coprime(%d,%d) == %v, want %v", p, q, got, want)
			}
		}
	}

	// Numbers sharing a common factor are never relatively coprime
	factors := []int64{2, 5, 19, 47, 137, 877, 2089, 9931}
	for i, p := range ps {
		for j := 0; j <= i; j++ {
			q := ps[j]
			for _, f := range factors {
				a := f * int64(p)
				b := f * int64(q)
				if a <= math.MaxInt32 && b <= math.MaxInt32 {
					if primes.Coprime(int(a), int(b)) {
						t.Errorf("Coprime(%d,%d) == true, want false", a, b)
					}
				}
			}
		}
	}
}

func TestSieve(t *testing.T) {
	cases := []struct {
		n    int // input to Sieve(n)
		want int // expected largest prime <= n (0 if it does not exist)
	}{
		{-1, 0},
		{0, 0},
		{1, 0},
		{2, 2},
		{3, 3},
		{4, 3},
		{1229, 1229},
		{100, 97},
		{1000, 997},
		{10000, 9973},
		{100000, 99991},
		{200000, 199999},
		{500000, 499979},
		{1000000, 999983},
	}
	for _, c := range cases {
		ps := primes.Sieve(c.n)
		if c.want > 0 {
			if len(ps) == 0 {
				t.Errorf("Sieve(%d) == [], want something", c.n)
			} else if last := ps[len(ps)-1]; last != c.want {
				t.Errorf("Sieve(%d)[-1] == %d, want %d", c.n, last, c.want)
			}
		} else {
			if l := len(ps); l > 0 {
				t.Errorf("|Sieve(%d)| == %d, want 0", c.n, l)
			}
		}
	}
}

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
	"fmt"
	"math"

	"github.com/fxtlabs/primes"
)

func ExampleSieve() {
	// Generate the prime numbers less than or equal to 20
	ps := primes.Sieve(20)
	fmt.Println(ps)

	// Output: [2 3 5 7 11 13 17 19]
}

func ExampleIsPrime() {
	// Check which numbers are prime
	ns := []int{2, 23, 49, 73, 98765, 1000003, 10007 * 10009, math.MaxInt32}
	for _, n := range ns {
		if primes.IsPrime(n) {
			fmt.Printf("%d is prime\n", n)
		} else {
			fmt.Printf("%d is composite\n", n)
		}
	}

	// Output:
	// 2 is prime
	// 23 is prime
	// 49 is composite
	// 73 is prime
	// 98765 is composite
	// 1000003 is prime
	// 100160063 is composite
	// 2147483647 is prime
}

func ExampleCoprime() {
	// Check which combinations are coprime
	ns := []int{2, 3, 4, 5, 6}
	for i, a := range ns {
		for _, b := range ns[i+1:] {
			if primes.Coprime(a, b) {
				fmt.Printf("%d and %d are coprime\n", a, b)
			}
		}
	}

	// Output:
	// 2 and 3 are coprime
	// 2 and 5 are coprime
	// 3 and 4 are coprime
	// 3 and 5 are coprime
	// 4 and 5 are coprime
	// 5 and 6 are coprime
}

func ExamplePi() {
	// Check how many prime numbers are less than or equal to n
	ns := []int{6, 11, 23, 12345}
	for _, n := range ns {
		pi, ok := primes.Pi(n)
		if ok {
			fmt.Printf("There are %d primes in [0,%d]\n", pi, n)
		} else {
			fmt.Printf("There are approximately %d primes in [0,%d]\n", pi, n)
		}
	}

	// Output:
	// There are 3 primes in [0,6]
	// There are 5 primes in [0,11]
	// There are 9 primes in [0,23]
	// There are approximately 1465 primes in [0,12345]
}

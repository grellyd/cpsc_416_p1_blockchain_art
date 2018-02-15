package tests
import (
	"fmt"
	"strings"
	"testing"
	"minerlib/compute"
)
// Inconsistent, as the naturing of random hash finding is.
// Worst: times out.
// Best: Data 1-5 + DataConcurrent 1-7 in 215s with c@8

func TestData(t *testing.T) {
	var tests = []struct {
		data    []byte
		difficulty uint8
	}{
		{
			[]byte("base case"),
			1,
		},
		{
			[]byte("test"),
			3,
		},
		// ^^ completes in 1.6s -> 5s @ c=3 -> <1s @ c=8
		{
			[]byte("thisisastring"),
			5,
		},
		// ^^ completes in 18s -> 13s @ c=3 -> <2s @ c=8
		// {
		// 	[]byte("bestchai"),
		// 	6,
		// },
		// {
		// 	[]byte("bestchai"),
		// 	7,
		// },
		// {
		// 	[]byte("Ioe948%*(F)4"),
		// 	10,
		// },
		// ^^ fails after 600s
		// {
		// 	[]byte("nonce-ahoy"),
		// 	32,
		// },
		// {
		// 	[]byte("bestchai"),
		// 	63,
		// },
	}
	for _, test := range tests {
		nonce, err := compute.Data(test.data, test.difficulty)
		if err != nil {
			t.Errorf("Bad Exit: \"TestData(%v)\" produced err: %v", test, err)
		}
		hash := compute.MD5Hash(test.data, nonce)
		fmt.Printf("Difficulty: %d, Nonce: %d, Hash: %s\n", test.difficulty, nonce, hash)
		// sanity check num zeros as using same validity functions as mine.Data
		numPresentZeros := strings.Count(hash, "0")
		if !compute.Valid(hash, test.difficulty) || numPresentZeros < int(test.difficulty) {
			t.Errorf("Bad Exit: Not enough zeros with %d! %d instead of %d", nonce, numPresentZeros, test.difficulty)
		}
	}
}

func TestDataConcurrent(t *testing.T) {
	var tests = []struct {
		data    []byte
		difficulty uint8
	}{
		{
			[]byte("base case"),
			1,
		},
		{
			[]byte("test"),
			3,
		},
		// ^^ completes in 1.6s -> 5s @ c=3 -> <1s @ c=8
		{
			[]byte("thisisastring"),
			5,
		},
		// ^^ completes in 18s -> 13s @ c=3 -> <2s @ c=8
		{
			[]byte("bestchai"),
			6,
		},
		{
			[]byte("bestchai"),
			7,
		},
		// {
		// 	[]byte("Ioe948%*(F)4"),
		// 	10,
		// },
		// ^^ fails after 600s
		// {
		// 	[]byte("nonce-ahoy"),
		// 	32,
		// },
		// {
		// 	[]byte("bestchai"),
		// 	63,
		// },
	}
	for _, test := range tests {
		nonce, err := compute.DataConcurrent(test.data, test.difficulty)
		if err != nil {
			t.Errorf("Bad Exit: \"TestData(%v)\" produced err: %v", test, err)
		}
		hash := compute.MD5Hash(test.data, nonce)
		fmt.Printf("Difficulty: %d, Nonce: %d, Hash: %s\n", test.difficulty, nonce, hash)
		// sanity check num zeros as using same validity functions as mine.Data
		numPresentZeros := strings.Count(hash, "0")
		if !compute.Valid(hash, test.difficulty) || numPresentZeros < int(test.difficulty) {
			t.Errorf("Bad Exit: Not enough zeros with %d! %d instead of %d", nonce, numPresentZeros, test.difficulty)
		}
	}
}

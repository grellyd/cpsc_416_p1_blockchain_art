package mine_test
import (
	"minerlib/mine"
	"strings"
	"testing"
)

func TestData(t *testing.T) {
	var tests = []struct {
		data    []byte
		difficulty uint8
	}{
		{
			[]byte("base case"),
			0,
		},
		{
			[]byte("test"),
			3,
		},
		// ^^ completes in 1.6s
		{
			[]byte("thisisastring"),
			5,
		},
		// ^^ completes in 18s
		{
			[]byte("Ioe948%*(F)4"),
			10,
		},
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
		nonce, err := mine.Data(test.data, test.difficulty)
		if err != nil {
			t.Errorf("Bad Exit: \"TestData(%v)\" produced err: %v", test, err)
		}
		hash := mine.MD5Hash(test.data, nonce)
		// sanity check num zeros as using same validity functions as mine.Data
		numPresentZeros := strings.Count(hash, "0")
		if !mine.Valid(hash, test.difficulty) || numPresentZeros < int(test.difficulty) {
			t.Errorf("Bad Exit: Not enough zeros with %d! %d instead of %d", nonce, numPresentZeros, test.difficulty)
		}
	}
}

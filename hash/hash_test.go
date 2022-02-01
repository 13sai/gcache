package hash

import (
	"testing"

	"github.com/spf13/cast"
)

func TestHash(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		return cast.ToUint32(string(data))
	})

	hash.Add("6", "2", "4")

	cases := map[string]string{
		"2":  "2",
		"11": "3",
		"23": "5",
		"34": "10",
	}

	for k, v := range cases {
		if hash.Get(k) != v {
			t.Errorf("k=%s, v=%s", k, v)
		}
	}

	hash.Add("9")

	cases["27"] = "9"

	for k, v := range cases {
		if hash.Get(k) != v {
			t.Errorf("2k=%s, v=%s", k, v)
		}
	}
}

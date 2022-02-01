package gcache

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	var db = map[string]string{
		"Bob":  "250",
		"Jack": "380",
		"Tom":  "517",
	}

	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			t.Logf("search key %s", key)

			if v, found := db[key]; found {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not found", key)
		},
	))

	addr := "127.0.0.1:9000"
	peers := NewHTTPPool(addr)
	t.Logf("gcache is running at %s", addr)
	t.Fatal(http.ListenAndServe(
		addr, peers,
	))
}

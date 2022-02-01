package gcache

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	var db = map[string]string{
		"Bob":  "250",
		"Jack": "380",
		"Tom":  "517",
	}

	loadCounts := make(map[string]int, len(db))
	g := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			t.Logf("search key %s", key)

			if v, found := db[key]; found {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not found", key)
		},
	))

	for name, score := range db {
		if view, err := g.Get(name); err != nil || view.String() != score {
			t.Fatalf("%s no score", name)
		}

		if _, err := g.Get(name); err != nil || loadCounts[name] > 1 {
			t.Fatalf("%s cache missed", name)
		}
	}

	if _, err := g.Get("abcd"); err != nil {
		t.Fatal("not get")
	}
}

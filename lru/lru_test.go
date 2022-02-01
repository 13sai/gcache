package lru

import "testing"

type String string

func (d String) Len() int {
	return len(d)
}

func (d String) GetString() string {
	return string(d)
}

func TestLRU(t *testing.T) {
	lru := New(100, nil)

	lru.Add("k1", String("k1-v"))
	if v, ok := lru.Get("k1"); !ok {
		t.Log("k1 not found")
	} else {
		t.Logf("k1 found=%s", v.(String).GetString())
	}
}

func TestLRU2(t *testing.T) {
	k1, k2, k3, v1, v2, v3 := "k1", "k2", "k3", "v1", "v2", "v3"
	l := len(k1 + v1 + k2 + v2)
	lru := New(uint32(l), func(s string, v Value) {
		t.Logf("remove %s-%s", s, v.(String).GetString())
	})
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("k1"); !ok {
		t.Log(lru.Len())
	}
}

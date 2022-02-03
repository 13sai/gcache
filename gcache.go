package gcache

import (
	"errors"
	"sync"

	pb "github.com/13sai/gcache/gcachepb"

	"github.com/13sai/gcache/singleflight"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name   string
	getter Getter
	cache  cache
	peers  PeerPicker
	loader *singleflight.Group
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	v := ByteView{b: cloneBytes(bytes)}
	g.cache.add(key, v)
	return v, nil
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: res.Value}, nil
}

func (g *Group) load(key string) (ByteView, error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if v, err := g.getFromPeer(peer, key); err != nil {
					return v, nil
				}
			}
		}
		return g.getLocally(key)
	})

	if err == nil {
		return view.(ByteView), nil
	}

	return ByteView{}, nil
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is nil")
	}
	if v, ok := g.cache.get(key); ok {
		return v, nil
	}

	return g.load(key)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	g.peers = peers
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes uint32, getter Getter) *Group {
	if getter == nil {
		panic("Getter is nil")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		cache:  cache{cacheBytes: cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

package gcache

type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

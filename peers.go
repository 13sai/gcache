package gcache

import pb "github.com/13sai/gcache/gcachepb"

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

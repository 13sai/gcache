package gcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	pb "github.com/13sai/gcache/gcachepb"
	"github.com/13sai/gcache/hash"
	"google.golang.org/protobuf/proto"
)

const defaultBasePath = "/cache/"
const defaultReplicas = 30

type HTTPPool struct {
	self          string
	basePath      string
	mu            sync.Mutex
	peers         *hash.Map
	httpGetterMap map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (hp *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, hp.basePath) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("HTTPPool serving unexpected path: " + r.URL.Path))
		return
	}

	parts := strings.SplitN(r.URL.Path[len(hp.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&pb.Response{Value: view.Clone()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(body)
}

func (hp *HTTPPool) Set(peers ...string) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.peers = hash.New(defaultReplicas, nil)
	hp.peers.Add(peers...)
	hp.httpGetterMap = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		hp.httpGetterMap[peer] = &httpGetter{baseURL: peer + hp.basePath}
	}
}

func (hp *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	if peer := hp.peers.Get(key); peer != "" && peer != hp.self {
		return hp.httpGetterMap[peer], true
	}

	return nil, false
}

type httpGetter struct {
	baseURL string
}

func (hg *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	path := fmt.Sprintf("%s%s/%s",
		hg.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	res, err := http.Get(path)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server status %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading resp body %v", err)
	}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding resp body: %v", err)
	}
	return nil
}

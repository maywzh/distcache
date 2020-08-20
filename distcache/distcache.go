package distcache

import (
	pb "distcache/distcachepb"
	"fmt"
	"log"
	"sync"
)

// A Node is a cache namespace and associated data loaded spread over
type Node struct {
	name      string // The instance name
	getter    Getter // The missing callback
	mainCache cache  // The concurrent cache
	peers     PeerPicker
}

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc implements Getter with function
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu    sync.RWMutex
	nodes = make(map[string]*Node)
)

// NewNode create a new instance of Node
func NewNode(name string, cacheBytes int64, getter Getter) *Node {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Node{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	nodes[name] = g
	return g
}

// GetNode returns the named node previously created with NewNode, or
// nil if there's no such node.
func GetNode(name string) *Node {
	mu.RLock()
	g := nodes[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Node) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[DistCache] hit")
		return v, nil
	}
	return g.load(key)
}

// Get from locally database
func (g *Node) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache add bytes to cache
func (g *Node) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Node) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Node) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}
	return g.getLocally(key)
}

func (g *Node) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Node: g.name,
		Key:  key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

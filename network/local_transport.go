package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
}

// NewLocalTransport initializes a LocalTransport instance
func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// Addr returns node address
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

// Connect adds new node to current node's peer list
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

// Peers returns a list of a node's peers
func (t *LocalTransport) Peers() []*LocalTransport {
	t.lock.Lock()
	defer t.lock.Unlock()

	var peers []*LocalTransport
	for _, peer := range t.peers {
		if p, ok := t.peers[peer.addr]; ok {
			peers = append(peers, p)
		}
	}

	return peers
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s could not connect to %s: ", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(payload),
	}

	return nil
}

func (t *LocalTransport) Broadcast(msg []byte) error {
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), msg); err != nil {
			return err
		}
	}

	return nil
}

package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC),
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

// TODO: get node peers
//func (t *LocalTransport) Peers() []*Transport {
//	t.lock.Lock()
//	defer t.lock.Unlock()
//
//	var peers []*Transport
//	for _, peer := range t.peers {
//		if p, ok := t.peers[peer.addr]; ok {
//			peers = append(peers, p)
//		}
//	}
//
//	return peers
//}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s could not connect to %s: ", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: payload,
	}

	return nil
}

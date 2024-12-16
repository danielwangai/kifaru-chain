package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalTransportConnect(t *testing.T) {
	// create new local transports
	addrA := NetAddr("A")
	addrB := NetAddr("B")
	trA := NewLocalTransport(addrA)
	trB := NewLocalTransport(addrB)

	assert.Equal(t, trA.Addr(), addrA)
	assert.Equal(t, trB.Addr(), addrB)

	// connect nodes
	trA.Connect(trB)
	trB.Connect(trA)

	//assert.Equal(t, trA.Peers()[0], trB.Peers()[0])

	//assert.Equal(t, trA.peers[trB.addr], trB)
	//assert.Equal(t, trB.peers[trA.addr], trA)
}

//func TestLocalTransportSendMessage(t *testing.T) {
//	trA := NewLocalTransport("A")
//	trB := NewLocalTransport("B")
//
//	trA.Connect(trB)
//	trB.Connect(trA)
//
//	assert.Equal(t, trA.peers[trB.addr], trB)
//	assert.Equal(t, trB.peers[trA.addr], trA)
//
//	payload := []byte("Hello")
//	assert.Nil(t, trA.SendMessage(trB.addr, payload))
//	//
//	//rpc := <-trA.Consume()
//	//assert.Equal(t, rpc.Payload, payload)
//	//assert.Equal(t, rpc.From, trA.addr)
//}

package network

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, trA.Connect(trB))
	assert.Nil(t, trB.Connect(trA))

	assert.Equal(t, trA.peers[trB.addr], trB)
	assert.Equal(t, trB.peers[trA.addr], trA)
}

func TestLocalTransportSendMessage(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")

	assert.Nil(t, trA.Connect(trB))
	assert.Nil(t, trB.Connect(trA))

	//assert.Equal(t, len(trA.Peers()), 1)
	//assert.Equal(t, len(trB.Peers()), 1)

	payload := []byte("Hello")
	assert.Nil(t, trA.SendMessage(trB.addr, payload))

	rpc := <-trB.Consume()
	buf := make([]byte, len(payload))
	n, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, n, len(payload))

	assert.Equal(t, buf, payload)
	assert.Equal(t, rpc.From, trA.Addr())
}

func TestBroadcast(t *testing.T) {
	trA := NewLocalTransport("A")
	trB := NewLocalTransport("B")
	trC := NewLocalTransport("C")

	assert.Nil(t, trA.Connect(trB))
	assert.Nil(t, trA.Connect(trC))

	msg := []byte("Hello")
	assert.Nil(t, trA.Broadcast(msg))

	rpcB := <-trB.Consume()
	b, err := ioutil.ReadAll(rpcB.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcC := <-trC.Consume()
	b, err = ioutil.ReadAll(rpcC.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
}

package main

import (
	"fmt"
	"github.com/danielwangai/kifaru-block/network"
	"time"
)

func main() {
	fmt.Println("Kifaru Chain")
	// local node
	trLocal := network.NewLocalTransport("LocalTR")
	// remote node
	trRemote := network.NewLocalTransport("RemoteTR")

	// connect nodes
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			// send message from remote to local node
			trRemote.SendMessage(trLocal.Addr(), []byte("Hello World"))
			time.Sleep(time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trRemote, trLocal},
	}
	server := network.NewServer(opts)
	server.Start()
}

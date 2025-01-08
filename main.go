package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"math/rand"

	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/danielwangai/kifaru-block/network"
	"github.com/sirupsen/logrus"
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
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(time.Second)
		}
	}()

	privKey := crypto.GeneratePrivateKey()

	opts := network.ServerOpts{
		ID:         "LOCAL",
		Transports: []network.Transport{trLocal},
		PrivateKey: privKey,
	}
	server := network.NewServer(opts)
	server.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000000000)), 10))
	tx := crypto.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(crypto.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

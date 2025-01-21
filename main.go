package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/danielwangai/kifaru-block/network"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Kifaru Chain")
	trLocal := network.NewLocalTransport("LocalTR")
	trRemote1 := network.NewLocalTransport("RemoteTR1")
	trRemote2 := network.NewLocalTransport("RemoteTR2")
	trRemote3 := network.NewLocalTransport("RemoteTR3")

	// connect nodes
	trLocal.Connect(trRemote1)
	trRemote1.Connect(trRemote2)
	trRemote2.Connect(trRemote3)
	trRemote1.Connect(trLocal)

	initRemoteServers([]network.Transport{trRemote1, trRemote2, trRemote3})

	go func() {
		for {
			// send message from remote to local node
			if err := sendTransaction(trRemote1, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)
		trLate := network.NewLocalTransport("LATE_REMOTE")
		trRemote3.Connect(trLate)
		lateServer := makeServer(string(trLate.Addr()), trLate, nil)
		go lateServer.Start()
	}()

	// start local server
	privKey := crypto.GeneratePrivateKey()

	localServer := makeServer("LOCAL", trLocal, privKey)
	localServer.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i, tr := range trs {
		id := fmt.Sprintf("REMOTE-%d", i)
		s := makeServer(id, tr, nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		ID:         id,
		Transports: []network.Transport{tr},
		PrivateKey: pk,
	}

	server, err := network.NewServer(opts)
	if err != nil {
		logrus.Fatal(err)
	}

	return server
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

package network

import (
	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/sirupsen/logrus"

	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	blockTime   time.Duration
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

// NewServer initializes new server
func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		blockTime:   opts.BlockTime,
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}
}

// Start the server
func (s *Server) Start() {
	s.InitTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("Server stopped")
}

func (s *Server) createNewBlock() {
	// take all transaction, hash and create a new block
	fmt.Println("create new block")
}

func (s *Server) handleTransaction(tx *crypto.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(crypto.TxHasher{})
	fmt.Println("Hash: ", hash)
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("tx already in mempool")
		return nil
	}

	// fmt.Println("Has Hash: ", hash)
	logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Info("tx has been added to mempool")
	return s.memPool.Add(tx)
}

func (s *Server) InitTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

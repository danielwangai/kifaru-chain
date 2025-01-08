package network

import (
	"bytes"
	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/sirupsen/logrus"

	"fmt"
	"time"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID            string
	Logger        *logrus.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
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
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.Logger == nil {
		logger := logrus.New()
		opts.Logger = logger
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		blockTime:   opts.BlockTime,
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}
	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}

	if s.isValidator {
		go s.startValidatorBlockProducer()
	}
	s.ServerOpts = opts

	return s
}

// Start the server
func (s *Server) Start() {
	if s.blockTime == time.Duration(0) {
		s.blockTime = defaultBlockTime
	}
	s.InitTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			// decode message from rpc
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.WithError(err)
			}

			// process message
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.WithError(err)
			}
		case <-s.quitCh:
			break free
		}
	}

	s.Logger.Info("Server stopped")
}

func (s *Server) startValidatorBlockProducer() {
	ticker := time.NewTicker(s.blockTime)
	s.Logger.Infof("msg=starting validator, loop blocktime=%s", s.blockTime)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) createNewBlock() {
	// take all transaction, hash and create a new block
	fmt.Println("create new block")
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *crypto.Transaction:
		return s.processTransaction(t)
	}

	return nil
}

// handles checks before adding a new transaction to the mempool
func (s *Server) processTransaction(tx *crypto.Transaction) error {
	hash := tx.Hash(crypto.TxHasher{})
	if s.memPool.Has(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Infof("msg=adding new tx to mempool,hash=%s", hash)

	// broadcast to peers
	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

func (s *Server) broadcast(msg []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(msg); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) broadcastTx(tx *crypto.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(crypto.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
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

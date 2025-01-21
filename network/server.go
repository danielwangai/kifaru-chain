package network

import (
	"bytes"

	"github.com/danielwangai/kifaru-block/crypto"
	"github.com/danielwangai/kifaru-block/types"
	"github.com/sirupsen/logrus"

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
	blockchain  *crypto.Blockchain
	blockTime   time.Duration
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

// NewServer initializes new server
func NewServer(opts ServerOpts) (*Server, error) {
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

	chain, err := crypto.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(1000),
		blockchain:  chain,
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

	return s, nil
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

// adds new block to the blockchain
// currently all transactions in the mempool are added to the block.
// TODO: figure out a how many transactions are to be added to the block
func (s *Server) createNewBlock() error {
	currentHeader, err := s.blockchain.GetHeaderByHeight(s.blockchain.Height())
	if err != nil {
		return err
	}

	txs := s.memPool.Pending()

	block, err := crypto.NewBlockFromPrevHeader(currentHeader, txs)
	if err != nil {
		return err
	}

	block.Sign(s.PrivateKey)

	if err := s.blockchain.AddBlock(block); err != nil {
		return err
	}

	s.memPool.ClearPending()

	go s.broadcastBlock(block)

	return nil
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *crypto.Transaction:
		return s.processTransaction(t)
	case *crypto.Block:
		return s.processBlock(t)
	}

	return nil
}

// handles checks before adding a new transaction to the mempool
func (s *Server) processTransaction(tx *crypto.Transaction) error {
	hash := tx.Hash(crypto.TxHasher{})
	if s.memPool.Contains(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	s.Logger.Infof("msg=adding new tx to mempool,hash=%s, mempoolPending: %d", hash, s.memPool.PendingCount())

	// broadcast to peers
	go s.broadcastTx(tx)

	s.memPool.Add(tx)

	return nil
}

func (s *Server) processBlock(b *crypto.Block) error {
	if err := s.blockchain.AddBlock(b); err != nil {
		return err
	}
	go s.broadcastBlock(b)
	return nil
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

func (s *Server) broadcastBlock(b *crypto.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(crypto.NewGobBlockEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeBlock, buf.Bytes())
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

// temp helpers
func genesisBlock() *crypto.Block {
	header := &crypto.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 1691622800,
	}
	return crypto.NewBlock(header, []*crypto.Transaction{})
}

package crypto

import (
	"bytes"
	"github.com/danielwangai/kifaru-block/examples"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHeader_Encode_Decode(t *testing.T) {
	h := &Header{
		Version:   1,
		PrevBlock: examples.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    10,
		Nonce:     100000,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

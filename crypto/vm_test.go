package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVMAdd(t *testing.T) {
	// This test performs the following operation => 1 + 2 = 3
	// 1 = 0x1, 2 = 0x2
	// InstrPush opcode to push into the stack
	// InstrAdd = opcode to add item to stack
	data := []byte{0x01, byte(InstrPush), 0x02, byte(InstrPush), byte(InstrAdd)}
	vm := NewVM(data)
	vm.Run()

	assert.Equal(t, byte(3), byte(vm.stack.Pop().(int)))
}

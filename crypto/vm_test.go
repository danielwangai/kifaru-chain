package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVMAdd(t *testing.T) {
	// performs the following operation => 1 + 2 = 3
	// 1 => 0x1,
	// 2 => 0x2
	// InstrPush opcode to push into the stack
	// OpCodeAdd = opcode to add item to stack
	data := []byte{0x01, byte(OpCodePushInt), 0x02, byte(OpCodePushInt), byte(OpCodeAdd)}

	vm := NewVM(data)
	err := vm.Run()
	assert.Nil(t, err)

	assert.Equal(t, byte(3), byte(vm.stack.Pop().(int)))
}

func TestVMPack(t *testing.T) {
	// concatenate byte representation of GOOD
	// G => 0x47
	// O => 0x4F
	// O => 0x4F
	// D => 0x44
	// OpCodePack => 0x0b to combine the hex characters
	data := []byte{
		0x47, byte(OpCodePushByte),
		0x4F, byte(OpCodePushByte),
		0x4F, byte(OpCodePushByte),
		0x44, byte(OpCodePushByte),
		byte(OpCodePack),
	}

	vm := NewVM(data)
	err := vm.Run()
	assert.Nil(t, err)

	result := vm.stack.Pop().([]byte)
	assert.Equal(t, "GOOD", string(result))
}

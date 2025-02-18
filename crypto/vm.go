package crypto

type OpCode byte

const (
	OpCodePushInt  OpCode = 0x0a // 11
	OpCodePushByte OpCode = 0x0b
	OpCodePack     OpCode = 0x0c
	OpCodeAdd      OpCode = 0x0d
	OpCodeSub      OpCode = 0x0e
	OpCodeMul      OpCode = 0x0f
	OpCodeDiv      OpCode = 0x10
)

type Stack struct {
	data []any // data store
	ptr  int   // pointer position where an operation is in the stack
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		ptr:  0,
	}
}

// Push adds item to stack
// Fixme: behaves like a queue
func (s *Stack) Push(v any) {
	// Check if the stack is full
	if s.ptr >= len(s.data) {
		panic("stack overflow")
	}

	s.data[s.ptr] = v
	s.ptr++
}

// Pop deletes element at the top of the stack
// Fixme: behaves like a queue
func (s *Stack) Pop() any {
	if s.ptr == 0 {
		panic("stack underflow")
	}
	// delete item at index 0
	head := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)
	s.ptr--

	return head
}

type VM struct {
	data  []byte
	ip    int // instruction pointer - points to where the current instruction is
	stack *Stack
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: NewStack(1024),
	}
}

func (vm *VM) Run() error {
	for {
		instr := OpCode(vm.data[vm.ip])

		if err := vm.Exec(instr); err != nil {
			return err
		}
		vm.ip++

		if vm.ip > len(vm.data)-1 {
			break
		}
	}

	return nil
}

func (vm *VM) Exec(instr OpCode) error {
	switch instr {
	case OpCodePushInt:
		// add to stack
		vm.stack.Push(int(vm.data[vm.ip-1]))
	case OpCodePushByte:
		// add to stack
		vm.stack.Push(byte(vm.data[vm.ip-1]))
	case OpCodePack:
		n := vm.stack.ptr
		b := make([]byte, n)

		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)
	case OpCodeAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a + b)
	case OpCodeSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a - b)
	case OpCodeMul:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a * b)
	case OpCodeDiv:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a / b)
	}

	return nil
}

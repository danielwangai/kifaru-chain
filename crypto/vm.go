package crypto

type Instruction byte

const (
	InstrPush Instruction = 0x0a // 11
	InstrAdd  Instruction = 0x0b
)

type Stack struct {
	data []any // data store
	ptr  int   // pointer position where an operation is in the stack
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
	}
}

// Push adds item to stack
// Fixme: behaves like a queue
func (s *Stack) Push(v any) {
	s.data[s.ptr] = v
	s.ptr++
}

// Pop deletes element at the top of the stack
// Fixme: behaves like a queue
func (s *Stack) Pop() any {
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
		instr := Instruction(vm.data[vm.ip])

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

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrPush:
		// add to stack
		vm.stack.Push(int(vm.data[vm.ip-1]))
	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a + b)
	}

	return nil
}

package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
)

const (
	// Op codes
	opNop = iota
	opPush
	opPop
	opSum
	opSub
	opMul
	opPrint
	opPrintln
	opRand
	opCmp
	opDup
	opSwap
	opJmp
	opHalt
)

type operation struct {
	name string
	args int
}

var ops = map[int]operation{
	opNop:     operation{"noop", 0},
	opPush:    operation{"push", 1},
	opPop:     operation{"pop", 0},
	opSum:     operation{"sum", 0},
	opSub:     operation{"subtract", 0},
	opMul:     operation{"multiply", 0},
	opPrint:   operation{"print", 0},
	opPrintln: operation{"println", 0},
	opRand:    operation{"random", 0},
	opCmp:     operation{"compare", 0},
	opDup:     operation{"duplicate", 0},
	opSwap:    operation{"swap", 0},
	opJmp:     operation{"jump", 1},
	opHalt:    operation{"halt", 0},
}

// VM represents the state of the virtual machine.
type VM struct {
	instructions []int // slice containing all the code to execute
	stack        []int
	pc           int // program counter
}

// Just for debugging
func (v *VM) trace() {
	op := ops[v.instructions[v.pc]]
	args := v.instructions[v.pc+1 : v.pc+op.args+1]
	fmt.Printf("%04d: %s %v \t%v\n", v.pc, op.name, args, v.stack)
}

// Push an int onto the stack.
func (v *VM) push(n int) {
	v.stack = append(v.stack, n)
}

// Pop an int off the stack.
func (v *VM) pop() int {
	if len(v.stack) < 1 {
		panic("stack underflow")
	}
	val := v.stack[len(v.stack)-1]
	v.stack = v.stack[:len(v.stack)-1]
	return val
}

// Math functions, two inputs are popped from the stack for use
// in the operation and the result is pushed back on.
func (v *VM) sum() {
	a, b := v.pop(), v.pop()
	v.push(a + b)
}

func (v *VM) subtract() {
	a, b := v.pop(), v.pop()
	v.push(b - a)
}

func (v *VM) multiply() {
	a, b := v.pop(), v.pop()
	v.push(a * b)
}

// Seed the PRNG
func seed() error {
	var s int64
	if err := binary.Read(crand.Reader, binary.BigEndian, &s); err != nil {
		return err
	}
	rand.Seed(s)
	return nil
}

// pushes a pseudorandom int to the stack
func (v *VM) rand() {
	v.push(rand.Int())
}

// Pop and print value
func (v *VM) print() {
	val := v.pop()
	fmt.Print(val)
}

// pop and print value with a new line suffix
func (v *VM) println() {
	val := v.pop()
	fmt.Println(val)
}

func (v *VM) compare() {
	var result int
	a, b := v.pop(), v.pop()
	if b < a {
		result = -1
	}
	if b > a {
		result = 1
	}
	v.push(result)
}

func (v *VM) dup() {
	val := v.pop()
	v.push(val)
	v.push(val)
}

func (v *VM) swap() {
	a, b := v.pop(), v.pop()
	v.push(a)
	v.push(b)
}

// run a program
func (v *VM) run(instructions []int) {
	v.instructions = instructions

	// reset program counter and stack so that programs can be run
	// one after another in main()
	v.pc = 0
	v.stack = []int{}

	for {
		v.trace()

		op := v.instructions[v.pc]
		v.pc++

		switch op {
		case opNop:
			continue
		case opPush:
			v.push(v.instructions[v.pc])
			v.pc++
		case opPop:
			v.pop()
		case opPrint:
			v.print()
		case opPrintln:
			v.println()
		case opSum:
			v.sum()
		case opSub:
			v.subtract()
		case opMul:
			v.multiply()
		case opRand:
			v.rand()
		case opCmp:
			v.compare()
		case opDup:
			v.dup()
		case opSwap:
			v.swap()
		case opJmp:
			v.pc = v.instructions[v.pc]
		case opHalt:
			return
		}
	}
}

func init() {
	if err := seed(); err != nil {
		panic(err)
	}
}

func main() {
	// Some example "programs" just feeding them in as a slice of ints for now,
	// can read them in from a file later.
	code := []int{
		opPush, 10,
		opPush, 15,
		opSum,
		opPrintln,
		opPush, 100,
		opPush, 150,
		opSum,
		opPrintln,
		opHalt,
	}
	fmt.Printf("code: %v\n", code)

	v := &VM{}
	v.run(code)
	fmt.Println()

	code2 := []int{
		opPush, 10,
		opPush, 5,
		opMul,
		opPrintln,
		opHalt,
	}
	fmt.Printf("code: %v\n", code2)

	v.run(code2)
	fmt.Println()

	code3 := []int{
		opRand,
		opPrintln,
		opHalt,
	}
	fmt.Printf("code: %v\n", code3)
	v.run(code3)
}

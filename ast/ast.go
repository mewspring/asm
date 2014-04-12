// Package ast declares the types used to represent abstract syntax trees of
// assembly source code.
package ast

// Inst represent an instruction which consists of an operation (opcode) and
// zero or more operands (arguments). Some instructions are not implemented in
// hardware; these pseudo-instructions represent one or more native
// instructions.
type Inst struct {
	// Op represents the operation of the instruction.
	Op Op
	// Args is a slice of zero or more operands.
	Args []Arg
}

// An Op specifies the operation to perform of an instruction; such as ADD, BEQ,
// XOR, etc.
type Op int

// Arg represent an instruction operand which specifies an address, a label, a
// register or an immediate value.
type Arg interface{}

// Addr represent an address which specifies a memory location.
type Addr uint32

// Ident represent an identifier which specifies a memory location.
type Ident string

// Reg represent a register from r0 through r31.
type Reg uint8

// Int represent a 16-bit signed integer value.
type Int int16

// Uint represent a 16-bit unsigned integer value.
type Uint uint16

// A Directive is a command to the assembler which specifies memory alignment,
// section names and program origin among others.
type Directive int

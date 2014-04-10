// Package ast declares the types used to represent abstract syntax trees of
// assembly source code.
package ast

// An Instruction consists of an operation (opcode) and zero or more operands
// (arguments). Some instructions are not implemented in hardware; these
// pseudo-instructions represent one or more native instructions.
type Instruction struct{}

// A Directive is a command to the assembler which specifies memory alignment,
// section names and program origin among others.
type Directive struct{}

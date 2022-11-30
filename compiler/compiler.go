package compiler

import (
	"curryLang/ast"
	"curryLang/code"
	"curryLang/object"
	"curryLang/token"
	"fmt"
)

type EmittedInstruction struct {
	Code code.Opcode
	Pos  int
}

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	previousInstr *EmittedInstruction
	currentInstr  *EmittedInstruction
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) CompileStatements(statements []ast.Statement) error {
	for _, s := range statements {
		err := c.Compile(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		err := c.CompileStatements(node.Statements)
		if err != nil {
			return err
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

		if node.Token.Type != token.IF {
			c.emit(code.OpPop)
		}

	case *ast.InfixExpression:

		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)

			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PrefixExpression:

		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)

		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IfElseExpression:
		err := c.compileIfExpression(node)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileIfExpression(ifExpr *ast.IfElseExpression) error {
	err := c.Compile(ifExpr.Condition)
	if err != nil {
		return err
	}

	conditionJumpPos := c.emit(code.OpJumpIfFalse, 0)

	err = c.CompileStatements(ifExpr.Consequence)
	if err != nil {
		return err
	}

	endJumpPos := c.emit(code.OpJump, 0)
	c.updateInstruction(conditionJumpPos, code.OpJumpIfFalse, c.currentInstr.Pos+3-conditionJumpPos)

	err = c.CompileStatements(ifExpr.Alternative)
	if err != nil {
		return err
	}

	c.updateInstruction(endJumpPos, code.OpJump, c.currentInstr.Pos-endJumpPos)

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.previousInstr = c.currentInstr
	c.currentInstr = &EmittedInstruction{
		Code: op,
		Pos:  pos,
	}

	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) updateInstruction(pos int, op code.Opcode, operands ...int) {
	ins := code.Make(op, operands...)

	for i := 0; i < len(ins); i++ {
		c.instructions[pos+i] = ins[i]
	}
}

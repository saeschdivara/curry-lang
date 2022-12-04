package vm

import (
	"curryLang/code"
	"curryLang/compiler"
	"curryLang/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 65536

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	globals      []object.Object
	sp           int // Always points to the next value. Top of stack is stack[sp-1]
	DebugMode    bool
}

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		globals:      make([]object.Object, GlobalsSize),
		sp:           0,
		DebugMode:    false,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	if vm.DebugMode {
		vm.DumpByteCode()
	}

	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		if vm.DebugMode {
			opDef, _ := code.Lookup(byte(op))
			fmt.Println(ip, " > ", opDef.Name)
		}

		switch op {

		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpMinus:
			right := vm.pop()
			rightType := right.Type()

			if rightType != object.INTEGER_OBJ {
				return fmt.Errorf("%s does not support minus operator", rightType)
			}

			intVal := right.(*object.Integer)
			intVal.Value *= -1

			err := vm.push(intVal)
			if err != nil {
				return err
			}

		case code.OpBang:
			right := vm.pop()
			rightType := right.Type()

			if rightType != object.BOOLEAN_OBJ {
				return fmt.Errorf("%s does not support bang operator", rightType)
			}

			boolValue := right.(*object.Boolean)
			err := vm.push(nativeBooleanToVmBoolean(!boolValue.Value))
			if err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}

		case code.OpGreaterThan, code.OpEqual, code.OpNotEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpJumpIfFalse:
			jumpVal := code.ReadUint16(vm.instructions[ip+1:])
			conditionVal := vm.pop()

			if conditionVal.Type() != object.BOOLEAN_OBJ {
				return fmt.Errorf("unsupported type for boolean jump: %s", conditionVal.Type())
			}

			boolVal, _ := conditionVal.(*object.Boolean)

			if boolVal == False {
				ip += int(jumpVal) - 1
			} else {
				ip += 2
			}

		case code.OpJump:
			jumpVal := code.ReadUint16(vm.instructions[ip+1:])
			ip += int(jumpVal)

		case code.OpSetGlobal:
			variableIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			vm.globals[variableIndex] = vm.pop()

		case code.OpGetGlobal:
			variableIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.globals[variableIndex])
			if err != nil {
				return err
			}

		default:
			def, err := code.Lookup(byte(op))
			if err != nil {
				return err
			}

			return fmt.Errorf("no implementation for %s", def.Name)
		}
	}

	if vm.DebugMode {
		fmt.Println()
		fmt.Println()
	}

	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result int64
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeComparisonInteger(op, left, right)
	}

	if leftType == object.BOOLEAN_OBJ && rightType == object.BOOLEAN_OBJ {
		return vm.executeComparisonBoolean(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeComparisonInteger(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result bool
	switch op {
	case code.OpGreaterThan:
		result = leftValue > rightValue
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue
	default:
		def, _ := code.Lookup(byte(op))
		return fmt.Errorf("unknown integer operator: %s", def.Name)
	}

	return vm.push(nativeBooleanToVmBoolean(result))
}

func (vm *VM) executeComparisonBoolean(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Boolean)
	rightValue := right.(*object.Boolean)
	var result bool
	switch op {
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue
	default:
		def, _ := code.Lookup(byte(op))
		return fmt.Errorf("unknown integer operator: %s", def.Name)
	}

	return vm.push(nativeBooleanToVmBoolean(result))
}

func nativeBooleanToVmBoolean(val bool) *object.Boolean {
	if val {
		return True
	} else {
		return False
	}
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) DumpByteCode() {
	fmt.Println(vm.instructions)
}

package cpu

// 加载指令实现
func (cpu *CPU) LDA() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	cpu.A = value
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return cycles
}

func (cpu *CPU) LDX() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	cpu.X = value
	cpu.SetZeroAndNegativeFlags(cpu.X)
	return cycles
}

func (cpu *CPU) LDY() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	cpu.Y = value
	cpu.SetZeroAndNegativeFlags(cpu.Y)
	return cycles
}

// 存储指令实现
func (cpu *CPU) STA() int {
	addr, cycles := cpu.GetOperandAddress(ZP)
	cpu.WriteByte(addr, cpu.A)
	return cycles
}

func (cpu *CPU) STX() int {
	addr, cycles := cpu.GetOperandAddress(ZP)
	cpu.WriteByte(addr, cpu.X)
	return cycles
}

func (cpu *CPU) STY() int {
	addr, cycles := cpu.GetOperandAddress(ZP)
	cpu.WriteByte(addr, cpu.Y)
	return cycles
}

// 算术指令实现
func (cpu *CPU) ADC() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)

	// 带进位的加法
	carry := uint16(0)
	if cpu.GetFlag(FLAG_C) {
		carry = 1
	}

	result := uint16(cpu.A) + uint16(value) + carry

	// 设置进位标志
	cpu.SetFlag(FLAG_C, result > 0xFF)

	// 设置溢出标志
	overflow := ((cpu.A^value)&0x80 == 0) && ((cpu.A^byte(result))&0x80 != 0)
	cpu.SetFlag(FLAG_V, overflow)

	cpu.A = byte(result)
	cpu.SetZeroAndNegativeFlags(cpu.A)

	return cycles
}

func (cpu *CPU) SBC() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)

	// 带借位的减法
	carry := uint16(0)
	if cpu.GetFlag(FLAG_C) {
		carry = 1
	}

	result := uint16(cpu.A) - uint16(value) - (1 - carry)

	// 设置进位标志
	cpu.SetFlag(FLAG_C, result <= 0xFF)

	// 设置溢出标志
	overflow := ((cpu.A^value)&0x80 != 0) && ((cpu.A^byte(result))&0x80 != 0)
	cpu.SetFlag(FLAG_V, overflow)

	cpu.A = byte(result)
	cpu.SetZeroAndNegativeFlags(cpu.A)

	return cycles
}

// 逻辑指令实现
func (cpu *CPU) AND() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	cpu.A &= value
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return cycles
}

func (cpu *CPU) EOR() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	cpu.A ^= value
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return cycles
}

func (cpu *CPU) ORA() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	cpu.A |= value
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return cycles
}

// 移位指令实现
func (cpu *CPU) ASL() int {
	carry := (cpu.A & 0x80) != 0
	cpu.A <<= 1
	cpu.SetFlag(FLAG_C, carry)
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 2
}

func (cpu *CPU) LSR() int {
	carry := (cpu.A & 0x01) != 0
	cpu.A >>= 1
	cpu.SetFlag(FLAG_C, carry)
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 2
}

func (cpu *CPU) ROL() int {
	oldCarry := cpu.GetFlag(FLAG_C)
	carry := (cpu.A & 0x80) != 0
	cpu.A <<= 1
	if oldCarry {
		cpu.A |= 0x01
	}
	cpu.SetFlag(FLAG_C, carry)
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 2
}

func (cpu *CPU) ROR() int {
	oldCarry := cpu.GetFlag(FLAG_C)
	carry := (cpu.A & 0x01) != 0
	cpu.A >>= 1
	if oldCarry {
		cpu.A |= 0x80
	}
	cpu.SetFlag(FLAG_C, carry)
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 2
}

// 比较指令实现
func (cpu *CPU) CMP() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	result := cpu.A - value
	cpu.SetFlag(FLAG_C, cpu.A >= value)
	cpu.SetZeroAndNegativeFlags(byte(result))
	return cycles
}

func (cpu *CPU) CPX() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	result := cpu.X - value
	cpu.SetFlag(FLAG_C, cpu.X >= value)
	cpu.SetZeroAndNegativeFlags(byte(result))
	return cycles
}

func (cpu *CPU) CPY() int {
	addr, cycles := cpu.GetOperandAddress(IMM)
	value := cpu.ReadByte(addr)
	result := cpu.Y - value
	cpu.SetFlag(FLAG_C, cpu.Y >= value)
	cpu.SetZeroAndNegativeFlags(byte(result))
	return cycles
}

// 分支指令实现
func (cpu *CPU) BCC() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if !cpu.GetFlag(FLAG_C) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BCS() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if cpu.GetFlag(FLAG_C) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BEQ() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if cpu.GetFlag(FLAG_Z) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BMI() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if cpu.GetFlag(FLAG_N) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BNE() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if !cpu.GetFlag(FLAG_Z) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BPL() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if !cpu.GetFlag(FLAG_N) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BVC() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if !cpu.GetFlag(FLAG_V) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

func (cpu *CPU) BVS() int {
	addr, cycles := cpu.GetOperandAddress(REL)
	if cpu.GetFlag(FLAG_V) {
		cpu.PC = addr
		cycles++
	}
	return cycles
}

// 栈指令实现
func (cpu *CPU) PHA() int {
	cpu.Push(cpu.A)
	return 3
}

func (cpu *CPU) PHP() int {
	cpu.Push(cpu.Status | FLAG_B | FLAG_U)
	return 3
}

func (cpu *CPU) PLA() int {
	cpu.A = cpu.Pop()
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 4
}

func (cpu *CPU) PLP() int {
	cpu.Status = cpu.Pop()
	cpu.SetFlag(FLAG_U, true)
	return 4
}

// 寄存器操作指令实现
func (cpu *CPU) TAX() int {
	cpu.X = cpu.A
	cpu.SetZeroAndNegativeFlags(cpu.X)
	return 2
}

func (cpu *CPU) TXA() int {
	cpu.A = cpu.X
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 2
}

func (cpu *CPU) TAY() int {
	cpu.Y = cpu.A
	cpu.SetZeroAndNegativeFlags(cpu.Y)
	return 2
}

func (cpu *CPU) TYA() int {
	cpu.A = cpu.Y
	cpu.SetZeroAndNegativeFlags(cpu.A)
	return 2
}

func (cpu *CPU) TSX() int {
	cpu.X = cpu.SP
	cpu.SetZeroAndNegativeFlags(cpu.X)
	return 2
}

func (cpu *CPU) TXS() int {
	cpu.SP = cpu.X
	return 2
}

// 标志位操作指令实现
func (cpu *CPU) CLC() int {
	cpu.SetFlag(FLAG_C, false)
	return 2
}

func (cpu *CPU) SEC() int {
	cpu.SetFlag(FLAG_C, true)
	return 2
}

func (cpu *CPU) CLI() int {
	cpu.SetFlag(FLAG_I, false)
	return 2
}

func (cpu *CPU) SEI() int {
	cpu.SetFlag(FLAG_I, true)
	return 2
}

func (cpu *CPU) CLV() int {
	cpu.SetFlag(FLAG_V, false)
	return 2
}

func (cpu *CPU) CLD() int {
	cpu.SetFlag(FLAG_D, false)
	return 2
}

func (cpu *CPU) SED() int {
	cpu.SetFlag(FLAG_D, true)
	return 2
}

// 递增递减指令实现
func (cpu *CPU) INC() int {
	addr, cycles := cpu.GetOperandAddress(ZP)
	value := cpu.ReadByte(addr)
	value++
	cpu.WriteByte(addr, value)
	cpu.SetZeroAndNegativeFlags(value)
	return cycles
}

func (cpu *CPU) INX() int {
	cpu.X++
	cpu.SetZeroAndNegativeFlags(cpu.X)
	return 2
}

func (cpu *CPU) INY() int {
	cpu.Y++
	cpu.SetZeroAndNegativeFlags(cpu.Y)
	return 2
}

func (cpu *CPU) DEC() int {
	addr, cycles := cpu.GetOperandAddress(ZP)
	value := cpu.ReadByte(addr)
	value--
	cpu.WriteByte(addr, value)
	cpu.SetZeroAndNegativeFlags(value)
	return cycles
}

func (cpu *CPU) DEX() int {
	cpu.X--
	cpu.SetZeroAndNegativeFlags(cpu.X)
	return 2
}

func (cpu *CPU) DEY() int {
	cpu.Y--
	cpu.SetZeroAndNegativeFlags(cpu.Y)
	return 2
}

// 跳转指令实现
func (cpu *CPU) JMP() int {
	addr, cycles := cpu.GetOperandAddress(ABS)
	cpu.PC = addr
	return cycles
}

func (cpu *CPU) JSR() int {
	addr, cycles := cpu.GetOperandAddress(ABS)
	cpu.PushWord(cpu.PC - 1)
	cpu.PC = addr
	return cycles
}

func (cpu *CPU) RTS() int {
	cpu.PC = cpu.PopWord() + 1
	return 6
}

// 中断指令实现
func (cpu *CPU) BRK() int {
	cpu.PushWord(cpu.PC + 1)
	cpu.Push(cpu.Status | FLAG_B | FLAG_U)
	cpu.SetFlag(FLAG_I, true)
	cpu.PC = cpu.ReadWord(0xFFFE)
	return 7
}

func (cpu *CPU) RTI() int {
	cpu.Status = cpu.Pop()
	cpu.SetFlag(FLAG_U, true)
	cpu.PC = cpu.PopWord()
	return 6
}

// 空操作指令实现
func (cpu *CPU) NOP() int {
	return 2
}

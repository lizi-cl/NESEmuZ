package cpu

import (
	"fmt"
)

// CPU状态标志位
const (
	FLAG_C = 0x01 // Carry
	FLAG_Z = 0x02 // Zero
	FLAG_I = 0x04 // Interrupt Disable
	FLAG_D = 0x08 // Decimal Mode
	FLAG_B = 0x10 // Break Command
	FLAG_U = 0x20 // Unused
	FLAG_V = 0x40 // Overflow
	FLAG_N = 0x80 // Negative
)

// 寻址模式
type AddressingMode int

const (
	IMM  AddressingMode = iota // Immediate
	ZP                         // Zero Page
	ZPX                        // Zero Page X
	ZPY                        // Zero Page Y
	ABS                        // Absolute
	ABSX                       // Absolute X
	ABSY                       // Absolute Y
	IND                        // Indirect
	INDX                       // Indirect X
	INDY                       // Indirect Y
	REL                        // Relative
	ACC                        // Accumulator
	IMP                        // Implied
)

// 指令结构
type Instruction struct {
	Name    string
	Mode    AddressingMode
	Cycles  int
	Execute func(*CPU) int
}

// CPU结构体
type CPU struct {
	A, X, Y byte        // 寄存器
	SP      byte        // 栈指针
	PC      uint16      // 程序计数器
	Status  byte        // 状态寄存器
	Memory  [65536]byte // 内存
	Cycles  int         // 当前周期数
}

// 初始化CPU
func NewCPU() *CPU {
	cpu := &CPU{}
	cpu.Reset()
	return cpu
}

// 重置CPU
func (cpu *CPU) Reset() {
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.SP = 0xFD
	cpu.Status = FLAG_U | FLAG_I
	// 从复位向量读取起始地址，如果没有设置则使用默认值
	resetVector := cpu.ReadWord(0xFFFC)
	if resetVector == 0 {
		cpu.PC = 0x8000 // 默认起始地址
	} else {
		cpu.PC = resetVector
	}
	cpu.Cycles = 0
}

// 读取内存中的一个字节
func (cpu *CPU) ReadByte(addr uint16) byte {
	return cpu.Memory[addr]
}

// 写入内存中的一个字节
func (cpu *CPU) WriteByte(addr uint16, value byte) {
	cpu.Memory[addr] = value
}

// 读取内存中的一个字（小端序）
func (cpu *CPU) ReadWord(addr uint16) uint16 {
	return uint16(cpu.ReadByte(addr)) | uint16(cpu.ReadByte(addr+1))<<8
}

// 写入内存中的一个字（小端序）
func (cpu *CPU) WriteWord(addr uint16, value uint16) {
	cpu.WriteByte(addr, byte(value&0xFF))
	cpu.WriteByte(addr+1, byte(value>>8))
}

// 推入栈
func (cpu *CPU) Push(value byte) {
	cpu.WriteByte(0x0100+uint16(cpu.SP), value)
	cpu.SP--
}

// 从栈弹出
func (cpu *CPU) Pop() byte {
	cpu.SP++
	return cpu.ReadByte(0x0100 + uint16(cpu.SP))
}

// 推入字到栈
func (cpu *CPU) PushWord(value uint16) {
	cpu.Push(byte(value >> 8))
	cpu.Push(byte(value & 0xFF))
}

// 从栈弹出字
func (cpu *CPU) PopWord() uint16 {
	low := uint16(cpu.Pop())
	high := uint16(cpu.Pop())
	return high<<8 | low
}

// 设置状态标志
func (cpu *CPU) SetFlag(flag byte, set bool) {
	if set {
		cpu.Status |= flag
	} else {
		cpu.Status &^= flag
	}
}

// 获取状态标志
func (cpu *CPU) GetFlag(flag byte) bool {
	return (cpu.Status & flag) != 0
}

// 设置零标志和负标志
func (cpu *CPU) SetZeroAndNegativeFlags(value byte) {
	cpu.SetFlag(FLAG_Z, value == 0)
	cpu.SetFlag(FLAG_N, (value&0x80) != 0)
}

// 寻址模式实现
func (cpu *CPU) GetOperandAddress(mode AddressingMode) (uint16, int) {
	switch mode {
	case IMM:
		addr := cpu.PC
		cpu.PC++
		return addr, 1
	case ZP:
		addr := uint16(cpu.ReadByte(cpu.PC))
		cpu.PC++
		return addr, 2
	case ZPX:
		addr := uint16(cpu.ReadByte(cpu.PC) + cpu.X)
		cpu.PC++
		return addr & 0xFF, 3
	case ZPY:
		addr := uint16(cpu.ReadByte(cpu.PC) + cpu.Y)
		cpu.PC++
		return addr & 0xFF, 3
	case ABS:
		addr := cpu.ReadWord(cpu.PC)
		cpu.PC += 2
		return addr, 3
	case ABSX:
		base := cpu.ReadWord(cpu.PC)
		cpu.PC += 2
		addr := base + uint16(cpu.X)
		cycles := 3
		if (base & 0xFF00) != (addr & 0xFF00) {
			cycles++ // 页面交叉
		}
		return addr, cycles
	case ABSY:
		base := cpu.ReadWord(cpu.PC)
		cpu.PC += 2
		addr := base + uint16(cpu.Y)
		cycles := 3
		if (base & 0xFF00) != (addr & 0xFF00) {
			cycles++ // 页面交叉
		}
		return addr, cycles
	case IND:
		ptr := cpu.ReadWord(cpu.PC)
		cpu.PC += 2
		// 6502的间接寻址bug
		low := cpu.ReadByte(ptr)
		high := cpu.ReadByte((ptr & 0xFF00) | ((ptr + 1) & 0x00FF))
		return uint16(high)<<8 | uint16(low), 5
	case INDX:
		base := cpu.ReadByte(cpu.PC)
		cpu.PC++
		ptr := uint16((base + cpu.X) & 0xFF)
		low := cpu.ReadByte(ptr)
		high := cpu.ReadByte((ptr + 1) & 0xFF)
		return uint16(high)<<8 | uint16(low), 6
	case INDY:
		base := cpu.ReadByte(cpu.PC)
		cpu.PC++
		ptr := uint16(base)
		low := cpu.ReadByte(ptr)
		high := cpu.ReadByte((ptr + 1) & 0xFF)
		addr := uint16(high)<<8 | uint16(low) + uint16(cpu.Y)
		cycles := 5
		if (uint16(high)<<8|uint16(low))&0xFF00 != addr&0xFF00 {
			cycles++ // 页面交叉
		}
		return addr, cycles
	case REL:
		offset := int8(cpu.ReadByte(cpu.PC))
		cpu.PC++
		addr := uint16(int32(cpu.PC) + int32(offset))
		return addr, 2
	default:
		return 0, 1
	}
}

// 指令集
var instructions = map[byte]Instruction{
	// 加载指令
	0xA9: {"LDA", IMM, 2, func(cpu *CPU) int { return cpu.LDA() }},
	0xA5: {"LDA", ZP, 3, func(cpu *CPU) int { return cpu.LDA() }},
	0xB5: {"LDA", ZPX, 4, func(cpu *CPU) int { return cpu.LDA() }},
	0xAD: {"LDA", ABS, 4, func(cpu *CPU) int { return cpu.LDA() }},
	0xBD: {"LDA", ABSX, 4, func(cpu *CPU) int { return cpu.LDA() }},
	0xB9: {"LDA", ABSY, 4, func(cpu *CPU) int { return cpu.LDA() }},
	0xA1: {"LDA", INDX, 6, func(cpu *CPU) int { return cpu.LDA() }},
	0xB1: {"LDA", INDY, 5, func(cpu *CPU) int { return cpu.LDA() }},

	0xA2: {"LDX", IMM, 2, func(cpu *CPU) int { return cpu.LDX() }},
	0xA6: {"LDX", ZP, 3, func(cpu *CPU) int { return cpu.LDX() }},
	0xB6: {"LDX", ZPY, 4, func(cpu *CPU) int { return cpu.LDX() }},
	0xAE: {"LDX", ABS, 4, func(cpu *CPU) int { return cpu.LDX() }},
	0xBE: {"LDX", ABSY, 4, func(cpu *CPU) int { return cpu.LDX() }},

	0xA0: {"LDY", IMM, 2, func(cpu *CPU) int { return cpu.LDY() }},
	0xA4: {"LDY", ZP, 3, func(cpu *CPU) int { return cpu.LDY() }},
	0xB4: {"LDY", ZPX, 4, func(cpu *CPU) int { return cpu.LDY() }},
	0xAC: {"LDY", ABS, 4, func(cpu *CPU) int { return cpu.LDY() }},
	0xBC: {"LDY", ABSX, 4, func(cpu *CPU) int { return cpu.LDY() }},

	// 存储指令
	0x85: {"STA", ZP, 3, func(cpu *CPU) int { return cpu.STA() }},
	0x95: {"STA", ZPX, 4, func(cpu *CPU) int { return cpu.STA() }},
	0x8D: {"STA", ABS, 4, func(cpu *CPU) int { return cpu.STA() }},
	0x9D: {"STA", ABSX, 5, func(cpu *CPU) int { return cpu.STA() }},
	0x99: {"STA", ABSY, 5, func(cpu *CPU) int { return cpu.STA() }},
	0x81: {"STA", INDX, 6, func(cpu *CPU) int { return cpu.STA() }},
	0x91: {"STA", INDY, 6, func(cpu *CPU) int { return cpu.STA() }},

	0x86: {"STX", ZP, 3, func(cpu *CPU) int { return cpu.STX() }},
	0x96: {"STX", ZPY, 4, func(cpu *CPU) int { return cpu.STX() }},
	0x8E: {"STX", ABS, 4, func(cpu *CPU) int { return cpu.STX() }},

	0x84: {"STY", ZP, 3, func(cpu *CPU) int { return cpu.STY() }},
	0x94: {"STY", ZPX, 4, func(cpu *CPU) int { return cpu.STY() }},
	0x8C: {"STY", ABS, 4, func(cpu *CPU) int { return cpu.STY() }},

	// 算术指令
	0x69: {"ADC", IMM, 2, func(cpu *CPU) int { return cpu.ADC() }},
	0x65: {"ADC", ZP, 3, func(cpu *CPU) int { return cpu.ADC() }},
	0x75: {"ADC", ZPX, 4, func(cpu *CPU) int { return cpu.ADC() }},
	0x6D: {"ADC", ABS, 4, func(cpu *CPU) int { return cpu.ADC() }},
	0x7D: {"ADC", ABSX, 4, func(cpu *CPU) int { return cpu.ADC() }},
	0x79: {"ADC", ABSY, 4, func(cpu *CPU) int { return cpu.ADC() }},
	0x61: {"ADC", INDX, 6, func(cpu *CPU) int { return cpu.ADC() }},
	0x71: {"ADC", INDY, 5, func(cpu *CPU) int { return cpu.ADC() }},

	0xE9: {"SBC", IMM, 2, func(cpu *CPU) int { return cpu.SBC() }},
	0xE5: {"SBC", ZP, 3, func(cpu *CPU) int { return cpu.SBC() }},
	0xF5: {"SBC", ZPX, 4, func(cpu *CPU) int { return cpu.SBC() }},
	0xED: {"SBC", ABS, 4, func(cpu *CPU) int { return cpu.SBC() }},
	0xFD: {"SBC", ABSX, 4, func(cpu *CPU) int { return cpu.SBC() }},
	0xF9: {"SBC", ABSY, 4, func(cpu *CPU) int { return cpu.SBC() }},
	0xE1: {"SBC", INDX, 6, func(cpu *CPU) int { return cpu.SBC() }},
	0xF1: {"SBC", INDY, 5, func(cpu *CPU) int { return cpu.SBC() }},

	// 逻辑指令
	0x29: {"AND", IMM, 2, func(cpu *CPU) int { return cpu.AND() }},
	0x25: {"AND", ZP, 3, func(cpu *CPU) int { return cpu.AND() }},
	0x35: {"AND", ZPX, 4, func(cpu *CPU) int { return cpu.AND() }},
	0x2D: {"AND", ABS, 4, func(cpu *CPU) int { return cpu.AND() }},
	0x3D: {"AND", ABSX, 4, func(cpu *CPU) int { return cpu.AND() }},
	0x39: {"AND", ABSY, 4, func(cpu *CPU) int { return cpu.AND() }},
	0x21: {"AND", INDX, 6, func(cpu *CPU) int { return cpu.AND() }},
	0x31: {"AND", INDY, 5, func(cpu *CPU) int { return cpu.AND() }},

	0x49: {"EOR", IMM, 2, func(cpu *CPU) int { return cpu.EOR() }},
	0x45: {"EOR", ZP, 3, func(cpu *CPU) int { return cpu.EOR() }},
	0x55: {"EOR", ZPX, 4, func(cpu *CPU) int { return cpu.EOR() }},
	0x4D: {"EOR", ABS, 4, func(cpu *CPU) int { return cpu.EOR() }},
	0x5D: {"EOR", ABSX, 4, func(cpu *CPU) int { return cpu.EOR() }},
	0x59: {"EOR", ABSY, 4, func(cpu *CPU) int { return cpu.EOR() }},
	0x41: {"EOR", INDX, 6, func(cpu *CPU) int { return cpu.EOR() }},
	0x51: {"EOR", INDY, 5, func(cpu *CPU) int { return cpu.EOR() }},

	0x09: {"ORA", IMM, 2, func(cpu *CPU) int { return cpu.ORA() }},
	0x05: {"ORA", ZP, 3, func(cpu *CPU) int { return cpu.ORA() }},
	0x15: {"ORA", ZPX, 4, func(cpu *CPU) int { return cpu.ORA() }},
	0x0D: {"ORA", ABS, 4, func(cpu *CPU) int { return cpu.ORA() }},
	0x1D: {"ORA", ABSX, 4, func(cpu *CPU) int { return cpu.ORA() }},
	0x19: {"ORA", ABSY, 4, func(cpu *CPU) int { return cpu.ORA() }},
	0x01: {"ORA", INDX, 6, func(cpu *CPU) int { return cpu.ORA() }},
	0x11: {"ORA", INDY, 5, func(cpu *CPU) int { return cpu.ORA() }},

	// 移位指令
	0x0A: {"ASL", ACC, 2, func(cpu *CPU) int { return cpu.ASL() }},
	0x06: {"ASL", ZP, 5, func(cpu *CPU) int { return cpu.ASL() }},
	0x16: {"ASL", ZPX, 6, func(cpu *CPU) int { return cpu.ASL() }},
	0x0E: {"ASL", ABS, 6, func(cpu *CPU) int { return cpu.ASL() }},
	0x1E: {"ASL", ABSX, 7, func(cpu *CPU) int { return cpu.ASL() }},

	0x4A: {"LSR", ACC, 2, func(cpu *CPU) int { return cpu.LSR() }},
	0x46: {"LSR", ZP, 5, func(cpu *CPU) int { return cpu.LSR() }},
	0x56: {"LSR", ZPX, 6, func(cpu *CPU) int { return cpu.LSR() }},
	0x4E: {"LSR", ABS, 6, func(cpu *CPU) int { return cpu.LSR() }},
	0x5E: {"LSR", ABSX, 7, func(cpu *CPU) int { return cpu.LSR() }},

	0x2A: {"ROL", ACC, 2, func(cpu *CPU) int { return cpu.ROL() }},
	0x26: {"ROL", ZP, 5, func(cpu *CPU) int { return cpu.ROL() }},
	0x36: {"ROL", ZPX, 6, func(cpu *CPU) int { return cpu.ROL() }},
	0x2E: {"ROL", ABS, 6, func(cpu *CPU) int { return cpu.ROL() }},
	0x3E: {"ROL", ABSX, 7, func(cpu *CPU) int { return cpu.ROL() }},

	0x6A: {"ROR", ACC, 2, func(cpu *CPU) int { return cpu.ROR() }},
	0x66: {"ROR", ZP, 5, func(cpu *CPU) int { return cpu.ROR() }},
	0x76: {"ROR", ZPX, 6, func(cpu *CPU) int { return cpu.ROR() }},
	0x6E: {"ROR", ABS, 6, func(cpu *CPU) int { return cpu.ROR() }},
	0x7E: {"ROR", ABSX, 7, func(cpu *CPU) int { return cpu.ROR() }},

	// 比较指令
	0xC9: {"CMP", IMM, 2, func(cpu *CPU) int { return cpu.CMP() }},
	0xC5: {"CMP", ZP, 3, func(cpu *CPU) int { return cpu.CMP() }},
	0xD5: {"CMP", ZPX, 4, func(cpu *CPU) int { return cpu.CMP() }},
	0xCD: {"CMP", ABS, 4, func(cpu *CPU) int { return cpu.CMP() }},
	0xDD: {"CMP", ABSX, 4, func(cpu *CPU) int { return cpu.CMP() }},
	0xD9: {"CMP", ABSY, 4, func(cpu *CPU) int { return cpu.CMP() }},
	0xC1: {"CMP", INDX, 6, func(cpu *CPU) int { return cpu.CMP() }},
	0xD1: {"CMP", INDY, 5, func(cpu *CPU) int { return cpu.CMP() }},

	0xE0: {"CPX", IMM, 2, func(cpu *CPU) int { return cpu.CPX() }},
	0xE4: {"CPX", ZP, 3, func(cpu *CPU) int { return cpu.CPX() }},
	0xEC: {"CPX", ABS, 4, func(cpu *CPU) int { return cpu.CPX() }},

	0xC0: {"CPY", IMM, 2, func(cpu *CPU) int { return cpu.CPY() }},
	0xC4: {"CPY", ZP, 3, func(cpu *CPU) int { return cpu.CPY() }},
	0xCC: {"CPY", ABS, 4, func(cpu *CPU) int { return cpu.CPY() }},

	// 分支指令
	0x90: {"BCC", REL, 2, func(cpu *CPU) int { return cpu.BCC() }},
	0xB0: {"BCS", REL, 2, func(cpu *CPU) int { return cpu.BCS() }},
	0xF0: {"BEQ", REL, 2, func(cpu *CPU) int { return cpu.BEQ() }},
	0x30: {"BMI", REL, 2, func(cpu *CPU) int { return cpu.BMI() }},
	0xD0: {"BNE", REL, 2, func(cpu *CPU) int { return cpu.BNE() }},
	0x10: {"BPL", REL, 2, func(cpu *CPU) int { return cpu.BPL() }},
	0x50: {"BVC", REL, 2, func(cpu *CPU) int { return cpu.BVC() }},
	0x70: {"BVS", REL, 2, func(cpu *CPU) int { return cpu.BVS() }},

	// 栈指令
	0x48: {"PHA", IMP, 3, func(cpu *CPU) int { return cpu.PHA() }},
	0x08: {"PHP", IMP, 3, func(cpu *CPU) int { return cpu.PHP() }},
	0x68: {"PLA", IMP, 4, func(cpu *CPU) int { return cpu.PLA() }},
	0x28: {"PLP", IMP, 4, func(cpu *CPU) int { return cpu.PLP() }},

	// 寄存器操作
	0xAA: {"TAX", IMP, 2, func(cpu *CPU) int { return cpu.TAX() }},
	0x8A: {"TXA", IMP, 2, func(cpu *CPU) int { return cpu.TXA() }},
	0xA8: {"TAY", IMP, 2, func(cpu *CPU) int { return cpu.TAY() }},
	0x98: {"TYA", IMP, 2, func(cpu *CPU) int { return cpu.TYA() }},
	0xBA: {"TSX", IMP, 2, func(cpu *CPU) int { return cpu.TSX() }},
	0x9A: {"TXS", IMP, 2, func(cpu *CPU) int { return cpu.TXS() }},

	// 标志位操作
	0x18: {"CLC", IMP, 2, func(cpu *CPU) int { return cpu.CLC() }},
	0x38: {"SEC", IMP, 2, func(cpu *CPU) int { return cpu.SEC() }},
	0x58: {"CLI", IMP, 2, func(cpu *CPU) int { return cpu.CLI() }},
	0x78: {"SEI", IMP, 2, func(cpu *CPU) int { return cpu.SEI() }},
	0xB8: {"CLV", IMP, 2, func(cpu *CPU) int { return cpu.CLV() }},
	0xD8: {"CLD", IMP, 2, func(cpu *CPU) int { return cpu.CLD() }},
	0xF8: {"SED", IMP, 2, func(cpu *CPU) int { return cpu.SED() }},

	// 递增递减
	0xE6: {"INC", ZP, 5, func(cpu *CPU) int { return cpu.INC() }},
	0xF6: {"INC", ZPX, 6, func(cpu *CPU) int { return cpu.INC() }},
	0xEE: {"INC", ABS, 6, func(cpu *CPU) int { return cpu.INC() }},
	0xFE: {"INC", ABSX, 7, func(cpu *CPU) int { return cpu.INC() }},

	0xE8: {"INX", IMP, 2, func(cpu *CPU) int { return cpu.INX() }},
	0xC8: {"INY", IMP, 2, func(cpu *CPU) int { return cpu.INY() }},

	0xC6: {"DEC", ZP, 5, func(cpu *CPU) int { return cpu.DEC() }},
	0xD6: {"DEC", ZPX, 6, func(cpu *CPU) int { return cpu.DEC() }},
	0xCE: {"DEC", ABS, 6, func(cpu *CPU) int { return cpu.DEC() }},
	0xDE: {"DEC", ABSX, 7, func(cpu *CPU) int { return cpu.DEC() }},

	0xCA: {"DEX", IMP, 2, func(cpu *CPU) int { return cpu.DEX() }},
	0x88: {"DEY", IMP, 2, func(cpu *CPU) int { return cpu.DEY() }},

	// 跳转指令
	0x4C: {"JMP", ABS, 3, func(cpu *CPU) int { return cpu.JMP() }},
	0x6C: {"JMP", IND, 5, func(cpu *CPU) int { return cpu.JMP() }},
	0x20: {"JSR", ABS, 6, func(cpu *CPU) int { return cpu.JSR() }},
	0x60: {"RTS", IMP, 6, func(cpu *CPU) int { return cpu.RTS() }},

	// 中断指令
	0x00: {"BRK", IMP, 7, func(cpu *CPU) int { return cpu.BRK() }},
	0x40: {"RTI", IMP, 6, func(cpu *CPU) int { return cpu.RTI() }},

	// 空操作
	0xEA: {"NOP", IMP, 2, func(cpu *CPU) int { return cpu.NOP() }},
}

// 执行指令
func (cpu *CPU) Execute() int {
	opcode := cpu.ReadByte(cpu.PC)
	cpu.PC++

	instruction, exists := instructions[opcode]
	if !exists {
		fmt.Printf("未知指令: %02X at PC=%04X\n", opcode, cpu.PC-1)
		return 1
	}

	// 根据寻址模式执行指令
	switch instruction.Mode {
	case IMM:
		// 立即寻址 - 操作数在指令中
		value := cpu.ReadByte(cpu.PC)
		cpu.PC++
		return cpu.executeImmediate(instruction, value)
	case ZP, ZPX, ZPY, ABS, ABSX, ABSY, IND, INDX, INDY:
		// 内存寻址 - 需要计算地址
		addr, cycles := cpu.GetOperandAddress(instruction.Mode)
		return cpu.executeMemory(instruction, addr, cycles)
	case REL:
		// 相对寻址 - 用于分支指令
		offset := int8(cpu.ReadByte(cpu.PC))
		cpu.PC++
		return cpu.executeRelative(instruction, offset)
	case ACC, IMP:
		// 累加器或隐含寻址
		return instruction.Execute(cpu)
	default:
		return 1
	}
}

// 执行立即寻址指令
func (cpu *CPU) executeImmediate(instruction Instruction, value byte) int {
	switch instruction.Name {
	case "LDA":
		cpu.A = value
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return 2
	case "LDX":
		cpu.X = value
		cpu.SetZeroAndNegativeFlags(cpu.X)
		return 2
	case "LDY":
		cpu.Y = value
		cpu.SetZeroAndNegativeFlags(cpu.Y)
		return 2
	case "ADC":
		return cpu.ADC()
	case "SBC":
		return cpu.SBC()
	case "AND":
		return cpu.AND()
	case "EOR":
		return cpu.EOR()
	case "ORA":
		return cpu.ORA()
	case "CMP":
		return cpu.CMP()
	case "CPX":
		return cpu.CPX()
	case "CPY":
		return cpu.CPY()
	default:
		return 2
	}
}

// 执行内存寻址指令
func (cpu *CPU) executeMemory(instruction Instruction, addr uint16, cycles int) int {
	switch instruction.Name {
	case "LDA":
		value := cpu.ReadByte(addr)
		cpu.A = value
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return cycles
	case "LDX":
		value := cpu.ReadByte(addr)
		cpu.X = value
		cpu.SetZeroAndNegativeFlags(cpu.X)
		return cycles
	case "LDY":
		value := cpu.ReadByte(addr)
		cpu.Y = value
		cpu.SetZeroAndNegativeFlags(cpu.Y)
		return cycles
	case "STA":
		cpu.WriteByte(addr, cpu.A)
		return cycles
	case "STX":
		cpu.WriteByte(addr, cpu.X)
		return cycles
	case "STY":
		cpu.WriteByte(addr, cpu.Y)
		return cycles
	case "ADC":
		value := cpu.ReadByte(addr)
		carry := uint16(0)
		if cpu.GetFlag(FLAG_C) {
			carry = 1
		}
		result := uint16(cpu.A) + uint16(value) + carry
		cpu.SetFlag(FLAG_C, result > 0xFF)
		overflow := ((cpu.A^value)&0x80 == 0) && ((cpu.A^byte(result))&0x80 != 0)
		cpu.SetFlag(FLAG_V, overflow)
		cpu.A = byte(result)
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return cycles
	case "SBC":
		value := cpu.ReadByte(addr)
		carry := uint16(0)
		if cpu.GetFlag(FLAG_C) {
			carry = 1
		}
		result := uint16(cpu.A) - uint16(value) - (1 - carry)
		cpu.SetFlag(FLAG_C, result <= 0xFF)
		overflow := ((cpu.A^value)&0x80 != 0) && ((cpu.A^byte(result))&0x80 != 0)
		cpu.SetFlag(FLAG_V, overflow)
		cpu.A = byte(result)
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return cycles
	case "AND":
		value := cpu.ReadByte(addr)
		cpu.A &= value
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return cycles
	case "EOR":
		value := cpu.ReadByte(addr)
		cpu.A ^= value
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return cycles
	case "ORA":
		value := cpu.ReadByte(addr)
		cpu.A |= value
		cpu.SetZeroAndNegativeFlags(cpu.A)
		return cycles
	case "CMP":
		value := cpu.ReadByte(addr)
		result := cpu.A - value
		cpu.SetFlag(FLAG_C, cpu.A >= value)
		cpu.SetZeroAndNegativeFlags(byte(result))
		return cycles
	case "INC":
		value := cpu.ReadByte(addr)
		value++
		cpu.WriteByte(addr, value)
		cpu.SetZeroAndNegativeFlags(value)
		return cycles
	case "DEC":
		value := cpu.ReadByte(addr)
		value--
		cpu.WriteByte(addr, value)
		cpu.SetZeroAndNegativeFlags(value)
		return cycles
	case "ASL":
		value := cpu.ReadByte(addr)
		carry := (value & 0x80) != 0
		value <<= 1
		cpu.WriteByte(addr, value)
		cpu.SetFlag(FLAG_C, carry)
		cpu.SetZeroAndNegativeFlags(value)
		return cycles
	case "LSR":
		value := cpu.ReadByte(addr)
		carry := (value & 0x01) != 0
		value >>= 1
		cpu.WriteByte(addr, value)
		cpu.SetFlag(FLAG_C, carry)
		cpu.SetZeroAndNegativeFlags(value)
		return cycles
	case "ROL":
		value := cpu.ReadByte(addr)
		oldCarry := cpu.GetFlag(FLAG_C)
		carry := (value & 0x80) != 0
		value <<= 1
		if oldCarry {
			value |= 0x01
		}
		cpu.WriteByte(addr, value)
		cpu.SetFlag(FLAG_C, carry)
		cpu.SetZeroAndNegativeFlags(value)
		return cycles
	case "ROR":
		value := cpu.ReadByte(addr)
		oldCarry := cpu.GetFlag(FLAG_C)
		carry := (value & 0x01) != 0
		value >>= 1
		if oldCarry {
			value |= 0x80
		}
		cpu.WriteByte(addr, value)
		cpu.SetFlag(FLAG_C, carry)
		cpu.SetZeroAndNegativeFlags(value)
		return cycles
	case "JMP":
		cpu.PC = addr
		return cycles
	case "JSR":
		cpu.PushWord(cpu.PC - 1)
		cpu.PC = addr
		return cycles
	default:
		return cycles
	}
}

// 执行相对寻址指令
func (cpu *CPU) executeRelative(instruction Instruction, offset int8) int {
	cycles := 2
	switch instruction.Name {
	case "BCC":
		if !cpu.GetFlag(FLAG_C) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BCS":
		if cpu.GetFlag(FLAG_C) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BEQ":
		if cpu.GetFlag(FLAG_Z) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BMI":
		if cpu.GetFlag(FLAG_N) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BNE":
		if !cpu.GetFlag(FLAG_Z) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BPL":
		if !cpu.GetFlag(FLAG_N) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BVC":
		if !cpu.GetFlag(FLAG_V) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	case "BVS":
		if cpu.GetFlag(FLAG_V) {
			cpu.PC = uint16(int32(cpu.PC) + int32(offset))
			cycles++
		}
	}
	return cycles
}

// 执行一个周期
func (cpu *CPU) Step() int {
	cycles := cpu.Execute()
	cpu.Cycles += cycles
	return cycles
}

// 获取CPU状态字符串
func (cpu *CPU) GetStatusString() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X SP:%02X PC:%04X Status:%02X",
		cpu.A, cpu.X, cpu.Y, cpu.SP, cpu.PC, cpu.Status)
}

package cpu

import (
	"testing"
)

func TestLDAImmediate(t *testing.T) {
	cpu := NewCPU()

	// 加载 LDA Immediate 指令和操作数到内存中
	cpu.Memory[0x8000] = 0xA9 // LDA Immediate
	cpu.Memory[0x8001] = 0x42 // 操作数

	// 执行指令
	cpu.Step()

	// 检查 A 寄存器的值是否正确
	if cpu.A != 0x42 {
		t.Errorf("期望 A 寄存器的值为 0x42，但得到 %02X", cpu.A)
	}

	// 检查零标志
	if cpu.Status&0x02 != 0 {
		t.Errorf("期望零标志未设置，但它被设置了")
	}

	// 检查负标志
	if cpu.Status&0x80 != 0 {
		t.Errorf("期望负标志未设置，但它被设置了")
	}
}

func TestLDAImmediateZeroFlag(t *testing.T) {
	cpu := NewCPU()

	// 加载 LDA Immediate 指令和操作数到内存中
	cpu.Memory[0x8000] = 0xA9 // LDA Immediate
	cpu.Memory[0x8001] = 0x00 // 操作数

	// 执行指令
	cpu.Step()

	// 检查 A 寄存器的值是否正确
	if cpu.A != 0x00 {
		t.Errorf("期望 A 寄存器的值为 0x00，但得到 %02X", cpu.A)
	}

	// 检查零标志
	if cpu.Status&0x02 == 0 {
		t.Errorf("期望零标志被设置，但它未被设置")
	}

	// 检查负标志
	if cpu.Status&0x80 != 0 {
		t.Errorf("期望负标志未设置，但它被设置了")
	}
}

func TestLDAImmediateNegativeFlag(t *testing.T) {
	cpu := NewCPU()

	// 加载 LDA Immediate 指令和操作数到内存中
	cpu.Memory[0x8000] = 0xA9 // LDA Immediate
	cpu.Memory[0x8001] = 0x80 // 操作数

	// 执行指令
	cpu.Step()

	// 检查 A 寄存器的值是否正确
	if cpu.A != 0x80 {
		t.Errorf("期望 A 寄存器的值为 0x80，但得到 %02X", cpu.A)
	}

	// 检查零标志
	if cpu.Status&0x02 != 0 {
		t.Errorf("期望零标志未设置，但它被设置了")
	}

	// 检查负标志
	if cpu.Status&0x80 == 0 {
		t.Errorf("期望负标志被设置，但它未被设置")
	}
}

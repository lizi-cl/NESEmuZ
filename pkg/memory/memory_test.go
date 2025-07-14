package memory

import (
	"testing"
)

// MockPPU 用于测试的PPU实现

type MockPPU struct {
	registers [8]byte
}

func (m *MockPPU) WriteRegister(addr uint16, value byte) {
	m.registers[addr] = value
}

func (m *MockPPU) ReadRegister(addr uint16) byte {
	return m.registers[addr]
}

// MockAPU 用于测试的APU实现

type MockAPU struct {
	registers [24]byte
}

func (m *MockAPU) WriteRegister(addr uint16, value byte) {
	m.registers[addr-0x4000] = value
}

func (m *MockAPU) ReadRegister(addr uint16) byte {
	return m.registers[addr-0x4000]
}

// MockInputs 用于测试的手柄实现

type MockInputs struct {
	registers [2]byte
}

func (m *MockInputs) WriteRegister(addr uint16, value byte) {
	m.registers[addr-0x4016] = value
}

func (m *MockInputs) ReadRegister(addr uint16) byte {
	return m.registers[addr-0x4016]
}

func TestMemory(t *testing.T) {
	ppu := &MockPPU{}
	apu := &MockAPU{}
	inputs := &MockInputs{}
	memory := &Memory{
		PPU:    ppu,
		APU:    apu,
		Inputs: inputs,
	}

	// 测试RAM写入和读取
	memory.Write(0x0000, 0x42)
	if memory.Read(0x0000) != 0x42 {
		t.Errorf("RAM read/write failed")
	}

	// 测试PPU寄存器写入和读取
	memory.Write(0x2002, 0x84)
	if memory.Read(0x2002) != 0x84 {
		t.Errorf("PPU register read/write failed")
	}

	// 测试APU寄存器写入和读取
	memory.Write(0x4004, 0x99)
	if memory.Read(0x4004) != 0x99 {
		t.Errorf("APU register read/write failed")
	}

	// 测试手柄寄存器写入和读取
	memory.Write(0x4016, 0x77)
	if memory.Read(0x4016) != 0x77 {
		t.Errorf("Inputs register read/write failed")
	}
}

package memory

import (
	"log"
	"os"
)

// NES内存总线大小
const (
	MemorySize = 0x10000 // 64KB
	RamSize    = 0x0800  // 2KB内部RAM
)

// Memory 表示NES主机的内存总线
// 可扩展为带有映射、IO、PPU、APU等
// 这里只实现最基础的RAM和ROM映射

type Memory struct {
	RAM    [RamSize]byte
	ROM    []byte       // PRG ROM
	SRAM   [0x2000]byte // 8KB SRAM
	PPU    PPUInterface
	APU    APUInterface
	Inputs InputInterface
}

// NewMemory 创建新的内存对象
func NewMemory(prgROM []byte) *Memory {
	return &Memory{
		ROM: prgROM,
	}
}

// 定义PPU、APU和手柄模块的接口

type PPUInterface interface {
	WriteRegister(addr uint16, value byte)
	ReadRegister(addr uint16) byte
}

type APUInterface interface {
	WriteRegister(addr uint16, value byte)
	ReadRegister(addr uint16) byte
}

type InputInterface interface {
	WriteRegister(addr uint16, value byte)
	ReadRegister(addr uint16) byte
}

// Read 读取内存
func (m *Memory) Read(addr uint16) byte {
	switch {
	case addr < 0x2000:
		// 2KB RAM镜像
		return m.RAM[addr%RamSize]
	case addr >= 0x2000 && addr < 0x4000:
		// PPU寄存器（8字节镜像）
		if m.PPU != nil {
			return m.PPU.ReadRegister(addr % 8)
		}
	case addr >= 0x4000 && addr < 0x4018:
		// APU和IO寄存器
		if addr >= 0x4016 && addr <= 0x4017 && m.Inputs != nil {
			return m.Inputs.ReadRegister(addr)
		} else if m.APU != nil {
			return m.APU.ReadRegister(addr)
		}
	case addr >= 0x6000 && addr < 0x8000:
		// SRAM
		return m.SRAM[addr-0x6000]
	case addr >= 0x8000:
		// PRG ROM区
		if len(m.ROM) == 0 {
			return 0
		}
		romAddr := int(addr - 0x8000)
		if len(m.ROM) == 0x4000 && romAddr >= 0x4000 {
			// 16KB ROM镜像
			romAddr = romAddr % 0x4000
		}
		if romAddr >= 0 && romAddr < len(m.ROM) {
			return m.ROM[romAddr]
		}
		return 0
	default:
		// 其他区域暂未实现
		return 0
	}
	return 0
}

// Write 写内存
func (m *Memory) Write(addr uint16, value byte) {
	switch {
	case addr < 0x2000:
		// 2KB RAM镜像
		m.RAM[addr%RamSize] = value
	case addr >= 0x2000 && addr < 0x4000:
		// PPU寄存器（8字节镜像）
		if m.PPU != nil {
			m.PPU.WriteRegister(addr%8, value)
		}
	case addr >= 0x4000 && addr < 0x4018:
		// APU和IO寄存器
		if addr >= 0x4016 && addr <= 0x4017 && m.Inputs != nil {
			m.Inputs.WriteRegister(addr, value)
		} else if m.APU != nil {
			m.APU.WriteRegister(addr, value)
		}
	case addr >= 0x6000 && addr < 0x8000:
		// SRAM
		m.SRAM[addr-0x6000] = value
	}
	// 0x8000及以上为ROM，通常不可写
}

func (m *Memory) LoadSRAM(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("无法打开存档文件: %v", err)
		return err
	}
	defer file.Close()

	buffer := make([]byte, len(m.SRAM))
	_, err = file.Read(buffer)
	if err != nil {
		log.Printf("读取存档文件失败: %v", err)
		return err
	}

	copy(m.SRAM[:], buffer)
	return nil
}

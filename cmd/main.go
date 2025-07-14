package main

import (
	"fmt"

	"com.lizi.emu/nesemuz/pkg/cpu"
)

func main() {
	// 创建CPU实例
	cpu := cpu.NewCPU()

	// 设置复位向量
	cpu.WriteWord(0xFFFC, 0x8000)

	// 加载一些测试程序到内存
	// LDA #$42 (加载立即数0x42到A寄存器)
	cpu.WriteByte(0x8000, 0xA9)
	cpu.WriteByte(0x8001, 0x42)

	// STA $00 (将A寄存器的值存储到零页地址0x00)
	cpu.WriteByte(0x8002, 0x85)
	cpu.WriteByte(0x8003, 0x00)

	// LDX #$10 (加载立即数0x10到X寄存器)
	cpu.WriteByte(0x8004, 0xA2)
	cpu.WriteByte(0x8005, 0x10)

	// LDY #$20 (加载立即数0x20到Y寄存器)
	cpu.WriteByte(0x8006, 0xA0)
	cpu.WriteByte(0x8007, 0x20)

	// ADC #$01 (A寄存器加1)
	cpu.WriteByte(0x8008, 0x69)
	cpu.WriteByte(0x8009, 0x01)

	// 执行指令
	fmt.Println("初始状态:", cpu.GetStatusString())

	for i := 0; i < 5; i++ {
		cycles := cpu.Step()
		fmt.Printf("执行指令 %d, 周期: %d, 状态: %s\n", i+1, cycles, cpu.GetStatusString())
	}

	// 检查内存中的值
	fmt.Printf("内存地址0x00的值: %02X\n", cpu.ReadByte(0x00))
}

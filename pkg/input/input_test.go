package input

import (
	"testing"
)

func TestControllerBasic(t *testing.T) {
	c := NewController()
	// 初始状态全部未按下
	for i := 0; i < NumButtons; i++ {
		if c.buttons[i] {
			t.Errorf("Button %d should be unpressed initially", i)
		}
	}

	// 按下A和Start
	c.SetButton(ButtonA, true)
	c.SetButton(ButtonStart, true)
	if !c.buttons[ButtonA] || !c.buttons[ButtonStart] {
		t.Error("ButtonA and ButtonStart should be pressed")
	}
	if c.buttons[ButtonB] {
		t.Error("ButtonB should not be pressed")
	}
}

func TestControllerStrobeAndRead(t *testing.T) {
	c := NewController()
	// 按下A, B, Up
	c.SetButton(ButtonA, true)
	c.SetButton(ButtonB, true)
	c.SetButton(ButtonUp, true)
	// 写入strobe=1，移位寄存器应加载
	c.WriteStrobe(1)
	// 停止strobe模式
	c.WriteStrobe(0)
	// 读取8次，检查bit顺序
	expected := []byte{1, 1, 0, 0, 1, 0, 0, 0} // A, B, Select, Start, Up, Down, Left, Right
	for i := 0; i < NumButtons; i++ {
		v := c.Read()
		if v != expected[i] {
			t.Errorf("Read bit %d: got %d, want %d", i, v, expected[i])
		}
	}
	// 继续读，应该返回1（未连接）
	for i := 0; i < 4; i++ {
		v := c.Read()
		if v != 1 {
			t.Errorf("Read overflow bit %d: got %d, want 1", i, v)
		}
	}
}

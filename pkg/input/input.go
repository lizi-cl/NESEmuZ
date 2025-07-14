package input

// NES手柄按键定义
const (
	ButtonA      = 0
	ButtonB      = 1
	ButtonSelect = 2
	ButtonStart  = 3
	ButtonUp     = 4
	ButtonDown   = 5
	ButtonLeft   = 6
	ButtonRight  = 7
	NumButtons   = 8
)

// Controller 表示一个NES手柄
type Controller struct {
	buttons   [NumButtons]bool // 当前每个按键的状态
	shiftReg  byte             // 8位移位寄存器，bit=1表示按下，0表示未按下
	strobe    bool             // 是否处于strobe模式
	readIndex int              // 读取到第几个bit
}

// NewController 创建一个新的手柄
func NewController() *Controller {
	return &Controller{}
}

// SetButton 设置某个按键的状态
func (c *Controller) SetButton(button int, pressed bool) {
	if button >= 0 && button < NumButtons {
		c.buttons[button] = pressed
	}
}

// WriteStrobe 写入strobe信号
func (c *Controller) WriteStrobe(value byte) {
	c.strobe = value&1 == 1
	if c.strobe {
		c.loadShiftRegister()
	}
}

// Read 读取当前手柄的一个bit
func (c *Controller) Read() byte {
	var result byte
	if c.readIndex < NumButtons {
		// 从移位寄存器读取当前bit
		if (c.shiftReg>>c.readIndex)&1 == 1 {
			result = 1
		} else {
			result = 0
		}
	} else {
		result = 1 // 超出部分恒为1
	}
	if !c.strobe {
		c.readIndex++
	}
	return result
}

// 加载当前按钮状态到移位寄存器
func (c *Controller) loadShiftRegister() {
	c.shiftReg = 0
	for i := 0; i < NumButtons; i++ {
		if c.buttons[i] {
			c.shiftReg |= (1 << i) // 按下为1
		}
		// 未按下的保持为0
	}
	c.readIndex = 0
}

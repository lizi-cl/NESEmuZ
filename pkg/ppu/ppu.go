package ppu

const (
	VRAMSize     = 0x800 // 2KB VRAM
	ScreenWidth  = 256
	ScreenHeight = 240
)

type Palette struct {
	Colors [64]byte // NES调色板，64种颜色
}

type PPU struct {
	VRAM        [VRAMSize]byte // 视频内存
	Registers   [8]byte        // PPU寄存器
	FrameBuffer [256][240]byte // 帧缓冲区，用于存储渲染的图像
}

// NewPPU 创建新的PPU对象
func NewPPU() *PPU {
	return &PPU{}
}

// WriteRegister 写入PPU寄存器
func (p *PPU) WriteRegister(addr uint16, value byte) {
	if addr < 8 {
		p.Registers[addr] = value
	}
}

// ReadRegister 读取PPU寄存器
func (p *PPU) ReadRegister(addr uint16) byte {
	if addr < 8 {
		return p.Registers[addr]
	}
	return 0
}

func (p *PPU) InitializePalette() Palette {
	return Palette{
		Colors: [64]byte{
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
			0x0F, 0x30, 0x33, 0x36, 0x0F, 0x30, 0x33, 0x36,
		},
	}
}

func (p *PPU) RenderBackground() {
	palette := p.InitializePalette()
	for tileY := 0; tileY < ScreenHeight/8; tileY++ {
		for tileX := 0; tileX < ScreenWidth/8; tileX++ {
			tileIndex := tileY*(ScreenWidth/8) + tileX
			tileData := p.VRAM[tileIndex]
			attributeIndex := tileIndex / 4
			attributeData := p.VRAM[attributeIndex]
			paletteIndex := (attributeData >> ((tileIndex % 4) * 2)) & 0x03

			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					frameX := tileX*8 + x
					frameY := tileY*8 + y
					color := palette.Colors[paletteIndex*4+tileData%4]
					p.FrameBuffer[frameX][frameY] = color
				}
			}
		}
	}
}

func (p *PPU) RenderSprites() {
	palette := p.InitializePalette()
	for i := 0; i < 64; i++ {
		spriteX := int(p.VRAM[i*4])
		spriteY := int(p.VRAM[i*4+1])
		spriteTile := p.VRAM[i*4+2]
		attributeData := p.VRAM[i*4+3]
		paletteIndex := attributeData & 0x03

		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				frameX := spriteX + x
				frameY := spriteY + y
				if frameX < ScreenWidth && frameY < ScreenHeight {
					color := palette.Colors[paletteIndex*4+spriteTile%4]
					p.FrameBuffer[frameX][frameY] = color
				}
			}
		}
	}
}

// RenderFrame 渲染一帧图像
func (p *PPU) RenderFrame() {
	p.RenderBackground()
	p.RenderSprites()
}

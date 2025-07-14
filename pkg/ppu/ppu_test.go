package ppu

import (
	"testing"
)

func TestRenderBackground(t *testing.T) {
	ppu := NewPPU()
	ppu.VRAM[0] = 1 // 设置Tile数据
	ppu.VRAM[1] = 2 // 设置Tile数据
	ppu.VRAM[2] = 3 // 设置Tile数据
	ppu.VRAM[3] = 4 // 设置Tile数据
	ppu.RenderBackground()

	if ppu.FrameBuffer[0][0] == 0 {
		t.Errorf("背景渲染失败，帧缓冲区未更新")
	}
}

func TestRenderSprites(t *testing.T) {
	ppu := NewPPU()
	ppu.VRAM[0] = 10 // 设置精灵X坐标
	ppu.VRAM[1] = 20 // 设置精灵Y坐标
	ppu.VRAM[2] = 5  // 设置精灵Tile数据
	ppu.VRAM[3] = 1  // 设置精灵属性数据
	ppu.RenderSprites()

	if ppu.FrameBuffer[10][20] == 0 {
		t.Errorf("精灵渲染失败，帧缓冲区未更新")
	}
}

func TestRenderFrame(t *testing.T) {
	ppu := NewPPU()
	ppu.VRAM[0] = 1  // 设置背景Tile数据
	ppu.VRAM[4] = 10 // 设置精灵X坐标
	ppu.VRAM[5] = 20 // 设置精灵Y坐标
	ppu.VRAM[6] = 5  // 设置精灵Tile数据
	ppu.VRAM[7] = 1  // 设置精灵属性数据
	ppu.RenderFrame()

	if ppu.FrameBuffer[0][0] == 0 || ppu.FrameBuffer[10][20] == 0 {
		t.Errorf("帧渲染失败，背景或精灵未正确渲染")
	}
}

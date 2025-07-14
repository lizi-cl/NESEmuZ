package cartrige

import (
	"errors"
	"os"
)

type Cartrige struct {
	PRG []byte // PRG ROM数据
	CHR []byte // CHR ROM数据
}

func LoadCartrige(filePath string) (*Cartrige, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("无法打开ROM文件")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, errors.New("无法获取ROM文件信息")
	}

	romData := make([]byte, fileInfo.Size())
	_, err = file.Read(romData)
	if err != nil {
		return nil, errors.New("读取ROM文件失败")
	}

	if len(romData) < 16 {
		return nil, errors.New("ROM文件格式错误")
	}

	prgSize := int(romData[4]) * 16 * 1024
	chrSize := int(romData[5]) * 8 * 1024

	prgData := romData[16 : 16+prgSize]
	chrData := romData[16+prgSize : 16+prgSize+chrSize]

	return &Cartrige{
		PRG: prgData,
		CHR: chrData,
	}, nil
}

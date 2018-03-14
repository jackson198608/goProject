package composite

import (
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type Composite struct {
	imgaePath     string
	watermarkPath string
}

func NewComposite(imgaePath string, watermarkPath string) *Composite {
	if imgaePath == "" || watermarkPath == "" {
		return nil
	}

	c := new(Composite)
	if c == nil {
		return nil
	}

	c.imgaePath = imgaePath
	c.watermarkPath = watermarkPath
	return c
}

func (c *Composite) Do() error {
	err := c.compositeImage(c.imgaePath, c.watermarkPath)
	if err == nil {
		logger.Info("[sucess] composite image path is ", c.imgaePath, " watermarkPath is ", c.watermarkPath)
		return nil
	}
	return err
}

func (c *Composite) compositeImage(filename string, watermarkPath string) error {
	dest := imagick.NewMagickWand()
	src := imagick.NewMagickWand()

	//背景图
	if err := dest.ReadImage(filename); err != nil {
		logger.Error("ReadImage ", filename, err)
		return err
	}

	//水印图
	if err := src.ReadImage(watermarkPath); err != nil {
		logger.Error("ReadImage ", filename, err)
		return err
	}
	//获取水印尺寸
	srcWidth := int(src.GetImageWidth())
	srcHeight := int(src.GetImageHeight())

	//水印位于背景图的位置
	destWidth := int(dest.GetImageWidth()) - srcWidth
	destHeight := int(dest.GetImageHeight()) - srcHeight
	// This does the src (overlay) over the dest (background)
	dest.CompositeImage(src, imagick.COMPOSITE_OP_OVER, destWidth, destHeight)
	dest.WriteImage(filename)
	return nil
}

package composite

import (
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strings"
)

type Composite struct {
	imgaePath     string
	watermarkPath string
	suffix        string
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

	c.parsePath()
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
	nw := imagick.NewMagickWand()

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

	if c.suffix == "gif" {
		dest = dest.CoalesceImages()

		for i := 0; i < int(dest.GetNumberImages()); i++ {
			dest.SetIteratorIndex(i)
			tw := dest.GetImage()

			tw.CompositeImage(src, imagick.COMPOSITE_OP_OVER, destWidth, destHeight)
			tw.WriteImage(filename)

			nw.AddImage(tw)
			tw.Destroy()
		}
		dest.ResetIterator()
		dest.Destroy()
		dest = nw.CompareImageLayers(imagick.IMAGE_LAYER_COMPARE_ANY)
		// -loop 0
		dest.SetOption("loop", "0")
		dest.WriteImages(filename, true)
	}else {
		dest.CompositeImage(src, imagick.COMPOSITE_OP_OVER, destWidth, destHeight)
		dest.WriteImage(filename)
	}
	return nil
}

func (c *Composite) parsePath() error {
	rawSlice := []byte(c.imgaePath)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(c.imgaePath, ".")
	c.suffix = string(rawSlice[lastIndex+1 : rawLen])

	return nil
}

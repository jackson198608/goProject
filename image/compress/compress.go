package compress

import (
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/gographics/imagick.v2/imagick"
	// "reflect"
	"os"
	"strconv"
	"strings"
)

type Compress struct {
	imgaePath string
	width     int
	height    int
	filename  string
	suffix    string
}

func NewCompress(imgaePath string, width int, height int) *Compress {
	if imgaePath == "" || width == 0 {
		return nil
	}

	c := new(Compress)
	if c == nil {
		return nil
	}

	c.imgaePath = imgaePath
	c.width = width
	c.height = height

	c.parsePath()
	return c
}

// exists returns whether the given file or directory exists or not
func (c *Compress) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (c *Compress) Do() (string, error) {
	path, err := c.resizeImage(c.imgaePath, c.width, c.height)

	if err == nil {
		status, err := c.exists(path)
		if status {
			logger.Info("[sucess] compress image path is ", c.imgaePath, " width is ", c.width, " height is ", c.height)
		} else {
			//如果压缩后，图片不存在，则再尝试5次压缩
			for i := 0; i < 5; i++ {
				path, err = c.resizeImage(c.imgaePath, c.width, c.height)
				if err == nil {
					status, err = c.exists(path)
					if status {
						logger.Info("[sucess] next ", i, " compress image path is ", c.imgaePath, " width is ", c.width, " height is ", c.height)
						break
					}
				}
			}
		}
	}
	return path, err
}

func (c *Compress) resizeImage(filename string, width int, height int) (string, error) {
	var err error
	var newimg string

	mw := imagick.NewMagickWand()

	err = mw.ReadImage(filename)

	if err != nil {
		logger.Error(err)
		return newimg, err
	}

	// Get original logo size
	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()

	hWidth := uint(width)
	if height == 0 {
		ratio := float64(originalHeight) / float64(originalWidth)
		height = int(float64(hWidth) * ratio)
	}
	hHeight := uint(height)

	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		logger.Error(err)
		return newimg, err
	}

	err = mw.SetImageCompressionQuality(80)
	if err != nil {
		logger.Error(err)
		return newimg, err
	}
	widthStr := strconv.Itoa(width)
	newimg = c.filename + "_" + widthStr + "." + c.suffix
	err = mw.WriteImage(newimg)
	if err != nil {
		logger.Error(err)
		return newimg, err
	}
	return newimg, nil
}

func (c *Compress) parsePath() error {
	rawSlice := []byte(c.imgaePath)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(c.imgaePath, ".")
	c.filename = string(rawSlice[0:lastIndex])
	c.suffix = string(rawSlice[lastIndex+1 : rawLen])

	return nil
}

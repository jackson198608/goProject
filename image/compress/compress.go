package compress

import (
	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"
	"gopkg.in/gographics/imagick.v2/imagick"
	"strconv"
	"strings"
)

type Compress struct {
	jobstr   string
	filename string
	suffix   string
	jsonData *JsonColumn
}

//json column
type JsonColumn struct {
	imgaePath string
	width     int
	height    int
}

func NewCompress(jobStr string) *Compress {
	if jobStr == "" {
		return nil
	}

	c := new(Compress)
	if c == nil {
		return nil
	}

	c.jobstr = jobStr
	jsonColumn, err := c.parseJson()
	if err != nil {
		return nil
	}
	c.jsonData = jsonColumn

	c.parsePath()
	return c
}

//change json colum to object private member
func (c *Compress) parseJson() (*JsonColumn, error) {
	var jsonC JsonColumn
	js, err := simplejson.NewJson([]byte(c.jobstr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.imgaePath, _ = js.Get("path").String()
	jsonC.width, _ = js.Get("width").Int()
	jsonC.height, _ = js.Get("height").Int()

	return &jsonC, nil
}

func (c *Compress) Do() error {
	filename := c.jsonData.imgaePath
	width := c.jsonData.width
	height := c.jsonData.height
	err := c.resizeImage(filename, width, height)
	if err == nil {
		logger.Info("[sucess] compress image path is ", filename, " width is ", width, " height is ", height)
	}
	return err
}

func (c *Compress) resizeImage(filename string, width int, height int) error {
	imagick.Initialize()
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()

	err = mw.ReadImage(filename)

	if err != nil {
		logger.Error(err)
		return err
	}
	hWidth := uint(width)
	hHeight := uint(height)

	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = mw.SetImageCompressionQuality(80)
	if err != nil {
		logger.Error(err)
		return err
	}
	widthStr := strconv.Itoa(width)
	heightStr := strconv.Itoa(height)
	newimg := c.filename + "_" + widthStr + "_" + heightStr + "." + c.suffix
	err = mw.WriteImage(newimg)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (c *Compress) parsePath() error {
	rawSlice := []byte(c.jsonData.imgaePath)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(c.jsonData.imgaePath, ".")
	c.filename = string(rawSlice[0:lastIndex])
	c.suffix = string(rawSlice[lastIndex+1 : rawLen])

	return nil

}

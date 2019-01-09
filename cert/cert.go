package cert

import (
	"fmt"
	"icmm2019cert/database"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	_ "strings"

	"github.com/globalsign/mgo/bson"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	backgroundWidth  = 960
	backgroundHeight = 720
	utf8FontFile     = "./font.ttf"
	utf8FontSize     = float64(45.0)
	spacing          = float64(1.5)
	dpi              = float64(72)
	ctx              = new(freetype.Context)
	ctx4Test         = new(freetype.Context)
	utf8Font         = new(truetype.Font)
	red              = color.RGBA{255, 0, 0, 255}
	blue             = color.RGBA{0, 0, 255, 255}
	white            = color.RGBA{255, 255, 255, 255}
	black            = color.RGBA{0, 0, 0, 255}
	certBG           interface{}
	colorBlue        *image.Uniform
	colorBlack       *image.Uniform
	colorWhite       *image.Uniform
	img4TxtWidthTest *image.RGBA
)

// Runner is a runner's result
type Runner struct {
	ID       bson.ObjectId `bson:"_id"`
	BibNO    int           `bson:"bibNumber"`
	FName    string        `bson:"firstname"`
	LName    string        `bson:"lastname"`
	ChipTime float64       `bson:"chiptime"`
	GunTime  float64       `bson:"guntime"`
}

// func CertImgDebug(bibNO, txt, x, y string, w http.ResponseWriter) {

//Image gen cert img from bib
func Image(bibNO int, w http.ResponseWriter) {

	db := database.GetDB()
	defer db.Session.Close()
	runner := Runner{}
	err := db.C("results").Find(bson.M{"bib_number": bibNO}).One(&runner)
	if err != nil {
		http.Error(w, "runner not found", 404)
		return
	}

	var xInt, yInt int
	var fontSize float64
	background := image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))
	draw.Draw(background, background.Bounds(), colorWhite, image.ZP, draw.Src)
	bg, _ := certBG.(image.Image)
	draw.Draw(background, background.Bounds(), bg, image.Point{0, 0}, draw.Src)

	ctx.SetDst(background)

	runnerName := runner.FName + " " + runner.LName
	// draw Name
	xInt = 480
	yInt = 150
	fontSize = getFontSize(runnerName, 100, 400)
	xOffset := int(getWidth(runnerName, fontSize) / 2)
	ctx.SetFontSize(fontSize)
	pt := freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(utf8FontSize)>>6))
	ctx.DrawString(runnerName, pt)

	/*
		//draw chiptime
		// xInt, _ = strconv.Atoi(x)
		// yInt, _ = strconv.Atoi(y)
		xInt = 340
		yInt = 350
		fontSize = 50
		ctx.SetFontSize(fontSize)
		pt = freetype.Pt(xInt, yInt+int(ctx.PointToFixed(utf8FontSize)>>6))
		ctx.DrawString(formatTime(runner.ChipTime), pt)

		//draw guntime
		xInt = 420
		yInt = 435
		fontSize = 50
		ctx.SetFontSize(fontSize)
		pt = freetype.Pt(xInt, yInt+int(ctx.PointToFixed(utf8FontSize)>>6))
		ctx.DrawString(formatTime(runner.GunTime), pt)

	*/

	err = png.Encode(w, background)
	if err != nil {
		fmt.Println(err)
	}

}

//init load and instantiate reusable things
func init() {
	imgFile, err := os.Open("./cert.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_certBG, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	certBG = _certBG
	fontBytes, err := ioutil.ReadFile(utf8FontFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utf8Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	colorBlue = image.NewUniform(color.RGBA{0, 0, 255, 255})
	colorWhite = image.NewUniform(color.RGBA{255, 255, 255, 255})
	colorBlack = image.NewUniform(color.RGBA{0, 0, 0, 255})

	img4TxtWidthTest = image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))
	ctx = freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(utf8Font)
	ctx.SetClip(img4TxtWidthTest.Bounds())
	ctx.SetSrc(colorBlue)
	ctx4Test = freetype.NewContext()
	ctx4Test.SetDPI(dpi)
	ctx4Test.SetFont(utf8Font)
	ctx4Test.SetClip(img4TxtWidthTest.Bounds())
	ctx4Test.SetSrc(colorBlue)
	ctx4Test.SetDst(img4TxtWidthTest)
}

// getFontSize return largest possible FontSize that will make text length less than maxWidth , bad big-O but i don't give a shit
func getFontSize(txt string, fontSize float64, maxWidth int) float64 {
	goodSize := fontSize + 1.0
	currWidth := 1000000
	for currWidth > maxWidth {
		goodSize = goodSize - 1.0
		currWidth = getWidth(txt, goodSize)
	}
	return goodSize
}

// getWidth get width of txt that draw with specify fontSize
func getWidth(txt string, fontSize float64) int {
	pt := freetype.Pt(0, 0)
	ctx4Test.SetFontSize(fontSize)
	l, _ := ctx4Test.DrawString(txt, pt)
	return l.X.Floor()
}

// format time millisec to hh:mm:ss , mm:ss , ss
func formatTime(t float64) string {
	ret := ""
	ss := int(t / 1000)
	mm := int(ss / 60)
	ss = ss - mm*60
	hh := int(mm / 60)
	mm = mm - hh*60
	if hh > 0 {
		ret = fmt.Sprintf("%d", hh) + ":"
	}
	if hh > 0 || mm > 0 {
		if mm < 10 {
			ret = ret + "0"
		}
		ret = ret + fmt.Sprintf("%d", mm) + ":"
	}
	if hh > 0 || mm > 0 || ss > 0 {
		if ss < 10 {
			ret = ret + "0"
		}
		ret = ret + fmt.Sprintf("%d", ss)
	}
	return ret
}

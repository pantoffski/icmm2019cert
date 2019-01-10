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
	"math/rand"
	"net/http"
	"os"
	_ "strings"

	"github.com/globalsign/mgo/bson"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	backgroundWidth  = 1000
	backgroundHeight = 1000
	utf8FontFile     = "./font.ttf"
	utf8FontSize     = float64(45.0)
	spacing          = float64(1.5)
	dpi              = float64(72)
	ctx              = new(freetype.Context)
	utf8Font         = new(truetype.Font)
	certBG           interface{}
	colorPink        *image.Uniform
	colorBlue        *image.Uniform
	colorRed         *image.Uniform
	colorWhite       *image.Uniform
)

// Runner is a runner's result
type Runner struct {
	ID        bson.ObjectId `bson:"_id"`
	BibNO     int           `bson:"bibNumber"`
	FullBib   string        `bson:"fullBib"`
	Gender    string        `bson:"gender"`
	FName     string        `bson:"firstname"`
	LName     string        `bson:"lastname"`
	ChipTime  float64       `bson:"chiptime"`
	GunTime   float64       `bson:"guntime"`
	Challenge string        `bson:"challenge"`
	NameOnBib string        `bson:"nameOnBib"`
}

// func CertImgDebug(bibNO, txt, x, y, s string, w http.ResponseWriter) {

//Image gen cert img from bib
func Image(bibNO int, w http.ResponseWriter) {

	db := database.GetDB()
	defer db.Session.Close()
	runner := Runner{}
	err := db.C("bib_subscribers").Find(bson.M{"bibNumber": bibNO}).One(&runner)
	if err != nil {
		http.Error(w, "runner not found", 404)
		return
	}

	runner.GunTime = 2000000.0 + rand.Float64()*4000000
	runner.ChipTime = 2000000.0 + rand.Float64()*4000000
	var xInt, yInt int
	var fontSize float64
	background := image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))
	draw.Draw(background, background.Bounds(), colorWhite, image.ZP, draw.Src)
	bg, _ := certBG.(image.Image)
	draw.Draw(background, background.Bounds(), bg, image.Point{0, 0}, draw.Src)

	ctx.SetDst(background)
	ctx.SetClip(background.Bounds())

	runnerName := runner.FName + " " + runner.LName
	// draw Name
	xInt = 498
	yInt = 190
	fontSize = 80.0
	if len(runnerName) > 22 {
		fontSize = 70.0
	}
	if len(runnerName) > 26 {
		fontSize = 50.0
	}
	xOffset := int(getWidth(runnerName, fontSize) / 2)
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(colorPink)
	pt := freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(fontSize)>>6))
	ctx.DrawString(runnerName, pt)

	// draw bib
	xInt = 385
	yInt = 295
	fontSize = 35.0
	xOffset = int(getWidth(runner.FullBib, fontSize) / 2)
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(colorWhite)
	pt = freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(fontSize)>>6))
	ctx.DrawString(runner.FullBib, pt)

	// draw gender
	xInt = 790
	yInt = 295
	fontSize = 35.0
	ctx.SetFontSize(fontSize)
	gender := "FEMALE"
	if runner.Gender == "M" {
		gender = "MALE"
	}
	xOffset = int(getWidth(gender, fontSize) / 2)
	ctx.SetSrc(colorWhite)
	pt = freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(fontSize)>>6))
	ctx.DrawString(gender, pt)

	// draw guntime
	xInt = 500
	yInt = 430
	fontSize = 90.0
	xOffset = int(getWidth(formatTime(runner.GunTime), fontSize) / 2)
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(colorRed)
	pt = freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(fontSize)>>6))
	ctx.DrawString(formatTime(runner.GunTime), pt)

	// draw chiptime
	xInt = 635
	yInt = 560
	fontSize = 35.0
	xOffset = int(getWidth(formatTime(runner.ChipTime), fontSize) / 2)
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(colorWhite)
	pt = freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(fontSize)>>6))
	ctx.DrawString(formatTime(runner.ChipTime), pt)

	// //draw txt
	// xInt, _ = strconv.Atoi(x)
	// yInt, _ = strconv.Atoi(y)
	// sizeInt, _ := strconv.Atoi(s)

	// xOffset = int(getWidth(txt, float64(sizeInt)) / 2)
	// ctx.SetSrc(colorPink)
	// ctx.SetFontSize(float64(sizeInt))
	// pt = freetype.Pt(xInt-xOffset, yInt+int(ctx.PointToFixed(float64(sizeInt))>>6))
	// ctx.DrawString(txt, pt)

	err = png.Encode(w, background)
	if err != nil {
		fmt.Println(err)
	}

}

//init load and instantiate reusable things
func init() {
	imgFile, err := os.Open("./cert.png")
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
	colorRed = image.NewUniform(color.RGBA{129, 23, 25, 255})
	colorPink = image.NewUniform(color.RGBA{232, 133, 162, 255})

	ctx = freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(utf8Font)
	ctx.SetSrc(colorBlue)
}

// getWidth get width of txt that draw with specify fontSize
func getWidth(txt string, fontSize float64) int {
	pt := freetype.Pt(0, 2000)
	ctx.SetFontSize(fontSize)
	l, _ := ctx.DrawString(txt, pt)
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
	if hh < 10 {
		ret = ret + "0"
	}
	ret = ret + fmt.Sprintf("%d", hh) + ":"
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

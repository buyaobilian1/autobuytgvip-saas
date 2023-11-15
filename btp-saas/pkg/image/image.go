package image

import (
	"bytes"
	"encoding/base64"
	"github.com/skip2/go-qrcode"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
	"image/png"
	"net/url"
	"os"
	"strings"
)

func FromSvgBase64(data string) string {

	svgBase64 := strings.Replace(data, "data:image/svg+xml;base64,", "", 1)
	//svgBase64 := data
	unescape, err := url.QueryUnescape(svgBase64)
	if err != nil {
		panic(err)
	}
	svgBytes, err := base64.RawStdEncoding.DecodeString(unescape)
	if err != nil {
		panic(err)
	}

	//w, h := 160, 160
	icon, _ := oksvg.ReadIconStream(bytes.NewReader(svgBytes))
	w := int(icon.ViewBox.W)
	h := int(icon.ViewBox.H)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
	//buf := new(bytes.Buffer)
	f, _ := os.Create("out.png")
	defer f.Close()
	err = png.Encode(f, rgba)
	//os.WriteFile("out.png", buf.Bytes(), 0644)

	return string(svgBytes)
}

func GenQrcode(data string) *bytes.Buffer {
	b, _ := qrcode.Encode(data, qrcode.Medium, 128)
	buf := new(bytes.Buffer)
	buf.Write(b)
	return buf
}

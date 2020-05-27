package middleware

import (
	"fmt"
	"strings"
	"image"
	"github.com/disintegration/imaging"
	"github.com/syyongx/php2go"
)


//加油图片
func LightTagImg(headImgUrl ,tagName string) string {

 	lidx := strings.LastIndex(headImgUrl, "/")
	//
	bm, err := imaging.Open(saveImg(headImgUrl[0:lidx]+"/0"))
	if err != nil {
		fmt.Printf("open file failed," + err.Error())
	}
	dst := imaging.Resize(bm, 640, 640, imaging.Lanczos)     // 图片按比例缩放
	//
	tag, err := imaging.Open("img/tag/"+tagName+".png")
	if err != nil {
		fmt.Printf("open file failed")
	}

	result := imaging.Overlay(dst, tag, image.Pt(0, 0), 1)

	//
	fileName := fmt.Sprintf("./img/wuhan/%s.jpg", php2go.Md5(headImgUrl))
	err = imaging.Save(result, fileName)
	if err != nil {
		return ""
	}

	return fileName
}


//分享图片
func LightShareImg(qrcode string) string {

	bm, err := imaging.Open("img/tag/fengxiangwuhan.jpg")
	if err != nil {
		fmt.Printf("open file failed," + err.Error())
		return ""
	}
	//
	tag, err := imaging.Open(saveImg(qrcode))
	if err != nil {
		fmt.Printf("open qrcode file failed")
		return ""
	}
	qr := imaging.Resize(tag, 150, 150, imaging.Lanczos)     // 图片按比例缩放
	//
	result := imaging.Overlay(bm, qr, image.Pt(780, 740), 1)

	//
	fileName := fmt.Sprintf("./img/wuhan/whs_%s.jpg", php2go.Md5(qrcode))
	err = imaging.Save(result, fileName)
	if err != nil {
		return ""
	}

	return fileName
}




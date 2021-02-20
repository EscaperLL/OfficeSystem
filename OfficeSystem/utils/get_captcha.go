package utils

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type Captcha struct {
	Id string
	Base64 string
	Code int
}

var g_store = base64Captcha.DefaultMemStore

func GetCaptcha() (id,base string,err error) {
	fots:=[]string{"wqy-microhei.ttc"}
	driver:=base64Captcha.NewDriverMath(50,140,20,20,&color.RGBA{0,0,0,0},fots)
	captcha :=base64Captcha.NewCaptcha(driver,g_store)
	id,base,err = captcha.Generate()
	return id,base,err
}

func VerifyCaptcha(id,base64 string)bool  {
	return g_store.Verify(base64,id,true)
}
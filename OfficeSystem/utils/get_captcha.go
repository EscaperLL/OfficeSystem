package utils

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type Captcha struct {
	Id string
	Base64 string
}

func GetCaptcha() (id,base string,err error) {
	fots:=[]string{"wqy-microhei.ttc"}
	driver:=base64Captcha.NewDriverMath(50,140,20,20,&color.RGBA{0,0,0,0},fots)
	store := base64Captcha.DefaultMemStore
	captcha :=base64Captcha.NewCaptcha(driver,store)
	id,base,err = captcha.Generate()
	return id,base,err
}

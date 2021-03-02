package utils

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

func String2Bytes(str string)[]byte  {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))

	bytes :=reflect.SliceHeader{
		Data: stringHeader.Data,
		Len: stringHeader.Len,
		Cap: stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bytes))
}

func Bytes2String(b []byte)string  {
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	stringHeader := reflect.StringHeader{
		Data: bytesHeader.Data,
		Len: bytesHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&stringHeader))
}


func GetMd5(input string)string  {
	byte_str := ([]byte)(input)
	result:= md5.Sum(byte_str)
	fmt.Println("md5",result)
	result_str :=string(result[:])
	return  result_str
}

func GetDateTimeStr()  string{
	now:= time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d   ",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second())
}
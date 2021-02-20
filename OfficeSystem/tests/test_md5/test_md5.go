package test_md5


import (
	"OfficeSystem/utils"
	"crypto/md5"
	"fmt"
)

func GetMd5(input string)string  {
	byte_str := utils.String2Bytes(input)
	result:= md5.Sum(byte_str)
	result_str :=string(result[:])
	return  result_str
}

func main()  {
	fmt.Println(GetMd5("asdasdasd"))
}

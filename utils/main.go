package tools

import (
	"math/rand"
	"time"
)

var (
	//ProductName 产品名
	ProductName string
)

//SetProductName 设置产品名
func SetProductName(name string) {
	ProductName = name
}

func init() {
	rand.Seed(time.Now().UnixNano())
	SetProductName("UnknowProduct")
}

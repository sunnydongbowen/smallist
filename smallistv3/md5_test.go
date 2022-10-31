package main

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	pwd := md5Secret("123456")
	fmt.Println(pwd)

}

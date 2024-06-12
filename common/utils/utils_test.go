package utils

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func Test_CalcHash(t *testing.T) {
	filePath := "J:\\桌面快捷入口\\待阅读文件\\金舒航_开发_15088282732.pdf"
	f, e := os.Open(filePath)
	if e != nil {
		log.Fatalln(e)
	}
	defer f.Close()
	hash := CalculateHash(f)
	fmt.Println(hash)
}

package service

import (
	"fmt"
	"os"
)

var Imfile imfile

type imfile struct {
}

func (i *imfile) ImportFile() (string, error) {
	//filepath := "E:\\Users\\18206\\Downloads\\kubeagent.yaml"
	filepath := "yaml/kubeagent.yaml"
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("打开文件出错")
		return "", err
	}
	defer f.Close()

	var file []byte
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("文件读取完毕")
				break
			} else {
				fmt.Println("文件读取出错")
				break
			}
		}
		file = append(file, buf[:n]...)
	}
	//fmt.Println("读取：", string(file))
	return string(file), nil
}

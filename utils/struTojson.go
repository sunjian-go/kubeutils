package utils

import (
	"bytes"
	"encoding/json"
	"errors"
)

var Stj stj

type stj struct {
}

// 结构体转json
func (s *stj) StructToJson(stru interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(stru)
	if err != nil {
		Logg.Error("编码结构体为 JSON 时出错：" + err.Error())
		return nil, errors.New("结构体编码结构体为 JSON 时出错：" + err.Error())
	}
	// 创建一个包含 JSON 数据的 io.Reader
	jsonReader := bytes.NewBuffer(data)
	return jsonReader, nil
}

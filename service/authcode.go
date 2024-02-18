package service

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Auth auth

type auth struct {
}

// 生成随机公式
func (a *auth) GetAuthCode() string {
	//设置种子
	rand.Seed(time.Now().UnixNano())
	//随机生成数字和符号
	num1 := rand.Intn(9) + 1
	num2 := rand.Intn(9) + 1
	opt := []string{"+", "-", "*"}
	index := rand.Intn(3)
	formula := strconv.Itoa(num1) + opt[index] + strconv.Itoa(num2) + "=?"
	//fmt.Println("随机公式：", formula)
	return formula
}

// 计算结果进行对比
func (a *auth) CountResult(formula, res string) error {
	newformula := strings.Split(formula, "=")
	newFormula := strings.Split(newformula[0], "")
	//fmt.Println("公式：", newFormula[0], newFormula[1], newFormula[2])
	var result int
	num1, _ := strconv.Atoi(newFormula[0])
	num2, _ := strconv.Atoi(newFormula[2])
	switch newFormula[1] {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	}
	//fmt.Println("结果：", result)

	if strconv.Itoa(result) != res {
		return errors.New("验证不通过")
	}
	return nil
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"unsafe"
)

// 表示四则算式中的数字个数
// 2,3,4的比例不同是考虑到 2 3年级的学生算数能力
var levelTwo = []int{2, 2, 2, 2, 3, 3, 3, 4}
var levelThree = []int{2, 3, 3, 3, 4, 4, 4, 4, 4, 4}

// 用来判断算式唯一性并存储算式结果
var secondGradeFourArithmeticMap = make(map[string]int, 0)
var thirdGradefourArithmeticMap = make(map[string]int, 0)

// 存储三年级四则运算算式便于按顺序输出
var secondGradeFourArithmeticSlice = make([]string, 0)
var thirdGradefourArithmeticSlice = make([]string, 0)

type fourArithmetic struct {
	arithmetic []interface{}
	len        int
}

// level 表示符号的优先级
// 0表示 “(” 优先级最低
// 遇到 “)” 需要一直判断直到遇到 “(” 所以不设置优先级也可以
var level = map[string]int{
	"(": 0,
	"+": 1,
	"-": 1,
	"×": 2,
	"÷": 2,
}

func main() {
	// 生成算式
	for len(secondGradeFourArithmeticMap) <= 1000 {
		secondGrade()
	}
	for len(thirdGradefourArithmeticMap) <= 1000 {
		thirdGrade()
	}
	// 输出算式
	//1.创建一个新文件，写入内容 hello,Gardon
	//1.打开文件 d:/adc.txt
	secondGradefilePath := "C:/Users/xl/Desktop/SecondGrade.txt"
	secondFile, err := os.OpenFile(secondGradefilePath, os.O_RDWR|os.O_CREATE, 0777) //参数3只要在linux系统下才会生效
	if err != nil {
		fmt.Printf("open file err=%v", err)
		return
	}
	//及时关闭file句柄 避免文件泄露
	defer secondFile.Close()
	//写入时，使用带缓存的 *Writer
	secondWriter := bufio.NewWriter(secondFile)
	for _, s := range secondGradeFourArithmeticSlice {
		secondWriter.WriteString(s + "=" + "\n")
	}
	secondWriter.Flush()
	thirdGradefilePath := "C:/Users/xl/Desktop/ThirdGrade.txt"
	thirdFile, err := os.OpenFile(thirdGradefilePath, os.O_RDWR|os.O_CREATE, 0777) //参数3只要在linux系统下才会生效
	if err != nil {
		fmt.Printf("open file err=%v", err)
		return
	}
	//及时关闭file句柄 避免文件泄露
	defer thirdFile.Close()
	//写入时，使用带缓存的 *Writer
	thirdWriter := bufio.NewWriter(thirdFile)
	for _, s := range thirdGradefourArithmeticSlice {
		thirdWriter.WriteString(s + "=" + "\n")
	}
	thirdWriter.Flush()
	//1.创建一个新文件，写入内容 hello,Gardon
	//1.打开文件 d:/adc.txt
	secondGradeAnswerfilePath := "C:/Users/xl/Desktop/SecondGradeAnswer.txt"
	secondAnswerFile, err := os.OpenFile(secondGradeAnswerfilePath, os.O_RDWR|os.O_CREATE, 0777) //参数3只要在linux系统下才会生效
	if err != nil {
		fmt.Printf("open file err=%v", err)
		return
	}
	//及时关闭file句柄 避免文件泄露
	defer secondAnswerFile.Close()
	//写入时，使用带缓存的 *Writer
	secondAnswerWriter := bufio.NewWriter(secondAnswerFile)
	for _, s := range secondGradeFourArithmeticSlice {
		secondAnswerWriter.WriteString(s + "=" + strconv.Itoa(secondGradeFourArithmeticMap[s]) + "\n")
	}
	secondAnswerWriter.Flush()
	thirdGradeAnswerfilePath := "C:/Users/xl/Desktop/ThirdGradeAnswer.txt"
	thirdAnswerFile, err := os.OpenFile(thirdGradeAnswerfilePath, os.O_RDWR|os.O_CREATE, 0777) //参数3只要在linux系统下才会生效
	if err != nil {
		fmt.Printf("open file err=%v", err)
		return
	}
	//及时关闭file句柄 避免文件泄露
	defer thirdAnswerFile.Close()
	//写入时，使用带缓存的 *Writer
	thirdAnswerWriter := bufio.NewWriter(thirdAnswerFile)
	for _, s := range thirdGradefourArithmeticSlice {
		thirdAnswerWriter.WriteString(s + "=" + strconv.Itoa(thirdGradefourArithmeticMap[s]) + "\n")
	}
	thirdAnswerWriter.Flush()
}

// 生成二年级算式并存储，只有 + -
func secondGrade() (s []interface{}) {
	midds := []string{"+", "-"}
	count := levelTwo[rand.Intn(len(levelTwo))] // 表示算式中数字的个数
	for i := 0; i < count; i++ {
		num := rand.Intn(100) + 1
		s = append(s, num)
		midd := midds[rand.Intn(len(midds))]
		if i != count-1 {
			s = append(s, midd)
		}
	}
	resultNum, err := convert(s)
	if err == nil && resultNum <= 100 {
		value := convertString(s, resultNum)
		if _, ok := secondGradeFourArithmeticMap[value]; !ok {
			secondGradeFourArithmeticMap[value] = resultNum
			secondGradeFourArithmeticSlice = append(secondGradeFourArithmeticSlice, value)
		}
	}
	return
}

// 生成符合规则的三年级四则运算算式并存储
// 1. 保证算式中间结果都是正整数
// 2. 保证算式中的 （） 都是有效果的
// 3. 保证算式的唯一性
// 三年级式子，有 + - * / （）
// （）外侧符号为 * / 且内侧相邻符号为 + -的情况（）才能生效
func thirdGrade() {
	midds := []string{"+", "-", "×", "÷"}
	count := levelThree[rand.Intn(len(levelThree))] // 表示算式中数字的个数
	// 获得原始算式
	s := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		num := rand.Intn(100) + 1
		s = append(s, num)
		midd := midds[rand.Intn(len(midds))]
		if i != count-1 {
			s = append(s, midd)
		}
	}
	// 根据原始算式得到该算式加括号的所有情况
	// 如果长度为 3 表示为 a+b 形式 不需要加小括号，直接存储
	if len(s) == 3 {
		//fmt.Println(s)
		resultNum, err := convert(s)
		if err == nil {
			value := convertString(s, resultNum)
			if _, ok := thirdGradefourArithmeticMap[value]; !ok {
				thirdGradefourArithmeticMap[value] = resultNum
				thirdGradefourArithmeticSlice = append(thirdGradefourArithmeticSlice, value)
			}
		} else {
			//fmt.Println(s,"不符合规则")
		}
	} else {
		//fmt.Println(s, "长度：", len(s))
		fourArithmeticSet(s)
	}
	return
}

// 利用 原始四则算式 得出所有的 加小括号的四则算式并存储
func fourArithmeticSet(s []interface{}) {
	//before := len(fourArithmeticMap)
	//fmt.Println("原始：",s)
	if len(s) == 5 {
		// 存储一份原始运算式
		s0 := *(*[5]interface{})(unsafe.Pointer(&s[0]))
		arithmeticCheck(s)
		// 加括号 (a+b)*c
		s = sliceInster(s, 0, "(")
		s = sliceInster(s, 4, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
		// 加括号 a+(b*c)
		s = sliceInster(s, 2, "(")
		s = append(s, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
	} else if len(s) == 7 {
		// 存储一份原始运算式
		s0 := *(*[7]interface{})(unsafe.Pointer(&s[0]))
		arithmeticCheck(s)
		// 加括号 (a+b)*c-d
		s = sliceInster(s, 0, "(")
		s = sliceInster(s, 4, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
		// (a+b*c)*d
		s = sliceInster(s, 0, "(")
		s = sliceInster(s, 6, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
		// a+(b*c)*d
		s = sliceInster(s, 2, "(")
		s = sliceInster(s, 6, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
		// a+(b*c*d)
		s = sliceInster(s, 2, "(")
		s = append(s, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
		// a+b*(c*d)
		s = sliceInster(s, 4, "(")
		s = append(s, ")")
		arithmeticCheck(s)
		// 初始化 s
		s = s0[:]
	}
	//after := len(fourArithmeticMap)
	//if after-before == 0 {
	//	fmt.Println(s, "没有符合条件的式子")
	//}
}

// 通过后缀表达式计算四则表达式的出算式结果
// 检查中间计算过程是否是 正整数 主要是 除法 和 减法
func result(s []interface{}) (resultNum int, err error) {
	middResults := make([]int, 0)
	for _, value := range s {
		// 数字处理
		if value1, isInt := value.(int); isInt {
			middResults = append(middResults, value1)
		} else if value1, isString := value.(string); isString { // 运算符处理
			//fmt.Println(value,"是运算符")
			right := middResults[len(middResults)-1]
			left := middResults[len(middResults)-2]
			middResults = middResults[:len(middResults)-2]
			var middResult int
			resultBool := true
			if value1 == "+" {
				middResult = left + right
			} else if value1 == "-" {
				middResult = left - right
				if middResult < 0 { // 验证中间结果是否为正数
					resultBool = false
				}
			} else if value1 == "×" {
				middResult = left * right
			} else if value1 == "÷" {
				if right == 0 { // 保证 ÷ 右边为 非零
					resultBool = false
				} else if left%right != 0 { // 验证中间结果是否为正整数
					resultBool = false
				} else {
					middResult = left / right
				}
			}
			// 判断中间结果是否 >= 0 且是整数
			if !resultBool {
				err = errors.New("中间结果为非正整数")
				return 0, err
			} else {
				middResults = append(middResults, middResult)
			}
		}
	}
	resultNum = middResults[0]
	if resultNum > 10000 {
		err = errors.New("结果太大，不符合")
		return 0, err
	}
	return resultNum, err
}

// 中缀表达式转后缀表达式并返回计算结果
func convert(s []interface{}) (resultNum int, err error) {
	// 存储后缀表达式
	postfixExpression := make([]interface{}, 0)
	// 作为运算符的中转站
	midd := make([]string, 0) //这里不能用 chan  需要换成别的，比如切片
	for _, value := range s {
		// 数字处理
		if value1, isInt := value.(int); isInt {
			postfixExpression = append(postfixExpression, value1)
		} else if value2, isString := value.(string); isString { // 运算符处理
			//	} else if  { // 运算符处理
			if len(midd) == 0 {
				midd = append(midd, value2)
			} else if value2 == "(" {
				midd = append(midd, value2)
			} else if value2 == "+" {
				for a := midd[len(midd)-1]; level["+"] <= level[a] && len(midd) != 0; {
					// 表示栈顶运算符优先级 >= 新运算符优先级，将其放入后缀表达式
					midd = midd[:len(midd)-1]
					postfixExpression = append(postfixExpression, a)
					if len(midd) != 0 {
						a = midd[len(midd)-1]
					}
				}
				midd = append(midd, "+")
			} else if value2 == "-" {
				for a := midd[len(midd)-1]; level["-"] <= level[a] && len(midd) != 0; {
					// 表示栈顶运算符优先级 >= 新运算符优先级，将其出栈并放入后缀表达式
					midd = midd[:len(midd)-1]
					postfixExpression = append(postfixExpression, a)
					if len(midd) != 0 {
						a = midd[len(midd)-1]
					}
				}
				midd = append(midd, "-")
			} else if value2 == "×" {
				for a := midd[len(midd)-1]; level["×"] <= level[a] && len(midd) != 0; {
					// 表示栈顶运算符优先级 >= 新运算符优先级，将其放入后缀表达式
					midd = midd[:len(midd)-1]
					postfixExpression = append(postfixExpression, a)
					if len(midd) != 0 {
						a = midd[len(midd)-1]
					}
				}
				midd = append(midd, "×")
			} else if value2 == "÷" {
				for a := midd[len(midd)-1]; level["÷"] <= level[a] && len(midd) != 0; {
					// 表示栈顶运算符优先级 >= 新运算符优先级，将其放入后缀表达式
					midd = midd[:len(midd)-1]
					postfixExpression = append(postfixExpression, a)
					if len(midd) != 0 {
						a = midd[len(midd)-1]
					}
				}
				midd = append(midd, "÷")
			} else if value2 == ")" {
				for a := midd[len(midd)-1]; a != "("; {
					// 表示栈顶运算符优先级 >= 新运算符优先级，将其放入后缀表达式
					midd = midd[:len(midd)-1]
					postfixExpression = append(postfixExpression, a)
					if len(midd) != 0 {
						a = midd[len(midd)-1]
					}
				}
				// 删除栈中的 (
				midd = midd[:len(midd)-1]
			}
		}
	}
	for _, value := range midd {
		postfixExpression = append(postfixExpression, value)
	}
	resultNum, err = result(postfixExpression)
	return resultNum, err
}

// 往切片任意位置插入元素 index 表示插入的位置 s1 表示要插入的元素
func sliceInster(s []interface{}, index int, s1 interface{}) []interface{} {
	s = append(s, 0) //先把原来的切片长度+1
	copy(s[index+1:], s[index:])
	s[index] = s1 //新元素的值是0
	return s
}

// 验算式子是否符合规范
// 1. 小括号是否有效果
func arithmeticCheck(s []interface{}) (isOk bool) {
	isOk = true
	// midd 存储运算符
	midd := make([]string, 0)
	// isEfficient 来代表是否有效 0代表初始值 1 代表一方有效或者没有小括号 2 代表两方有效
	isEfficient := 0
	for index, value := range s {
		if value1, isString := value.(string); isString {
			if value1 != ")" && value1 != "(" {
				midd = append(midd, value1)
			} else if value1 == "(" && len(midd) != 0 &&
				(midd[len(midd)-1] == "×" || midd[len(midd)-1] == "÷") &&
				(s[index+2] == "+" || s[index+2] == "-") {
				isEfficient += 1 // "(" 有效
			} else if value1 == ")" && index != len(s)-1 &&
				(midd[len(midd)-1] == "+" || midd[len(midd)-1] == "-") &&
				(s[index+1] == "×" || s[index+1] == "÷") {
				isEfficient += 1 // ")" 有效
			}
		}
	}
	// 表示式子不含有小括号
	if len(midd) == len(s)/2 {
		isEfficient += 1
	}
	// 如果小括号有效果 就往下进行
	if isEfficient > 0 {
		resultNum, err := convert(s)
		if err == nil {
			value := convertString(s, resultNum)
			// 利用 map 确保唯一性，为了保证加括号的算式的顺序，存入 slice
			if _, ok := thirdGradefourArithmeticMap[value]; !ok {
				thirdGradefourArithmeticMap[value] = resultNum
				thirdGradefourArithmeticSlice = append(thirdGradefourArithmeticSlice, value)
			}
		} else {
			isOk = false
			//fmt.Println(s, err)
		}
	} else {
		isOk = false
		//fmt.Println(s, "中的小括号没用")
	}
	return isOk
}

// convertString 将式子转换为 string 形式；便于存入 map
func convertString(s []interface{}, resultNum int) string {
	var value string
	for _, i2 := range s {
		if value1, isString := i2.(string); isString {
			value += value1
		} else if value1, isInt := i2.(int); isInt {
			value += strconv.Itoa(value1)
		}
	}
	return value
}

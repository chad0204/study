package main

import (
	"flag"
	"fmt"
)

type Celsius float64    //摄氏度
type Fahrenheit float64 //华氏度

/**
实现
type Value interface {
	String() string
	Set(string) error
}
*/
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) //将s 解析出一个浮点数和一个字符串
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

//Celsius实现了String, celsiusFlag内嵌了Celsius, 就不用实现了
func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}

func CToF(celsius Celsius) Fahrenheit {
	//T(x), 将x转为T类型, 只有x和T的底层类型一样才能转换
	return Fahrenheit(celsius*9/5 + 32)
}

func FToC(fahrenheit Fahrenheit) Celsius {
	return Celsius((fahrenheit - 32) * 5 / 9)
}

var temp = CelsiusFlag("temp", 20.0, "the temperature")

//flag.Value: 给命令行标记定义新的符号
func main() {
	//输入                        输出
	//.\celsius.exe              20°C
	//.\celsius.exe -temp -18C   -18°C
	//.\celsius.exe -temp -18F   -7.777777777777778°C

	flag.Parse()
	fmt.Println(*temp)
}

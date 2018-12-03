package main

//假设某公司有两个员工，一个普通员工和一个高级员工， 但是基本薪资是相同的，高级员工多拿奖金。计算公司为员工的总开支。
import (
	"fmt"
	"reflect"
)

// 薪资计算器接口
type SalaryCalculator interface {
	CalculateSalary() int
}

// 普通挖掘机员工
type Contract struct {
	empId    int
	basicpay int
}

// 有蓝翔技校证的员工
type Permanent struct {
	empId    int
	basicpay int
	jj       int // 奖金
}

func (p Permanent) CalculateSalary() int {
	return p.basicpay + p.jj
}

func (c Contract) CalculateSalary() int {
	return c.basicpay
}

// 总开支
func totalExpense(s []SalaryCalculator) {
	expense := 0

	for _, v := range s {

		getType := reflect.TypeOf(v)
		fmt.Println("get Type is :", getType.Name())

		getValue := reflect.ValueOf(v)
		fmt.Println("get all Fields is:", getValue)
		// 获取方法字段
		// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
		// 2. 再通过reflect.Type的Field获取其Field
		// 3. 最后通过Field的Interface()得到对应的value
		// for i := 0; i < getType.NumField(); i++ {
		// 	field := getType.Field(i)
		// 	value := getValue.Field(i)
		// 	fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		// }

		fmt.Printf("个人开支id->%+v:$%v \n", getValue.Field(0), v.CalculateSalary())
		expense = expense + v.CalculateSalary()
	}

	fmt.Printf("总开支 $%d \n", expense)
}

func main() {
	pemp1 := Permanent{1, 3000, 10000}
	//fmt.Printf("%v \n", pemp1.empId)
	pemp2 := Permanent{2, 3000, 20000}
	cemp1 := Contract{3, 3000}
	employees := []SalaryCalculator{pemp1, pemp2, cemp1}
	totalExpense(employees)
}

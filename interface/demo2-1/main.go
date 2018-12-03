package main

//假设某公司有两个员工，一个普通员工和一个高级员工， 但是基本薪资是相同的，高级员工多拿奖金。计算公司为员工的总开支。
import (
	"fmt"
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
		//断言
		switch f := v.(type) {
		case SalaryCalculator:
			a := f.CalculateSalary()
			fmt.Println(a)
		default:
			fmt.Printf("unknown type\n")
		}
		//fmt.Printf("每个人开支:%+v-$%v \n", v, v.CalculateSalary())
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

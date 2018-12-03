package main

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
		expense = expense + v.CalculateSalary()
	}
	fmt.Printf("总开支 $%d", expense)
}

func main() {
	pemp1 := Permanent{1, 3000, 10000}
	pemp2 := Permanent{2, 3000, 20000}
	cemp1 := Contract{3, 3000}
	employees := []SalaryCalculator{pemp1, pemp2, cemp1}
	fmt.Printf("总开支 $%d", employees)

	totalExpense(employees)
}

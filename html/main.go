package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var myTemplate *template.Template

type Person struct {
	Name string
	age  string
}

func userInfo(w http.ResponseWriter, r *http.Request) {

	p := Person{Name: "safly", age: "30"}
	myTemplate.Execute(w, p)

}

func initTemplate(fileName string) (err error) {
	myTemplate, err = template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("parse file err:", err)
		return
	}
	return
}

/*
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
*/

func main() {
	initTemplate("e:/golang/go_pro/src/safly/index.html")
	http.HandleFunc("/user/info", userInfo)
	err := http.ListenAndServe("0.0.0.0:8880", nil)
	if err != nil {
		fmt.Println("http listen failed")
	}
}

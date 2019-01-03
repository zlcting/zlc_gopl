# **反射基础和反射三大定律**

反射的目标之一是获取变量的类型信息，例如这个类型的名称、占用字节数、所有的方法列表、所有的内部字段结构、它的底层存储类型等等

>func TypeOf(v interface{}) Type

>func ValueOf(v interface{}) Value

```
package main

import "fmt"
import "reflect"

func main() {
    var s int = 42
    fmt.Println(reflect.TypeOf(s))
    fmt.Println(reflect.ValueOf(s))
}
--------
int
42
```
>interface{} 类型的结构如下图

[点击跳转看大图](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mK3Ext7BFASMhf5gkM0ZzMZMElibqsuo4vOU58xBjAa9r3fjfTjniaYibrGd68YBLuUtfzpDVK7LJdicA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)
![图片](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mK3Ext7BFASMhf5gkM0ZzMZMElibqsuo4vOU58xBjAa9r3fjfTjniaYibrGd68YBLuUtfzpDVK7LJdicA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

>reflect.Type

是一个接口类型，里面定义了非常多的方法用于获取和这个类型相关的一切信息。这个接口的结构体实现隐藏在 reflect 包里，每一种类型都有一个相关的类型结构体来表达它的结构信息。

```
type Type interface {
  ...
  Method(i int) Method  // 获取挂在类型上的第 i'th 个方法
  ...
  NumMethod() int  // 该类型上总共挂了几个方法
  Name() string // 类型的名称
  PkgPath() string // 所在包的名称
  Size() uintptr // 占用字节数
  String() string // 该类型的字符串形式
  Kind() Kind // 元类型
  ...
  Bits() // 占用多少位
  ChanDir() // 通道的方向
  ...
  Elem() Type // 数组，切片，通道，指针，字典(key)的内部子元素类型
  Field(i int) StructField // 获取结构体的第 i'th 个字段
  ...
  In(i int) Type  // 获取函数第 i'th 个参数类型
  Key() Type // 字典的 key 类型
  Len() int // 数组的长度
  NumIn() int // 函数的参数个数
  NumOut() int // 函数的返回值个数
  Out(i int) Type // 获取函数 第 i'th 个返回值类型
  common() *rtype // 获取类型结构体的共同部分
  uncommon() *uncommonType // 获取类型结构体的不同部分
}
```
所有的类型结构体都包含一个共同的部分信息，这部分信息使用 rtype 结构体描述，rtype 实现了 Type 接口的所有方法。剩下的不同的部分信息各种特殊类型结构体都不一样。可以将 rtype 理解成父类，特殊类型的结构体是子类，会有一些不一样的字段信息。
```
// 基础类型 rtype 实现了 Type 接口
type rtype struct {
  size uintptr // 占用字节数
  ptrdata uintptr
  hash uint32 // 类型的hash值
  ...
  kind uint8 // 元类型
  ...
}

// 切片类型
type sliceType struct {
  rtype
  elem *rtype // 元素类型
}

// 结构体类型
type structType struct {
  rtype
  pkgPath name  // 所在包名
  fields []structField  // 字段列表
}
```
>reflect.Value
不同于 reflect.Type 接口，reflect.Value 是结构体类型，一个非常简单的结构体。
```
type Value struct {
  typ *rtype  // 变量的类型结构体
  ptr unsafe.Pointer // 数据指针
  flag uintptr // 标志位
}

```
这个接口体包含变量的类型结构体指针、数据的地址指针和一些标志位信息。里面的类型结构体指针字段就是上面的 rtype 结构体地址，存储了变量的类型信息。标志位里有几个位存储了值的「元类型」。下面我们看个简单的例子

```
package main

import "reflect"
import "fmt"

func main() {
    type SomeInt int
    var s SomeInt = 42
    var t = reflect.TypeOf(s)
    var v = reflect.ValueOf(s)
 // reflect.ValueOf(s).Type() 等价于 reflect.TypeOf(s)
 fmt.Println(t == v.Type())
    fmt.Println(v.Kind() == reflect.Int) // 元类型
 // 将 Value 还原成原来的变量
 var is = v.Interface()
 fmt.Println(is.(SomeInt))
}

----------
true
true
42
```
Value 结构体的 Type() 方法也可以返回变量的类型信息，它可以作为 reflect.TypeOf() 函数的替代品，没有区别。通过 Value 结构体提供的 Interface() 方法可以将 Value 还原成原来的变量值。

将上面的各种关系整理一下，可以得到下面这张图

[点击跳转看大图](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mK3Ext7BFASMhf5gkM0ZzMZUQDSZadtU8kL0z7sSkQrBIEgq4TpCx97HEEfZwwnNLlnLPf4xRCrXQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)
![图片](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mK3Ext7BFASMhf5gkM0ZzMZUQDSZadtU8kL0z7sSkQrBIEgq4TpCx97HEEfZwwnNLlnLPf4xRCrXQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

Value 这个结构体虽然很简单，但是附着在 Value 上的方法非常之多，主要是用来方便用户读写 ptr 字段指向的数据内存。虽然我们也可以通过 unsafe 包来精细操控内存，但是使用过于繁琐，使用 Value 结构体提供的方法会更加简单直接。
```
 func (v Value) SetLen(n int)  // 修改切片的 len 属性
 func (v Value) SetCap(n int) // 修改切片的 cap 属性
 func (v Value) SetMapIndex(key, val Value) // 修改字典 kv
 func (v Value) Send(x Value) // 向通道发送一个值
 func (v Value) Recv() (x Value, ok bool) // 从通道接受一个值
 // Send 和 Recv 的非阻塞版本
 func (v Value) TryRecv() (x Value, ok bool)
 func (v Value) TrySend(x Value) bool

 // 获取切片、字符串、数组的具体位置的值进行读写
 func (v Value) Index(i int) Value
 // 根据名称获取结构体的内部字段值进行读写
 func (v Value) FieldByName(name string) Value
 // 将接口变量装成数组，一个是类型指针，一个是数据指针
 func (v Value) InterfaceData() [2]uintptr
 // 根据名称获取结构体的方法进行调用
 // Value 结构体的数据指针 ptr 可以指向方法体
 func (v Value) MethodByName(name string) Value
```

值得注意的是，观察 Value 结构体提供的很多方法，其中有不少会返回 Value 类型。比如反射数组类型的 Index(i int) 方法，它会返回一个新的 Value 对象，这个对象的类型指向数组内部子元素的类型，对象的数据指针会指向数组指定位置子元素所在的内存。

**理解 Go 语言官方的反射三大定律**
官方对 Go 语言的反射功能做了一个抽象的描述，总结出了三大定律，分别是

1.Reflection goes from interface value to reflection object.

2.Reflection goes from reflection object to interface value.

3.To modify a reflection object, the value must be settable.


第一个定律的意思是反射将接口变量转换成反射对象 Type 和 Value，这个很好理解，就是下面这两个方法的功能

```
func TypeOf(v interface{}) Type
func ValueOf(v interface{}) Value
```
第二个定律的意思是反射可以通过反射对象 Value 还原成原先的接口变量，这个指的就是 Value 结构体提供的 Interface() 方法。注意它得到的是一个接口变量，如果要换成成原先的变量还需要经过一次造型。

```
func (v Value) Interface() interface{}
```

前两个定律比较简单，它的意思可以使用前面画的反射关系图来表达。第三个定律的功能不是很好理解，它的意思是想用反射功能来修改一个变量的值，前提是这个值可以被修改。

值类型的变量是不可以通过反射来修改，因为在反射之前，传参的时候需要将值变量转换成接口变量，值内容会被浅拷贝，反射对象 Value 指向的数据内存地址不是原变量的内存地址，而是拷贝后的内存地址。这意味着如果值类型变量可以通过反射功能来修改，那么修改操作根本不会影响到原变量的值，那就白白修改了。所以 reflect 包就直接禁止了通过反射来修改值类型的变量。我们看个例子
```
package main

import "reflect"

func main() {
    var s int = 42
    var v = reflect.ValueOf(s)
    v.SetInt(43)
}

---------
panic: reflect: reflect.Value.SetInt using unaddressable value

goroutine 1 [running]:
reflect.flag.mustBeAssignable(0x82)
    /usr/local/go/src/reflect/value.go:234 +0x157
reflect.Value.SetInt(0x107a1a0, 0xc000016098, 0x82, 0x2b)
    /usr/local/go/src/reflect/value.go:1472 +0x2f
main.main()
    /Users/qianwp/go/src/github.com/pyloque/practice/main.go:8 +0xc0
exit status 2

```
尝试通过反射来修改整型变量失败了，程序直接抛出了异常。下面我们来尝试通过反射来修改指针变量指向的值，这个是可行的。
```
package main

import "fmt"
import "reflect"

func main() {
    var s int = 42
    // 反射指针类型
 var v = reflect.ValueOf(&s)
    // 要拿出指针指向的元素进行修改
 v.Elem().SetInt(43)
    fmt.Println(s)
}

-------
43
```
可以看到变量 s 的值确实被修改成功了，不过这个例子修改的是指针指向的值而不是修改指针变量本身，如果不使用 Elem() 方法进行修改也会抛出一样的异常。

结构体也是值类型，也必须通过指针类型来修改。下面我们尝试使用反射来动态修改结构体内部字段的值。

```
package main

import "fmt"
import "reflect"

type Rect struct {
    Width int
    Height int
}

func SetRectAttr(r *Rect, name string, value int) {
    var v = reflect.ValueOf(r)
    var field = v.Elem().FieldByName(name)
    field.SetInt(int64(value))
}

func main() {
    var r = Rect{50, 100}
    SetRectAttr(&r, "Width", 100)
    SetRectAttr(&r, "Height", 200)
    fmt.Println(r)
}

-----
{100 200}
```



>--摘自 - 《快学 Go 语言》
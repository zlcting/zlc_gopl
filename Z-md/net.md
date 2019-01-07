# **net/http包和高性能可扩展 HTTP 路由 httprouter**

Go语言(golang)的一个很大的优势，就是很容易的开发出网络后台服务，而且性能快，效率高。在开发后端HTTP网络应用服务的时候，我们需要处理很多HTTP的请求访问，比如常见的API服务，我们就要处理很多HTTP请求，然后把处理的信息返回给使用者。对于这类需求，Golang提供了内置的net/http包帮我们来处理这些HTTP请求，让我们可以比较方便的开发一个HTTP服务。

**net/http**

```go
func main() {
	http.HandleFunc("/",Index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Index(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w,"Blog:www.flysnow.org\nwechat:flysnow_org")
}
```
这是net/http包中一个经典的HTTP服务实现，我们运行后打开http://localhost:8080，就可以看到如下信息:

```
Blog:www.flysnow.org
wechat:flysnow_org
```
显示的关键就是我们http.HandleFunc函数，我们通过该函数注册了对路径/的处理函数Index，所有才会看到上面的显示信息。那么这个http.HandleFunc他是如何注册Index函数的呢？下面看看这个函数的源代码。

```go
// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool // whether any patterns contain hostnames
}
```
看以上的源代码，是存在一个默认的DefaultServeMux路由的，这个DefaultServeMux类型是ServeMux,我们注册的路径函数信息都被存入ServeMux的m字段中，以便处理HTTP请求的时候使用。

DefaultServeMux.HandleFunc函数最终会调用ServeMux.Handle函数。

```go
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	//省略加锁和判断代码

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	//把我们注册的路径和相应的处理函数存入了m字段中
	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}
```
这下应该明白了，注册的路径和相应的处理函数都存入了m字段中。

既然注册存入了相应的信息，那么在处理HTTP请求的时候，就可以使用了。Go语言的net/http更底层细节就不详细分析了，我们只要知道处理HTTP请求的时候，会调用Handler接口的ServeHTTP方法，而ServeMux正好实现了Handler。

```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	//省略一些无关代码
	
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}
```
上面代码中的mux.Handler会获取到我们注册的Index函数，然后执行它，具体mux.Handler的详细实现不再分析了，大家可以自己看下源代码。

现在我们可以总结下net/http包对HTTP请求的处理。

>**HTTP请求->ServeHTTP函数->ServeMux的Handler方法->Index函数**

这就是整个一条请求处理链，现在我们明白了net/http里对HTTP请求的原理。

**net/http 的不足**

我们自己在使用内置的net/http的默认路径处理HTTP请求的时候，会发现很多不足，比如：

不能单独的对请求方法(POST,GET等)注册特定的处理函数
不支持Path变量参数
不能自动对Path进行校准
性能一般
扩展性不足
……

那么如何解决以上问题呢？一个办法就是自己写一个处理HTTP请求的路由，因为从上面的源代码我们知道，net/http用的是默认的路径。

```go
//这个是我们启动HTTP服务的函数，最后一个handler参数是nil
http.ListenAndServe(":8080", nil)

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	
	//这个判断成立，因为我们传递的是nil
	if handler == nil {
		handler = DefaultServeMux
	}
	//省略了一些代码
	handler.ServeHTTP(rw, req)
}
```
通过以上的代码分析，我们自己在通过http.ListenAndServe函数启动一个HTTP服务的时候，最后一个handler的值是nil，所以上面的nil判断成立，使用的就是默认的路由DefaultServeMux。

现在我们就知道如何使用自己定义的路由了，那就是给http.ListenAndServe的最后一个参数handler传一个自定义的路由，比如：

```go
type CustomMux struct {

}

func (cm *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,"Blog:www.flysnow.org\nwechat:flysnow_org")
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", &CustomMux{}))
}
```
这个自定义的CustomMux就是我们的路由，它显示了和使用net/http演示的例子一样的功能。

现在我们改变下代码，只有GET方法才会显示以上信息。

```go
func (cm *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w,"Blog:www.flysnow.org\nwechat:flysnow_org")
	} else {
		fmt.Fprint(w,"bad http method request")
	}
}
```

只需要改变下ServeHTTP方法的处理逻辑即可，现在我们可以换不同的请求方法试试，就会显示不同的内容。

这个就是自定义，我们可以通过扩展ServeHTTP方法的实现来添加我们想要的任何功能，包括上面章节列出来的net/http的不足都可以解决，不过我们无需这么麻烦，因为开源大牛已经帮我们做了这些事情，它就是 github.com/julienschmidt/httprouter

**httprouter**

httprouter 是一个高性能、可扩展的HTTP路由，上面我们列举的net/http默认路由的不足，都被httprouter 实现，我们先用一个例子，认识下 httprouter 这个强大的 HTTP 路由。
```go
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Blog:%s \nWechat:%s","www.flysnow.org","flysnow_org")
}
func main() {
	router := httprouter.New()
	router.GET("/", Index)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
这个例子，实现了在GET请求/路径时，会显示如下信息：

```
Blog:www.flysnow.org
wechat:flysnow_org
```
在这个例子中，首先通过httprouter.New()生成了一个*Router路由指针,然后使用GET方法注册一个适配/路径的Index函数，最后*Router作为参数传给ListenAndServe函数启动HTTP服务即可。

其实不止是GET方法，httprouter 为所有的HTTP Method 提供了快捷的使用方式，只需要调用对应的方法即可。

```go
func (r *Router) GET(path string, handle Handle) {
	r.Handle("GET", path, handle)
}

func (r *Router) HEAD(path string, handle Handle) {
	r.Handle("HEAD", path, handle)
}

func (r *Router) OPTIONS(path string, handle Handle) {
	r.Handle("OPTIONS", path, handle)
}

func (r *Router) POST(path string, handle Handle) {
	r.Handle("POST", path, handle)
}

func (r *Router) PUT(path string, handle Handle) {
	r.Handle("PUT", path, handle)
}

func (r *Router) PATCH(path string, handle Handle) {
	r.Handle("PATCH", path, handle)
}

func (r *Router) DELETE(path string, handle Handle) {
	r.Handle("DELETE", path, handle)
}
```
以上这些方法都是 httprouter 支持的，我们可以非常灵活的根据需要，使用对应的方法，这样就解决了net/http默认路由的问题。

**httprouter 命名参数**

现代的API，基本上都是Restful API，httprouter提供的命名参数的支持，可以很方便的帮助我们开发Restful API。比如我们设计的API/user/flysnow，这这样一个URL，可以查看flysnow这个用户的信息，如果要查看其他用户的，比如zhangsan,我们只需要访问API/user/zhangsan即可。
现在我们可以发现，其实这是一种URL匹配模式，我们可以把它总结为/user/:name,这是一个通配符，看个例子。

```go
func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/user/:name",UserInfo)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

当我们运行，在浏览器里输入http://localhost:8080/user/flysnow时，就会显示hello, flysnow!.

通过上面的代码示例，可以看到，路径的参数是以:开头的，后面紧跟着变量名，比如:name，然后在UserInfo这个处理函数中，通过httprouter.Params的ByName获取对应的值。

:name这种匹配模式，是精准匹配的，同时只能匹配一个，比如：
```
Pattern: /user/:name

 /user/gordon              匹配
 /user/you                 匹配
 /user/gordon/profile      不匹配
 /user/                    不匹配
 ```
因为httprouter这个路由就是单一匹配的，所以当我们使用命名参数的时候，一定要注意，是否有其他注册的路由和命名参数的路由，匹配同一个路径，比如/user/new这个路由和/user/:name就是冲突的，不能同时注册。

这里稍微提下httprouter的另外一种通配符模式，就是把:换成*，也就是*name，这是一种匹配所有的模式，不常用，比如：
```
Pattern: /user/*name

 /user/gordon              匹配
 /user/you                 匹配
 /user/gordon/profile      匹配
 /user/                    匹配
```

因为是匹配所有的<font color=#FF0000>  * </font>   模式，所以只要<font color=#FF0000>  * </font>   前面的路径匹配，就是匹配的，不管路径多长，有几层，都匹配。

**httprouter兼容http.Handler**

通过上面的例子，我们应该已经发现，GET方法的handle，并不是我们熟悉的http.Handler，它是httprouter自定义的，相比http.Handler多了一个通配符参数的支持。
```go
type Handle func(http.ResponseWriter, *http.Request, Params)
```
自定义的Handle，唯一的目的就是支持通配符参数，如果你的HTTP服务里，有些路由没有用到通配符参数，那么可以使用原生的http.Handler，httprouter是兼容支持的，这也为我们从net/http的方式，升级为httprouter路由提供了方便，会高效很多。

```go
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, _ Params) {
			handler.ServeHTTP(w, req)
		},
	)
}

func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}
```

httprouter通过Handler和HandlerFunc两个函数，提供了兼容http.Handler和http.HandlerFunc的完美支持。从以上源代码中，我们可以看出，实现的方式也比较简单，就是做了一个http.Handler到httprouter.Handle的转换，舍弃了通配符参数的支持。

**Handler处理链**

得益于http.Handler的模式，我们可以把不同的http.Handler组成一个处理链，httprouter.Router也是实现了http.Handler的，所以它也可以作为http.Handler处理链的一部分，比如和Negroni、Gorilla handlers这两个库配合使用，关于这两个库的介绍，可以参考我以前写的文章。

[Go语言经典库使用分析（五）| Negroni 中间件（一）](https://www.flysnow.org/2017/08/20/go-classic-libs-negroni-one.html)

[ Go语言经典库使用分析（三）| Gorilla Handlers 详细介绍](https://www.flysnow.org/2017/08/06/go-classic-libs-gorilla-handlers-guide.html)
这里使用一个官方的例子，作为Handler处理链的演示。

比如对多个不同的二级域名，进行不同的路由处理。

```go
//一个新类型，用于存储域名对应的路由
type HostSwitch map[string]http.Handler

//实现http.Handler接口，进行不同域名的路由分发
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    //根据域名获取对应的Handler路由，然后调用处理（分发机制）
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func main() {
    //声明两个路由
	playRouter := httprouter.New()
	playRouter.GET("/", PlayIndex)
	
	toolRouter := httprouter.New()
	toolRouter.GET("/", ToolIndex)

    //分别用于处理不同的二级域名
	hs := make(HostSwitch)
	hs["play.flysnow.org:12345"] = playRouter
	hs["tool.flysnow.org:12345"] = toolRouter

    //HostSwitch实现了http.Handler,所以可以直接用
	log.Fatal(http.ListenAndServe(":12345", hs))
}
```
以上就是一个简单的，针对不同域名，使用不同路由的例子，代码中的注释比较详细了，这里就不一一解释了。这个例子中,HostSwitch和httprouter.Router这两个http.Handler就组成了一个http.Handler处理链。

**httprouter 静态文件服务**
httprouter提供了很方便的静态文件服务，可以把一个目录托管在服务器上，以供访问。
```go
router.ServeFiles("/static/*filepath",http.Dir("./"))
```
只需要这一句核心代码即可，这个就是把当前目录托管在服务器上，以供访问，访问路径是/static。

使用ServeFiles需要注意的是，第一个参数路径，必须要以/*filepath，因为要获取我们要访问的路径信息。

```go
func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	r.GET(path, func(w http.ResponseWriter, req *http.Request, ps Params) {
		req.URL.Path = ps.ByName("filepath")
		fileServer.ServeHTTP(w, req)
	})
}
```
**httprouter 异常捕获**
很少有路由支持这个功能的，httprouter允许使用者，设置PanicHandler用于处理HTTP请求中发生的panic。
```go
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("故意抛出的异常")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s",v)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
演示例子中，我们通过设置router.PanicHandler来处理发生的panic，处理办法是打印出来异常信息。然后故意在Index函数中抛出一个painc，然后我们运行测试，会看到异常信息。

这是一种非常好的方式，可以让我们对painc进行统一处理，不至于因为漏掉的panic影响用户使用。

**小结**
httprouter还有不少有用的小功能，比如对404进行处理，我们通过设置Router.NotFound来实现，我们看看Router这个结构体的配置，可以发现更多有用的功能。
```go
type Router struct {
    //是否通过重定向，给路径自定加斜杠
	RedirectTrailingSlash bool
    //是否通过重定向，自动修复路径，比如双斜杠等自动修复为单斜杠
	RedirectFixedPath bool
    //是否检测当前请求的方法被允许
	HandleMethodNotAllowed bool
	//是否自定答复OPTION请求
	HandleOPTIONS bool
    //404默认处理
	NotFound http.Handler
    //不被允许的方法默认处理
	MethodNotAllowed http.Handler
    //异常统一处理
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}
```
这些字段都是导出的(export)，我们可以直接设置，来达到我们的目的。

httprouter是一个高性能，低内存占用的路由，它使用radix tree实现存储和匹配查找，所以效率非常高，内存占用也很低。关于radix tree大家可以查看相关的资料。

httprouter因为实现了http.Handler，所以可扩展性非常好，可以和其他库、中间件结合使用，gin这个web框架就是使用的自定义的httprouter。


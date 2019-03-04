# **轻量级 Web 框架 Gin 结构分析**
Go 语言最流行了两个轻量级 Web 框架分别是 Gin 和 Echo，这两个框架大同小异，都是插件式轻量级框架，背后都有一个开源小生态来提供各式各样的小插件，这两个框架的性能也都非常好，裸测起来跑的飞快。本节我们只讲 Gin 的实现原理和使用方法，Gin 起步比 Echo 要早，市场占有率要高一些，生态也丰富一些。
```go
go get -u github.com/gin-gonic/gin
```
# **Hello World**
Gin 框架的 Hello World 只需要 10 行代码，比大多数动态脚本语言稍微多几行。
```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

代码中的 gin.H 是 map[string]interface{} 的一个快捷名称，写起来会更加简洁。
```go
type H map[string]interface{}
```
# **gin.Engine**
Engine 是 Gin 框架最重要的数据结构，它是框架的入口。我们通过 Engine 对象来定义服务路由信息、组装插件、运行服务。正如 Engine 的中文意思「引擎」一样，它就是框架的核心发动机，整个 Web 服务的都是由它来驱动的。

发动机属于精密设备，构造非常复杂，不过 Engine 对象很简单，因为引擎最重要的部分 —— 底层的 HTTP 服务器使用的是 Go 语言内置的 http server，Engine 的本质只是对内置的 HTTP 服务器的包装，让它使用起来更加便捷。

gin.Default() 函数会生成一个默认的 Engine 对象，里面包含了 2 个默认的常用插件，分别是 Logger 和 Recovery，Logger 用于输出请求日志，Recovery 确保单个请求发生 panic 时记录异常堆栈日志，输出统一的错误响应。

```go
func Default() *Engine {
    engine := New()
    engine.Use(Logger(), Recovery())
    return engine
}
```
# **路由树**
在 Gin 框架中，路由规则被分成了最多 9 棵前缀树，每一个 HTTP Method对应一棵「前缀树」，树的节点按照 URL 中的 / 符号进行层级划分，URL 支持 :name 形式的名称匹配，还支持 *subpath 形式的路径通配符 。

```go
// 匹配单节点 named
pattern = /book/:id
match /book/123
nomatch /book/123/10
nomatch /book/

// 匹配子节点 catchAll mode
/book/*subpath
match /book/
match /book/123
match /book/123/10
```
![图片](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mIB0nOGicskI7oBydr7uhPrQdPhP6iar5zvXNN1Ya37tyakC6xohvty9NiaDvdUBcEclvKzb2aRTMMyQ/640?wx_fmt=png)


每个节点都会挂接若干请求处理函数构成一个请求处理链 HandlersChain。当一个请求到来时，在这棵树上找到请求 URL 对应的节点，拿到对应的请求处理链来执行就完成了请求的处理。

```go
type Engine struct {
  ...
  trees methodTrees
  ...
}

type methodTrees []methodTree

type methodTree struct {
    method string
    root   *node  // 树根
}

type node struct {
  path string // 当前节点的路径
  ...
  handlers HandlersChain // 请求处理链
  ...
}

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc
```
Engine 对象包含一个 addRoute 方法用于添加 URL 请求处理器，它会将对应的路径和处理器挂接到相应的请求树中

```go
func (e *Engine) addRoute(method, path string, handlers HandlersChain)
```
# **gin.RouterGroup** 
RouterGroup 是对路由树的包装，所有的路由规则最终都是由它来进行管理。Engine 结构体继承了 RouterGroup ，所以 Engine 直接具备了 RouterGroup 所有的路由管理功能。这是为什么在 Hello World 的例子中，可以直接使用 Engine 对象来定义路由规则。同时 RouteGroup 对象里面还会包含一个 Engine 的指针，这样 Engine 和 RouteGroup 就成了「你中有我我中有你」的关系。
```go
type Engine struct {
  RouterGroup
  ...
}

type RouterGroup struct {
  ...
  engine *Engine
  ...
}
```
RouterGroup 实现了 IRouter 接口，暴露了一系列路由方法，这些方法最终都是通过调用 Engine.addRoute 方法将请求处理器挂接到路由树中。
```go
GET(string, ...HandlerFunc) IRoutes
POST(string, ...HandlerFunc) IRoutes
DELETE(string, ...HandlerFunc) IRoutes
PATCH(string, ...HandlerFunc) IRoutes
PUT(string, ...HandlerFunc) IRoutes
OPTIONS(string, ...HandlerFunc) IRoutes
HEAD(string, ...HandlerFunc) IRoutes
// 匹配所有 HTTP Method
Any(string, ...HandlerFunc) IRoutes
```
RouterGroup 内部有一个前缀路径属性，它会将所有的子路径都加上这个前缀再放进路由树中。有了这个前缀路径，就可以实现 URL 分组功能。Engine 对象内嵌的 RouterGroup 对象的前缀路径是 /，它表示根路径。RouterGroup 支持分组嵌套，使用 Group 方法就可以让分组下面再挂分组，于是子子孙孙无穷尽也。
```go
func main() {
    router := gin.Default()

    v1 := router.Group("/v1")
    {
        v1.POST("/login", loginEndpoint)
        v1.POST("/submit", submitEndpoint)
        v1.POST("/read", readEndpoint)
    }

    v2 := router.Group("/v2")
    {
        v2.POST("/login", loginEndpoint)
        v2.POST("/submit", submitEndpoint)
        v2.POST("/read", readEndpoint)
    }

    router.Run(":8080")
}
```
上面这个例子中实际上已经使用了分组嵌套，因为 Engine 对象里面的 RouterGroup 对象就是第一层分组，也就是根分组，v1 和 v2 都是根分组的子分组。

# **gin.Context**
这个对象里保存了请求的上下文信息，它是所有请求处理器的入口参数。
```go
type HandlerFunc func(*Context)

type Context struct {
  ...
  Request *http.Request // 请求对象
  Writer ResponseWriter // 响应对象
  Params Params // URL匹配参数
  ...
  Keys map[string]interface{} // 自定义上下文信息
  ...
}
```
Context 对象提供了非常丰富的方法用于获取当前请求的上下文信息，如果你需要获取请求中的 URL 参数、Cookie、Header 都可以通过 Context 对象来获取。这一系列方法本质上是对 http.Request 对象的包装。

```go
// 获取 URL 匹配参数  /book/:id
func (c *Context) Param(key string) string
// 获取 URL 查询参数 /book?id=123&page=10
func (c *Context) Query(key string) string
// 获取 POST 表单参数
func (c *Context) PostForm(key string) string
// 获取上传的文件对象
func (c *Context) FormFile(name string) (*multipart.FileHeader, error)
// 获取请求Cookie
func (c *Context) Cookie(name string) (string, error) 
...
```
Context 对象提供了很多内置的响应形式，JSON、HTML、Protobuf 、MsgPack、Yaml 等。它会为每一种形式都单独定制一个渲染器。通常这些内置渲染器已经足够应付绝大多数场景，如果你觉得不够，还可以自定义渲染器。

```go
func (c *Context) JSON(code int, obj interface{})
func (c *Context) Protobuf(code int, obj interface{})
func (c *Context) YAML(code int, obj interface{})
...
// 自定义渲染
func (c *Context) Render(code int, r render.Render)

// 渲染器通用接口
type Render interface {
    Render(http.ResponseWriter) error
    WriteContentType(w http.ResponseWriter)
}
```
所有的渲染器最终还是需要调用内置的 http.ResponseWriter（Context.Writer） 将响应对象转换成字节流写到套接字中。

```go
type ResponseWriter interface {
 // 容纳所有的响应头
 Header() Header
 // 写Body
 Write([]byte) (int, error)
 // 写Header
 WriteHeader(statusCode int)
}
```
# **插件与请求链**
我们编写业务代码时一般也就是一个处理函数，为什么路由节点需要挂接一个函数链呢？

```go
type node struct {
  path string // 当前节点的路径
  ...
  handlers HandlersChain // 请求处理链
  ...
}
type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc
```
这是因为 Gin 提供了插件，只有函数链的尾部是业务处理，前面的部分都是插件函数。在 Gin 中插件和业务处理函数形式是一样的，都是 func(*Context)。当我们定义路由时，Gin 会将插件函数和业务处理函数合并在一起形成一个链条结构。

```go
type Context struct {
  ...
  index uint8 // 当前的业务逻辑位于函数链的位置
  handlers HandlersChain // 函数链
  ...
}

// 挨个调用链条中的处理函数
func (c *Context) Next() {
    c.index++
    for s := int8(len(c.handlers)); c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}
```
Gin 在接收到客户端请求时，找到相应的处理链，构造一个 Context 对象，再调用它的 Next() 方法就正式进入了请求处理的全流程。
![图片](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mIB0nOGicskI7oBydr7uhPrQXWqyNaKCHcj5CmZN3K7PpWqaCgOSwPoywNfoLXrAI45kIABdMNAuSQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)
Gin 还支持 Abort() 方法中断请求链的执行，它的原理是将 Context.index 调整到一个比较大的数字，这样 Next() 方法中的调用循环就会立即结束。需要注意的 Abort() 方法并不是通过 panic 的方式中断执行流，执行 Abort() 方法之后，当前函数内后面的代码逻辑还会继续执行。

```go
const abortIndex = 127
func (c *Context) Abort() {
    c.index = abortIndex
}

func SomePlugin(c *Context) {
  ...
  if condition {
    c.Abort()
    // continue executing
  }
  ...
}
```
如果在插件中显示调用 Next() 方法，那么它就改变了正常的顺序执行流，变成了像洋葱一样的嵌套执行流。换个角度来理解，正常的执行流就是后续的处理器是在前一个处理器的尾部执行，而嵌套执行流是让后续的处理器在前一个处理器进行到一半的时候执行，待后续处理器完成执行后，再回到前一个处理器继续往下执行。
![图片](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mIB0nOGicskI7oBydr7uhPrQLdDicBrIiaErmghGicHkvj8LBNBjLzArE2NbUX6XHFfZueISeFM9LvQdQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)


要是你学过 Python 语言，这种嵌套结构很容易让人联想到装饰器 decorator。如果所有的插件都使用嵌套执行流，那么就会变成了下面这张图
![图片](https://mmbiz.qpic.cn/mmbiz_png/bGribGtYC3mIB0nOGicskI7oBydr7uhPrQoa6seZQUbWR85MFbTub960mQot4I1royA1B0g2nvueFAooaLABSMgQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)
RouterGroup 提供了 Use() 方法来注册插件，因为 RouterGroup 是一层套一层，不同层级的路由可能会注册不一样的插件，最终不同的路由节点挂接的处理函数链也不尽相同。
```go
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
    group.Handlers = append(group.Handlers, middleware...)
    return group.returnObj()
}

// 注册 Get 请求
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
    return group.handle("GET", relativePath, handlers)
}

func (g *RouterGroup) handle(method, path string, handlers HandlersChain) IRoutes {
 // 合并URL (RouterGroup有URL前缀)
    absolutePath := group.calculateAbsolutePath(relativePath)
    // 合并处理链条
 handlers = group.combineHandlers(handlers)
    // 注册路由树
 group.engine.addRoute(httpMethod, absolutePath, handlers)
    return group.returnObj()
}
```
# **HTTP 错误**
当 URL 请求对应的路径不能在路由树里找到时，就需要处理 404 NotFound 错误。当 URL 的请求路径可以在路由树里找到，但是 Method 不匹配，就需要处理 405 MethodNotAllowed 错误。Engine 对象为这两个错误提供了处理器注册的入口
```go
func (engine *Engine) NoMethod(handlers ...HandlerFunc)
func (engine *Engine) NoRoute(handlers ...HandlerFunc)
```
异常处理器和普通处理器一样，也需要和插件函数组合在一起形成一个调用链。如果没有提供异常处理器，Gin 就会使用内置的简易错误处理器。

注意这两个错误处理器是定义在 Engine 全局对象上，而不是 RouterGroup。对于非 404 和 405 错误，需要用户自定义插件来处理。对于 panic 抛出来的异常需要也需要使用插件来处理。
# **静态文件服务**
RouterGroup 对象里定义了下面三个用来服务静态文件的方法
```go
// 服务单个静态文件
StaticFile(relativePath, filePath string) IRoutes
// 服务静态文件目录
Static(relativePath, dirRoot string) IRoutes
// 服务虚拟静态文件系统
StaticFS(relativePath string, fs http.FileSystem) IRoutes
```
它不同于错误处理器，静态文件服务挂在 RouterGroup 上，支持嵌套。这三个方法中 StaticFS 方法比较特别，它对文件系统进行了抽象，你可以提供一个基于网络的静态文件系统，也可以提供一个基于内存的静态文件系统。FileSystem 接口也很简单，提供一个路径参数返回一个实现了 File 接口的文件对象。不同的虚拟文件系统使用不同的代码来实现 File 接口。

```go
type FileSystem interface {
 Open(path string) (File, error)
}

type File interface {
 io.Closer
 io.Reader
 io.Seeker
 Readdir(count int) ([]os.FileInfo, error)
 Stat() (os.FileInfo, error)
}
```
静态文件处理器和普通处理器一样，也需要经过插件的重重过滤。
# **表单处理**
当请求参数数量比较多时，使用 Context.Query() 和 Context.PostForm() 方法来获取参数就会显得比较繁琐。Gin 框架也支持表单处理，将表单参数和结构体字段进行直接映射。
```go
package main

import (
    "github.com/gin-gonic/gin"
)

type LoginForm struct {
    User     string `form:"user" binding:"required"`
    Password string `form:"password" binding:"required"`
}

func main() {
    router := gin.Default()
    router.POST("/login", func(c *gin.Context) {
  var form LoginForm
        if c.ShouldBind(&form) == nil {
            if form.User == "user" && form.Password == "password" {
                c.JSON(200, gin.H{"status": "you are logged in"})
            } else {
                c.JSON(401, gin.H{"status": "unauthorized"})
            }
        }
    })
    router.Run(":8080")
}
```

Context.ShouldBind 方法遇到校验不通过时，会返回一个错误对象告知调用者校验失败的原因。它支持多种数据绑定类型，如 XML、JSON、Query、Uri、MsgPack、Protobuf等，根据请求的 Content-Type 头来决定使用何种数据绑定方法。

```go
func (c *Context) ShouldBind(obj interface{}) error {
 // 获取绑定器
    b := binding.Default(c.Request.Method, c.ContentType())
    // 执行绑定
 return c.ShouldBindWith(obj, b)
}
```
默认内置的表单校验功能很强大，它通过结构体字段 tag 标注来选择相应的校验器进行校验。Gin 还提供了注册自定义校验器的入口，支持用户自定义一些通用的特殊校验逻辑。

Context.ShouldBind 是比较柔和的校验方法，它只负责校验，并将校验结果以返回值的形式传递给上层。Context 还有另外一个比较暴力的校验方法 Context.Bind，它和 ShouldBind 的调用形式一摸一样，区别是当校验错误发生时，它会调用 Abort() 方法中断调用链的执行，向客户端返回一个 HTTP 400 Bad Request 错误。
# **HTTPS**
Gin 不支持 HTTPS，官方建议是使用 Nginx 来转发 HTTPS 请求到 Gin。


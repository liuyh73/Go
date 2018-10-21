[TOC]

# struct

```go
package main

import (
	"fmt"
)

type parent struct {
	a, b int
}

type child struct {
	c, d int
	*parent
}

func main() {
	child_1 := child{1, 2, &parent{3, 4}}
	fmt.Println(child_1.a, child_1.d, child_1.parent.a, child_1.parent.b)
	// fmt.Printf("%#v\n", child_1)
}

// 输出结果：
// 3 2 3 4
```

​	struct利用type来表示自定义类型，struct可以继承（暂且这样说吧），在上面的例子中，child继承了parent的所有字段（a，b），parent作为了child的一个你名字段，在child实例中可以直接用成员操作符来访问a，b，也可以使用child.parent.a来访问属性。这就引出了这样一种情况：当parent和child包含同名字段e，此时child.e访问的即为child内部的字段，child.parent.e访问的即为parent内部的字段。

关于结构体定义和初始化的一些讨论：

```go
package main

import (
	"fmt"
)
// 此方法定义的结构体只可以用一次
var s1 = struct {
	name string
	age  int
}{
	name: "Sam",
	age:  18,
}

type student2 struct {
	name string
	age  int
}

func main() {
	s2 := student2{"Kate", 19}
	fmt.Printf("%v\n%v\n", s1, s2)
}
// PS E:\mygo\src> go run struct2.go
// {Sam 18}
// {Kate 19}
```

# interface

interface是一组method签名的组合，我们通过interface来定义对象的一组行为。 任意的类型都实现了空interface(我们这样定义：interface{})，也就是包含0个method的interface。 

- interface类型

```go
package main

import (
	"fmt"
	"math"
)
// 自定义类型type
type Abser interface {
	Abs() float64
}
// MyFloat也是自定义类型，不是float64的别名
type MyFloat float64
// MyFloat实现 Abs方法，则该自定义类型实现了Abser接口
func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}
// Vertex*实现了Abs方法，表示自定义类型Vertex的指针实现了Abser接口
func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := &Vertex{3, 4}

	a = f
	fmt.Println(a.Abs())
	a = v
	fmt.Println(a.Abs())
}

// 输出结果
// PS E:\mygo\src> go run interface.go
// 1.4142135623730951
// 5
```

- interface值

如果我们定义了一个interface的变量，那么这个变量里面可以存实现这个interface的任意类型的对象。如上述代码中，a中可以存储MyFloat类型的变量，也可以存储&Vertex类型的变量。

- 空interface

空interface(interface{})不包含任何的method，正因为如此，所有的类型都实现了空interface。空interface对于描述起不到任何的作用(因为它不包含任何的method），但是空interface在我们需要存储任意类型的数值的时候相当有用，因为它可以存储任意类型的数值。它有点类似于C语言的void*类型。 

```go
// 定义a为空接口
var a interface{}
var i int = 5
s := "Hello world"
// a可以存储任意类型的数值
a = i
a = s
```

- interface函数参数

interface的变量可以持有任意实现该interface类型的对象，这给我们编写函数(包括method)提供了一些额外的思考，我们是不是可以通过定义interface参数，让函数接受各种类型的参数。 比如fmt包中的Stringer接口：（详情在后续的stringers中）

```go
type Stringer interface {
	 String() string
}
```

- interface变量存储的类型

  判断interface变量里面存储的熟知的类型的方法：

  - Comma-ok断言：

    Go语言里面有一个语法，可以直接判断是否是该类型的变量： value, ok = element.(T)，这里value就是变量的值，ok是一个bool类型，element是interface变量，T是断言的类型。

    如果element里面确实存储了T类型的数值，那么ok返回true，否则返回false。

    ```go
    package main
    
    	import (
    		"fmt"
    		"strconv"
    	)
    
    	type Element interface{}
    	type List [] Element
    
    	type Person struct {
    		name string
    		age int
    	}
    
    	//定义了String方法，实现了fmt.Stringer
    	func (p Person) String() string {
    		return "(name: " + p.name + " - age: "+strconv.Itoa(p.age)+ " years)"
    	}
    
    	func main() {
    		list := make(List, 3)
    		list[0] = 1 // an int
    		list[1] = "Hello" // a string
    		list[2] = Person{"Dennis", 70}
    
    		for index, element := range list {
    			if value, ok := element.(int); ok {
    				fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
    			} else if value, ok := element.(string); ok {
    				fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
    			} else if value, ok := element.(Person); ok {
    				fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
    			} else {
    				fmt.Printf("list[%d] is of a different type\n", index)
    			}
    		}
    	}
    ```

  - switch测试：

    ```go
    package main
    
    	import (
    		"fmt"
    		"strconv"
    	)
    
    	type Element interface{}
    	type List [] Element
    
    	type Person struct {
    		name string
    		age int
    	}
    
    	//打印
    	func (p Person) String() string {
    		return "(name: " + p.name + " - age: "+strconv.Itoa(p.age)+ " years)"
    	}
    
    	func main() {
    		list := make(List, 3)
    		list[0] = 1 //an int
    		list[1] = "Hello" //a string
    		list[2] = Person{"Dennis", 70}
    
    		for index, element := range list{
    			switch value := element.(type) {
    				case int:
    					fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
    				case string:
    					fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
    				case Person:
    					fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
    				default:
    					fmt.Println("list[%d] is of a different type", index)
    			}
    		}
    	}
    ```

  这里有一点需要强调的是：`element.(type)`语法不能在switch外的任何逻辑里面使用，如果你要在switch外面判断一个类型就使用`comma-ok`。

- 嵌入interface

  如果一个interface1作为interface2的一个嵌入字段，那么interface2隐式的包含了interface1里面的method。 

  我们可以看到源码包container/heap里面有这样的一个定义

  ```go
  type Interface interface {
  	sort.Interface //嵌入字段sort.Interface
  	Push(x interface{}) //a Push method to push elements into the heap
  	Pop() interface{} //a Pop elements that pops elements from the heap
  }
  ```

  我们看到sort.Interface其实就是嵌入字段，把sort.Interface的所有method给隐式的包含进来了。也就是下面三个方法：

  ```go
  type Interface interface {
  	// Len is the number of elements in the collection.
  	Len() int
  	// Less returns whether the element with index i should sort
  	// before the element with index j.
  	Less(i, j int) bool
  	// Swap swaps the elements with indexes i and j.
  	Swap(i, j int)
  }
  ```

# errors

error 包实现了用于错误处理的函数. 

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
    // New 返回一个按给定文本格式化的错误。
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Println(err)
	}
	// fmt 包的 Errorf 函数让我们使用该包的格式化特性来创建描述性的错误信息。
	const name, id = "bimmler", 17
	err = fmt.Errorf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Println(err)
	}
}
// 输出结果
// PS E:\mygo\src> go run error.go
// CreateFile error.go: The system cannot find the file specified.
```

# http

- web使用示例：

  ```go
  package main
  
  import (
  	"fmt"
  	"net/http"
  	"strings"
  	"log"
  )
  
  func sayhelloName(w http.ResponseWriter, r *http.Request) {
  	r.ParseForm()  //解析参数，默认是不会解析的
  	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
  	fmt.Println("path", r.URL.Path)
  	fmt.Println("scheme", r.URL.Scheme)
  	fmt.Println(r.Form["url_long"])
  	for k, v := range r.Form {
  		fmt.Println("key:", k)
  		fmt.Println("val:", strings.Join(v, ""))
  	}
  	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
  }
  
  func main() {
  	http.HandleFunc("/", sayhelloName) //设置访问的路由
  	err := http.ListenAndServe(":9090", nil) //设置监听的端口
  	if err != nil {
  		log.Fatal("ListenAndServe: ", err)
  	}
  }
  ```

- 服务器端概念：

  Request： 用户请求的信息，用来解析用户请求，包括get、post、cookie、url等信息。

  Response：服务器需要反馈给客户端的信息

  Conn：用户的每次请求链接

  Handler：处理请求和生成返回信息的处理逻辑

- http包运行机制：

  如下图所示，是Go实现Web服务的工作模式的流程图 

  ![img](https://raw.githubusercontent.com/astaxie/build-web-application-with-golang/master/zh/images/3.3.http.png) 

  图3.9 http包执行流程

  1. 创建Listen Socket, 监听指定的端口, 等待客户端请求到来。
  2. Listen Socket接受客户端的请求, 得到Client Socket, 接下来通过Client Socket与客户端通信。
  3. 处理客户端的请求, 首先从Client Socket读取HTTP请求的协议头, 如果是POST方法, 还可能要读取客户端提交的数据, 然后交给相应的handler处理请求, handler处理完毕准备好客户端需要的数据, 通过Client Socket写给客户端。

  这整个的过程里面我们只要了解清楚下面三个问题，也就知道Go是如何让Web运行起来了

  - 如何监听端口？
  - 如何接收客户端请求？
  - 如何分配handler？

  前面小节的代码里面我们可以看到，Go是通过一个函数`ListenAndServe`来处理这些事情的，这个底层其实这样处理的：初始化一个server对象，然后调用了`net.Listen("tcp", addr)`，也就是底层用TCP协议搭建了一个服务，然后监控我们设置的端口。

  下面代码来自Go的http包的源码，通过下面的代码我们可以看到整个的http处理过程：

  ```go
  func (srv *Server) Serve(l net.Listener) error {
  	defer l.Close()
  	var tempDelay time.Duration // how long to sleep on accept failure
  	for {
  		rw, e := l.Accept()
  		if e != nil {
  			if ne, ok := e.(net.Error); ok && ne.Temporary() {
  				if tempDelay == 0 {
  					tempDelay = 5 * time.Millisecond
  				} else {
  					tempDelay *= 2
  				}
  				if max := 1 * time.Second; tempDelay > max {
  					tempDelay = max
  				}
  				log.Printf("http: Accept error: %v; retrying in %v", e, tempDelay)
  				time.Sleep(tempDelay)
  				continue
  			}
  			return e
  		}
  		tempDelay = 0
  		c, err := srv.newConn(rw)
  		if err != nil {
  			continue
  		}
  		go c.serve()
  	}
  }
  ```

  监控之后如何接收客户端的请求呢？上面代码执行监控端口之后，调用了`srv.Serve(net.Listener)`函数，这个函数就是处理接收客户端的请求信息。这个函数里面起了一个`for{}`，首先通过Listener接收请求，其次创建一个Conn，最后单独开了一个goroutine，把这个请求的数据当做参数扔给这个conn去服务：`go c.serve()`。这个就是高并发体现了，用户的每一次请求都是在一个新的goroutine去服务，相互不影响。

  那么如何具体分配到相应的函数来处理请求呢？conn首先会解析request:`c.readRequest()`,然后获取相应的handler:`handler := c.server.Handler`，也就是我们刚才在调用函数`ListenAndServe`时候的第二个参数，我们前面例子传递的是nil，也就是为空，那么默认获取`handler = DefaultServeMux`,那么这个变量用来做什么的呢？对，这个变量就是一个路由器，它用来匹配url跳转到其相应的handle函数，那么这个我们有设置过吗?有，我们调用的代码里面第一句不是调用了`http.HandleFunc("/", sayhelloName)`嘛。这个作用就是注册了请求`/`的路由规则，当请求uri为"/"，路由就会转到函数sayhelloName，DefaultServeMux会调用ServeHTTP方法，这个方法内部其实就是调用sayhelloName本身，最后通过写入response的信息反馈到客户端。

  详细的整个流程如下图所示： 

  ![img](https://github.com/astaxie/build-web-application-with-golang/raw/master/zh/images/3.3.illustrator.png?raw=true) 

  上述分析回答了之前的三个问题。

  【补充】当ListenAndServe的第二个参数不为nil时（外部实现的路由器）：

  ```go
  package main
  
  import (
  	"fmt"
  	"log"
  	"net/http"
  )
  
  type Hello struct{}
  
  func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  	fmt.Fprint(w, "Hello!")
  }
  
  func main() {
  	var h Hello
  	err := http.ListenAndServe("localhost:3000", h)
  	if err != nil {
  		log.Fatal(err)
  	}
  }
  ```

- Go的http包详解

  Go的http有两个核心功能：Conn、ServeMux 

  - Conn的goroutine

    与我们一般编写的http服务器不同, Go为了实现高并发和高性能, 使用了goroutines来处理Conn的读写事件, 这样每个请求都能保持独立，相互不会阻塞，可以高效的响应网络事件。这是Go高效的保证。

    Go在等待客户端请求里面是这样写的：

    ```go
    c, err := srv.newConn(rw)
    if err != nil {
    	continue
    }
    go c.serve()
    ```

    这里我们可以看到客户端的每次请求都会创建一个Conn，这个Conn里面保存了该次请求的信息，然后再传递到对应的handler，该handler中便可以读取到相应的header信息，这样保证了每个请求的独立性。

  - ServeMux的自定义

    我们前面小节讲述conn.server的时候，其实内部是调用了http包默认的路由器，通过路由器把本次请求的信息传递到了后端的处理函数。那么这个路由器是怎么实现的呢？

    它的结构如下：

    ```go
    type ServeMux struct {
    	mu sync.RWMutex   //锁，由于请求涉及到并发处理，因此这里需要一个锁机制
    	m  map[string]muxEntry  // 路由规则，一个string对应一个mux实体，这里的string就是注册的路由表达式
    	hosts bool // 是否在任意的规则中带有host信息
    }
    ```

    下面看一下muxEntry

    ```go
    type muxEntry struct {
    	explicit bool   // 是否精确匹配
    	h        Handler // 这个路由表达式对应哪个handler
    	pattern  string  //匹配字符串
    }
    ```

    接着看一下Handler的定义

    ```go
    type Handler interface {
    	ServeHTTP(ResponseWriter, *Request)  // 路由实现器
    }
    ```

    Handler是一个接口，但是前一小节中的`sayhelloName`函数并没有实现ServeHTTP这个接口，为什么能添加呢？原来在http包里面还定义了一个类型`HandlerFunc`,我们定义的函数`sayhelloName`就是这个HandlerFunc调用之后的结果，这个类型默认就实现了ServeHTTP这个接口，即我们调用了HandlerFunc(f),强制类型转换f成为HandlerFunc类型，这样f就拥有了ServeHTTP方法。

    ```go
    type HandlerFunc func(ResponseWriter, *Request)
    
    // ServeHTTP calls f(w, r).
    func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    	f(w, r)
    }
    ```

    路由器里面存储好了相应的路由规则之后，那么具体的请求又是怎么分发的呢？请看下面的代码，默认的路由器实现了`ServeHTTP`：

    ```go
    func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
    	if r.RequestURI == "*" {
    		w.Header().Set("Connection", "close")
    		w.WriteHeader(StatusBadRequest)
    		return
    	}
    	h, _ := mux.Handler(r)
    	h.ServeHTTP(w, r)
    }
    ```

    如上所示路由器接收到请求之后，如果是`*`那么关闭链接，不然调用`mux.Handler(r)`返回对应设置路由的处理Handler，然后执行`h.ServeHTTP(w, r)`

    也就是调用对应路由的handler的ServerHTTP接口，那么mux.Handler(r)怎么处理的呢？

    ```go
    func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
    	if r.Method != "CONNECT" {
    		if p := cleanPath(r.URL.Path); p != r.URL.Path {
    			_, pattern = mux.handler(r.Host, p)
    			return RedirectHandler(p, StatusMovedPermanently), pattern
    		}
    	}	
    	return mux.handler(r.Host, r.URL.Path)
    }
    
    func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
    	mux.mu.RLock()
    	defer mux.mu.RUnlock()
    
    	// Host-specific pattern takes precedence over generic ones
    	if mux.hosts {
    		h, pattern = mux.match(host + path)
    	}
    	if h == nil {
    		h, pattern = mux.match(path)
    	}
    	if h == nil {
    		h, pattern = NotFoundHandler(), ""
    	}
    	return
    }
    ```

    原来他是根据用户请求的URL和路由器里面存储的map去匹配的，当匹配到之后返回存储的handler，调用这个handler的ServeHTTP接口就可以执行到相应的函数了。

    通过上面这个介绍，我们了解了整个路由过程，Go其实支持外部实现的路由器 `ListenAndServe`的第二个参数就是用以配置外部路由器的，它是一个Handler接口，即外部路由器只要实现了Handler接口就可以,我们可以在自己实现的路由器的ServeHTTP里面实现自定义路由功能。

    如下代码所示，我们自己实现了一个简易的路由器

    ```go
    package main
    
    import (
    	"fmt"
    	"net/http"
    )
    
    type MyMux struct {
    }
    
    func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    	if r.URL.Path == "/" {
    		sayhelloName(w, r)
    		return
    	}
    	http.NotFound(w, r)
    	return
    }
    
    func sayhelloName(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintf(w, "Hello myroute!")
    }
    
    func main() {
    	mux := &MyMux{}
    	http.ListenAndServe(":9090", mux)
    }
    ```

- Go代码的执行流程

  通过对http包的分析之后，现在让我们来梳理一下整个的代码执行过程。

  - 首先调用Http.HandleFunc

    按顺序做了几件事：

    1 调用了DefaultServeMux的HandleFunc

    2 调用了DefaultServeMux的Handle

    3 往DefaultServeMux的map[string]muxEntry中增加对应的handler和路由规则

  - 其次调用http.ListenAndServe(":9090", nil)

    按顺序做了几件事情：

    1 实例化Server

    2 调用Server的ListenAndServe()

    3 调用net.Listen("tcp", addr)监听端口

    4 启动一个for循环，在循环体中Accept请求

    5 对每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()

    6 读取每个请求的内容w, err := c.readRequest()

    7 判断handler是否为空，如果没有设置handler（这个例子就没有设置handler），handler就设置为DefaultServeMux

    8 调用handler的ServeHttp

    9 在这个例子中，下面就进入到DefaultServeMux.ServeHttp

    10 根据request选择handler，并且进入到这个handler的ServeHTTP

    ```
      mux.handler(r).ServeHTTP(w, r)
    ```

    11 选择handler：

    A 判断是否有路由能满足这个request（循环遍历ServeMux的muxEntry）

    B 如果有路由满足，调用这个路由handler的ServeHTTP

    C 如果没有路由满足，调用NotFoundHandler的ServeHTTP

# select

golang在语言级别直接支持select，用于处理异步IO问题。 

`select`默认是阻塞的，只有当监听的channel中有发送或接收可以进行时才会运行，当多个channel都准备好的时候，select是随机的选择一个执行的。 

```go
package main

import "fmt"
// 主线程执行fibnacci函数，在select语句中，case c<-x: 尝试向channel c中写入数据，若成功写入，则执行x,y=y,x+y, 否则进行case <-quit: 判断，若quit可以读取数据，则进行return操作
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			//fmt.Printf("%d", c)
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

//goroutine 执行输出斐波那契数列的操作，由于channel c只有一个缓冲区，所以上面函数没写入一个数字，便进入阻塞状态，需要下面的函数输出一次。当输出10次之后，调用quit<-0，所以上方函数便进行推出操作
func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

// 输出结果
// PS E:\mygo\src> go run select.go
// 0
// 1
// 1
// 2
// 3
// 5
// 8
// 13
// 21
// 34
// quit
```

在`select`里面还有default语法，`select`其实就是类似switch的功能，default就是当监听的channel都没有准备好的时候，默认执行的（select不再阻塞等待channel）。 

- select处理阻塞

避免整个程序进入阻塞的情况，利用select来设置超时 

```go
func main() {
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
				case v := <- c:
					println(v)
				case <- time.After(5 * time.Second):
					println("timeout")
					o <- true
					break
			}
		}
	}()
	<- o
}
```

# runtime

runtime包中有几个处理goroutine的函数：

- Goexit

  退出当前执行的goroutine，但是defer函数还会继续调用

- Gosched

  让出当前goroutine的执行权限，调度器安排其他等待的任务运行，并在下次某个时候从该位置恢复执行。

- NumCPU

  返回 CPU 核数量

- NumGoroutine

  返回正在执行和排队的任务总数

- GOMAXPROCS

  用来设置可以并行计算的CPU核数的最大值，并返回之前的值。

# stringers

fmt包中定义Strings接口：

```go
type Stringer interface {
    String() string
}
```

`Stringer` 是一个可以用字符串描述自己的类型。 

```go
package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years old)", p.Name, p.Age)
}

func main() {
	a := Person{"Arthur Dent", 42}
	z := Person{"zaphod Beeblebrox", 34}
	fmt.Print(a,"\n", z)
}

// 输出结果
// PS E:\mygo\src> go run Stringers.go
// Arthur Dent (42 years old)
// zaphod Beeblebrox (34 years old)
```

# strings

```go
package main

import (
    "fmt"
    "strings"
    //"unicode/utf8"
)

func main() {
    fmt.Println("查找子串是否在指定的字符串中")
    fmt.Println(" Contains 函数的用法")
    fmt.Println(strings.Contains("seafood", "foo")) //true
    fmt.Println(strings.Contains("seafood", "bar")) //false
    fmt.Println(strings.Contains("seafood", ""))    //true
    fmt.Println(strings.Contains("", ""))           //true 这里要特别注意
    fmt.Println(strings.Contains("我是中国人", "我"))     //true

    fmt.Println("")
    fmt.Println(" ContainsAny 函数的用法")
    fmt.Println(strings.ContainsAny("team", "i"))        // false
    fmt.Println(strings.ContainsAny("failure", "u & i")) // true
    fmt.Println(strings.ContainsAny("foo", ""))          // false
    fmt.Println(strings.ContainsAny("", ""))             // false

    fmt.Println("")
    fmt.Println(" ContainsRune 函数的用法")
    fmt.Println(strings.ContainsRune("我是中国", '我')) // true 注意第二个参数，用的是字符

    fmt.Println("")
    fmt.Println(" Count 函数的用法")
    fmt.Println(strings.Count("cheese", "e")) // 3 
    fmt.Println(strings.Count("five", ""))    // before & after each rune result: 5 , 源码中有实现

    fmt.Println("")
    fmt.Println(" EqualFold 函数的用法")
    fmt.Println(strings.EqualFold("Go", "go")) //大小写忽略 

    fmt.Println("")
    fmt.Println(" Fields 函数的用法")
    fmt.Println("Fields are: %q", strings.Fields("  foo bar  baz   ")) //["foo" "bar" "baz"] 返回一个列表

    //相当于用函数做为参数，支持匿名函数
    for _, record := range []string{" aaa*1892*122", "aaa\taa\t", "124|939|22"} {
        fmt.Println(strings.FieldsFunc(record, func(ch rune) bool {
            switch {
            case ch > '5':
                return true
            }
            return false
        }))
    }

    fmt.Println("")
    fmt.Println(" HasPrefix 函数的用法")
    fmt.Println(strings.HasPrefix("NLT_abc", "NLT")) //前缀是以NLT开头的

    fmt.Println("")
    fmt.Println(" HasSuffix 函数的用法")
    fmt.Println(strings.HasSuffix("NLT_abc", "abc")) //后缀是以NLT开头的

    fmt.Println("")
    fmt.Println(" Index 函数的用法")
    fmt.Println(strings.Index("NLT_abc", "abc")) // 返回第一个匹配字符的位置，这里是4
    fmt.Println(strings.Index("NLT_abc", "aaa")) // 在存在返回 -1
    fmt.Println(strings.Index("我是中国人", "中"))     // 在存在返回 6

    fmt.Println("")
    fmt.Println(" IndexAny 函数的用法")
    fmt.Println(strings.IndexAny("我是中国人", "中")) // 在存在返回 6
    fmt.Println(strings.IndexAny("我是中国人", "和")) // 在存在返回 -1

    fmt.Println("")
    fmt.Println(" Index 函数的用法")
    fmt.Println(strings.IndexRune("NLT_abc", 'b')) // 返回第一个匹配字符的位置，这里是4
    fmt.Println(strings.IndexRune("NLT_abc", 's')) // 在存在返回 -1
    fmt.Println(strings.IndexRune("我是中国人", '中'))   // 在存在返回 6

    fmt.Println("")
    fmt.Println(" Join 函数的用法")
    s := []string{"foo", "bar", "baz"}
    fmt.Println(strings.Join(s, ", ")) // 返回字符串：foo, bar, baz 

    fmt.Println("")
    fmt.Println(" LastIndex 函数的用法")
    fmt.Println(strings.LastIndex("go gopher", "go")) // 3

    fmt.Println("")
    fmt.Println(" LastIndexAny 函数的用法")
    fmt.Println(strings.LastIndexAny("go gopher", "go")) // 4
    fmt.Println(strings.LastIndexAny("我是中国人", "中"))      // 6

    fmt.Println("")
    fmt.Println(" Map 函数的用法")
    rot13 := func(r rune) rune {
        switch {
        case r >= 'A' && r <= 'Z':
            return 'A' + (r-'A'+13)%26
        case r >= 'a' && r <= 'z':
            return 'a' + (r-'a'+13)%26
        }
        return r
    }
    fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))

    fmt.Println("")
    fmt.Println(" Repeat 函数的用法")
    fmt.Println("ba" + strings.Repeat("na", 2)) //banana 

    fmt.Println("")
    fmt.Println(" Replace 函数的用法")
    fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
    fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))

    fmt.Println("")
    fmt.Println(" Split 函数的用法")
    fmt.Printf("%q\n", strings.Split("a,b,c", ","))
    fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
    fmt.Printf("%q\n", strings.Split(" xyz ", ""))
    fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))

    fmt.Println("")
    fmt.Println(" SplitAfter 函数的用法")
    fmt.Printf("%q\n", strings.SplitAfter("/home/m_ta/src", "/")) //["/" "home/" "m_ta/" "src"]

    fmt.Println("")
    fmt.Println(" SplitAfterN 函数的用法")
    fmt.Printf("%q\n", strings.SplitAfterN("/home/m_ta/src", "/", 2))  //["/" "home/m_ta/src"]
    fmt.Printf("%q\n", strings.SplitAfterN("#home#m_ta#src", "#", -1)) //["#" "home#" "m_ta#" "src"]

    fmt.Println("")
    fmt.Println(" SplitN 函数的用法")
    fmt.Printf("%q\n", strings.SplitN("/home/m_ta/src", "/", 1))

    fmt.Printf("%q\n", strings.SplitN("/home/m_ta/src", "/", 2))  //["" "home/m_ta/src"]
    fmt.Printf("%q\n", strings.SplitN("/home/m_ta/src", "/", -1)) //["" "home" "m_ta" "src"]
    fmt.Printf("%q\n", strings.SplitN("home,m_ta,src", ",", 2))   //["home" "m_ta,src"]

    fmt.Printf("%q\n", strings.SplitN("#home#m_ta#src", "#", -1)) //["/" "home/" "m_ta/" "src"]

    fmt.Println("")
    fmt.Println(" Title 函数的用法") //这个函数，还真不知道有什么用
    fmt.Println(strings.Title("her royal highness"))

    fmt.Println("")
    fmt.Println(" ToLower 函数的用法")
    fmt.Println(strings.ToLower("Gopher")) //gopher 

    fmt.Println("")
    fmt.Println(" ToLowerSpecial 函数的用法")

    fmt.Println("")
    fmt.Println(" ToTitle 函数的用法")
    fmt.Println(strings.ToTitle("loud noises"))
    fmt.Println(strings.ToTitle("loud 中国"))

    fmt.Println("")
    fmt.Println(" Replace 函数的用法")
    fmt.Println(strings.Replace("ABAACEDF", "A", "a", 2)) // aBaACEDF
    //第四个参数小于0，表示所有的都替换， 可以看下golang的文档
    fmt.Println(strings.Replace("ABAACEDF", "A", "a", -1)) // aBaaCEDF

    fmt.Println("")
    fmt.Println(" ToUpper 函数的用法")
    fmt.Println(strings.ToUpper("Gopher")) //GOPHER

    fmt.Println("")
    fmt.Println(" Trim  函数的用法")
    fmt.Printf("[%q]", strings.Trim(" !!! Achtung !!! ", "! ")) // ["Achtung"]

    fmt.Println("")
    fmt.Println(" TrimLeft 函数的用法")
    fmt.Printf("[%q]", strings.TrimLeft(" !!! Achtung !!! ", "! ")) // ["Achtung !!! "]

    fmt.Println("")
    fmt.Println(" TrimSpace 函数的用法")
    fmt.Println(strings.TrimSpace(" \t\n a lone gopher \n\t\r\n")) // a lone gopher
}
```

# goroutine

- *goroutine* 是由 Go 运行时环境管理的轻量级线程。说到底其实就是协程，但是它比线程更小，十几个goroutine可能体现在底层就是五六个线程，Go语言内部帮你实现了这些goroutine之间的内存共享。执行goroutine只需极少的栈内存(大概是4~5KB)，当然会根据相应的数据伸缩。也正因为如此，可同时运行成千上万个并发任务。goroutine比thread更易用、更高效、更轻便
- goroutine是通过Go的runtime管理的一个线程管理器。goroutine通过`go`关键字实现了，其实就是一个普通的函数。

```go
package main

import (
	"fmt"
	// "time"
	// "runtime"
	"sync"
)
// sync.WaitGroup可以实现c语言中join类似的功能
// 调用wg.Wait()使得当所有goroutine都执行完毕时再继续向下执行。
var wg *sync.WaitGroup = &sync.WaitGroup{}

func print(s string) {
	for i := 0; i < 5; i++ {
		// time.Sleep(100 * time.Millisecond)
		// runtime.Gosched() - runtime.Gosched()表示让CPU把时间片让给别人,下次某个时候继续恢复执行该goroutine。
		fmt.Println(s)
	}
	wg.Done()
}

func main() {
	wg.Add(2)
	go print("world")
	print("hello")
	wg.Wait()
}
// 输出结果
// PS E:\mygo\src> go run goroutine.go
// world
// world
// world
// world
// world
// hello
// hello
// hello
// hello
// hello
```

# channel

- channel 是有类型的管道，可以用 channel 操作符 <-对其发送或者接收值。
- "<-"箭头表示数据流的方向
- 初始化需要用make：`ch:=make(chan int)`
- 默认情况下，在另一端准备好之前，发送和接收都会阻塞。这使得 goroutine 可以在没有明确的锁或竞态变量的情况下进行同步
- 所谓阻塞，也就是如果读取（value := <-ch）它将会被阻塞，直到有数据接收。其次，任何发送（ch<-5）将会被阻塞，直到数据被读出。无缓冲channel是在多个goroutine之间同步很棒的工具。 

```go
package main

import "fmt"

func sum(a []int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	c <- sum
}

func main() {
	a := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(a[:len(a)/2], c)
	go sum(a[len(a)/2:], c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}
// 输出结果
// PS E:\mygo\src> go run channel.go
// 17 -5 12
```

# bufferedChannel

- channel 可以是 *带缓冲的*。为 make提供第二个参数作为缓冲长度来初始化一个缓冲 channel：

- `ch := make(chan int, 100)`

- 向带缓冲的 channel 发送数据的时候，只有在缓冲区满的时候才会阻塞。 而当缓冲区为空的时候接收操作会阻塞。

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
```

输出结果：

```go
PS E:\mygo\src> go run bufferChannel.go
1
2
```

当channel的长度改为1时：

```go
PS E:\mygo\src> go run bufferChannel.go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
        E:/mygo/src/bufferChannel.go:8 +0x7a
exit status 2
```

- range 和 close

可以通过range，像操作slice或者map一样操作缓存类型的channel ：

```go
package main

import (
	"fmt"
)

func fibonacci(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x + y
	}
	close(c)
}

func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
```

`for i := range c`能够不断的读取channel里面的数据，直到该channel被显式的关闭。上面代码我们看到可以显式的关闭channel，生产者通过内置函数`close`关闭channel。关闭channel之后就无法再发送任何数据了，在消费方可以通过语法`v, ok := <-ch`测试channel是否被关闭。如果ok返回false，那么说明channel已经没有任何数据并且已经被关闭。 

> 记住应该在生产者的地方关闭channel，而不是消费的地方去关闭它，这样容易引起panic

# sync

- channel常用于各个goroutine间通信，顺便实现了同步互斥的特性

- 如果不需要进行通信，则可以采用sync标准库的Mutex的Lock和Unlock操作来完成

- 可以用 defer语句来保证互斥锁一定会被解锁

  ```go
  package main
  
  import (
      "fmt"
      "sync"
      "time"
  )
  
  // SafeCounter 的并发使用是安全的。
  type SafeCounter struct {
      v   map[string]int
      mux sync.Mutex
  }
  
  // Inc 增加给定 key 的计数器的值。
  func (c *SafeCounter) Inc(key string) {
      c.mux.Lock()
      // Lock 之后同一时刻只有一个 goroutine 能访问 c.v
      c.v[key]++
      c.mux.Unlock()
  }
  
  // Value 返回给定 key 的计数器的当前值。
  func (c *SafeCounter) Value(key string) int {
      c.mux.Lock()
      // Lock 之后同一时刻只有一个 goroutine 能访问 c.v
      defer c.mux.Unlock()
      return c.v[key]
  }
  
  func main() {
      c := SafeCounter{v: make(map[string]int)}
      for i := 0; i < 1000; i++ {
          go c.Inc("somekey")
      }
  
      time.Sleep(time.Second)
      fmt.Println(c.Value("somekey"))
  }
  // https://www.jianshu.com/p/bed39de53087
  // 输出结果：
  // PS E:\mygo\src> go run sync.go
  // 1000
  ```


# context

在go服务器中，对于每个请求的request都是在单独的goroutine中进行的，处理一个request也可能设计多个goroutine之间的交互， 使用context可以使开发者方便的在这些goroutine里传递request相关的数据、取消goroutine的signal或截止日期。 

- context结构：

  ```go
  // A Context carries a deadline, cancelation signal, and request-scoped values
  // across API boundaries. Its methods are safe for simultaneous use by multiple
  // goroutines.
  type Context interface {
      // Done returns a channel that is closed when this Context is canceled
      // or times out.
      Done() <-chan struct{}
  
      // Err indicates why this context was canceled, after the Done channel
      // is closed.
      Err() error
  
      // Deadline returns the time when this Context will be canceled, if any.
      Deadline() (deadline time.Time, ok bool)
  
      // Value returns the value associated with key or nil if none.
      Value(key interface{}) interface{}
  }
  ```

  **Done** 方法在Context被取消或超时时返回一个close的channel,close的channel可以作为广播通知，告诉给context相关的函数要停止当前工作然后返回。

  当一个父operation启动一个goroutine用于子operation，这些子operation不能够取消父operation。下面描述的WithCancel函数提供一种方式可以取消新创建的Context.

  Context可以安全的被多个goroutine使用。开发者可以把一个Context传递给任意多个goroutine然后cancel这个context的时候就能够通知到所有的goroutine。

  **Err**方法返回context为什么被取消。

  **Deadline**返回context何时会超时。

  **Value**返回context相关的数据。

- 继承的context

  - BackGround

    ```go
    // Background returns an empty Context. It is never canceled, has no deadline,
    // and has no values. Background is typically used in main, init, and tests,
    // and as the top-level Context for incoming requests.
    func Background() Context
    ```

    BackGound是所有Context的root，不能够被cancel。 

  - WithCancel

    ```go
    // WithCancel returns a copy of parent whose Done channel is closed as soon as
    // parent.Done is closed or cancel is called.
    func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
    ```

    WithCancel返回一个继承的Context,这个Context在父Context的Done被关闭时关闭自己的Done通道，或者在自己被Cancel的时候关闭自己的Done。 WithCancel同时还返回一个取消函数cancel，这个cancel用于取消当前的Context。 

    ```go
    package main
    
    import (
        "context"
        "log"
        "os"
        "time"
    )
    
    var logg *log.Logger
    
    func someHandler() {
        ctx, cancel := context.WithCancel(context.Background())
        go doStuff(ctx)
    
    //10秒后取消doStuff
        time.Sleep(10 * time.Second)
        cancel()
    
    }
    
    //每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
    func doStuff(ctx context.Context) {
        for {
            time.Sleep(1 * time.Second)
            select {
            case <-ctx.Done():
                logg.Printf("done")
                return
            default:
                logg.Printf("work")
            }
        }
    }
    
    func main() {
        logg = log.New(os.Stdout, "", log.Ltime)
        someHandler()
        logg.Printf("down")
    }
    
    // https://studygolang.com/articles/10155?fr=sidebar
    // 输出结果
    // PS E:\mygo\src> go run context1.go
    // 16:07:37 work
    // 16:07:38 work
    // 16:07:39 work
    // 16:07:40 work
    // 16:07:41 work
    // 16:07:42 work
    // 16:07:43 work
    // 16:07:44 work
    // 16:07:45 work
    // 16:07:46 down
    ```

  - WithDeadline & WithTimeout

    WithTimeout 等价于 WithDeadline(parent, time.Now().Add(timeout)). 

    ```go
    ctx, cancel := context.WithCancel(context.Background())
    //修改为
    ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
    //或者
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    // 输出结果
    // PS E:\mygo\src> go run context1.go
    // 16:09:12 work
    // 16:09:13 work
    // 16:09:14 work
    // 16:09:15 work
    // 16:09:16 done
    // 16:09:21 down
    ```

  - WithValue

    ```go
    func WithValue(parent context.Context, key interface{}, val interface{}) context.Context
    ```

    在parent的上下文中增加key，value条目

- context使用如下：

  ```go
  package main
  
  import (
  	"context"
  	"fmt"
  	"time"
  )
  
  func inc(sum int) int {
  	sum++
  	time.Sleep(1 * time.Second)
  	return sum
  }
  
  func Add(ctx context.Context, a, b int) int {
  	res := 0
  	for i := 0; i < a; i++ {
  		res = inc(res)
  		select {
  		case <-ctx.Done():
  			return res
  		default:
  		}
  	}
  
  	for i := 0; i < b; i++ {
  		res = inc(res)
  		select {
  		case <-ctx.Done():
  			return res
  		default:
  		}
  	}
  
  	return res
  }
  
  func main() {
  	a := 1
  	b := 2
  	{
  		timeout := 2 * time.Second
  		ctx, _ := context.WithTimeout(context.Background(), timeout)
  		res := Add(ctx, a, b)
  		fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
  	}
  	{
  		ctx, cancel := context.WithCancel(context.Background())
  		go func() {
  			time.Sleep(3 * time.Second)
  			cancel()
  		}()
  		res := Add(ctx, a, b)
  		fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
  	}
  }
  
  // 输出结果
  // PS E:\mygo\src> go run context.go
  // Compute: 1+2, result: 2
  // Compute: 1+2, result: 3
  ```


# flag

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	levelFlag = flag.Int("level", 0, "级别")
	bnFlag    int
)

func init() {
	flag.IntVar(&bnFlag, "bn", 3, "份数")
}

func main() {
	flag.Parse()
	count := len(os.Args)
	fmt.Println("参数总个数：", count)
	fmt.Println("参数详情：")
	for i := 0; i < count; i++ {
		fmt.Println(i, ":", os.Args[i])
	}
	fmt.Println("\n参数值：")
	fmt.Println("级别：", *levelFlag)
	fmt.Println("份数：", bnFlag)
}

// https://studygolang.com/articles/3365?t=1493776691081
// 输出结果
// PS E:\mygo\src> go run flag.go
// 参数总个数： 1
// 参数详情：
// 0 : C:\Users\liuyh73\AppData\Local\Temp\go-build381154242\b001\exe\flag.exe

// 参数值：
// 级别： 0
// 份数： 3
```

# json

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// 结构体内变量只有大写字母开头的成员才会被JSON处理到
type ColorGroup struct {
	ID       int
	Name     string
	Colors   []string
	Flag     bool
	MyFloat  float64
	MyStruct Strings
	MyMap    map[int]string
}

type Strings struct {
	Strings []string
}

var jsonBlob = []byte(`[
	{"Name": "Platypus", "Order": "Monostremata"},
	{"Name": "Quoll", "Order": "Dasyuromorphia"}
]`)

type Animal struct {
	Name string
	// Order string
}

// 加上tag位的序列化与反序列化
type Person struct {
	Name        string `json:"username"` //起一个别名
	Age         int
	Gender      bool `json:",omitempty"` //别名 字段为空，omitempty表示为空时，不进行序列化
	Profile     string
	OmitContent string `json:"-"`       // - 表示忽略该属性
	Count       int    `json:",string"` // 表示将该int序列化为string
}

func main() {

	group := ColorGroup{
		ID:       1,
		Name:     "Reds",
		Colors:   []string{"Crimson", "Red", "Ruby", "Maroon"},
		Flag:     true,
		MyFloat:  0.4,
		MyStruct: Strings{[]string{"test1", "test2"}},
		MyMap:    map[int]string{2: "test3", 3: "test4"},
	}

	b, err := json.Marshal(group)

	if err != nil {
		fmt.Println("error: ", err)
	}

	os.Stdout.Write(b)
	var animals []Animal
	json.Unmarshal(jsonBlob, &animals)

	fmt.Printf("%+v", animals)
	
	fmt.Println("\n***************\n序列化")
	var p *Person = &Person{
		Name:        "brainwu",
		Age:         21,
		Gender:      true,
		Profile:     "I am Wujunbin",
		OmitContent: "OmitConent",
	}

	bs, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(string(bs))
	}
	fmt.Println("反序列化")
	var p2 Person
	json.Unmarshal(bs, p2)
	fmt.Printf("%+v\n", p2)
}
//输出结果
// PS E:\mygo\src> go run json.go
// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"],"Flag":true,"MyFloat":0.4,"MyStruct":{"Strings":["test1","test2"]},"MyMap":{"2":"test3","3":"test4"}}[{Name:Platypus} {Name:Quoll}]
// ***************
// 序列化
// {"username":"brainwu","Age":21,"Gender":true,"Profile":"I am Wujunbin","Count":"0"}
// 反序列化
// {Name: Age:0 Gender:false Profile: OmitContent: Count:0}
```

# reflect

Golang的反射机制 

一般数据类型：

```go
package main 

import (
	"fmt"
	"reflect"
)

func main() {
	var x int = 1
	fmt.Println("x type: ", reflect.TypeOf(x))
	fmt.Println("x value: ", reflect.ValueOf(x))
	fmt.Println("x Kind: ", reflect.ValueOf(x).Kind())
	
	// 修改反射对象
	var y float64 = 3.4
	v1 := reflect.ValueOf(y)
	fmt.Println("settability of v: ", v1.CanSet())

	// 那么如何才能设置该值呢？
	// 这里需要考虑一个常见的问题，参数传递，传值还是传引用或地址？
	// 在上面的例子中，我们使用的是reflect.ValueOf(y)，这是一个值传递，传递的是y的值的一个副本，不是y本身，因此更新副本中的值是不允许的。如果使用reflect.ValueOf(&y)来替换刚才的值传递，就可以实现值的修改。
	p := reflect.ValueOf(&y) // 获取y的地址
	fmt.Println("settability of p: ", p.CanSet())
	v2 := p.Elem()
	fmt.Println("settability of v2: ", v2.CanSet())
	v2.SetFloat(7.1)
	// 通过Interface函数可以实现反射对象到接口值的转换
	fmt.Println(v2.Interface())
	fmt.Println(y)
}
//输出结果
// PS E:\mygo\src> go run reflect.go
// x type:  int
// x value:  1
// x Kind:  int
// settability of v:  false
// settability of p:  false
// settability of v2:  true
// 7.1
// 7.1
```

struct结构体反射：

```go
package main 

import (
	"fmt"
	"reflect"
)

type Person struct {
	name string
}

type T struct {
    A int
    B string
}

func main () {
	// 1. 了解reflect操作struct的ValueOf和TypeOf等
	person := new(Person)
	person.name = "lyh73"
	fmt.Println("type:", reflect.TypeOf(person))
	fmt.Println("value:", reflect.ValueOf(person))
	// CanSet当Value是可寻址的时候，返回true，否则返回false。
	//当前面的CanSet是一个指针的时候（p）它是不可寻址的，但是当是p.Elem()(实际上就是*p)，它就是可以寻址的
	// fmt.Println("name:", reflect.ValueOf(person).FieldByName("name"))

	person2 := Person{
		"lyh74",
	}
	fmt.Println("type:", reflect.TypeOf(person2))
	fmt.Println("value:", reflect.ValueOf(person2))
	fmt.Println("name:", reflect.ValueOf(person2).FieldByName("name"))

	// 2. 遍历结构体字段内容
	{
		t := T{12, "skidoo"}
		s := reflect.ValueOf(&t).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Printf("%d %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		}	
	}

	// 3. 获取结构体的标签内容
	{
		type S struct {
			F string `species:"gopher" color:"blue"`
		}
	
		s := S{}
		st := reflect.TypeOf(s)
		field := st.Field(0)
		fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
	}
}

// https://studygolang.com/articles/8742
// 输出结果
// PS E:\mygo\src> go run reflect_struct.go
// type: *main.Person
// value: &{lyh73}
// type: main.Person
// value: {lyh74}
// name: lyh74
// 0 A int = 12
// 1 B string = skidoo
// blue gopher
```








[TOC]

# session和cookie

## cookie

**cookie**，简而言之就是在本地计算机保存一些用户操作的历史信息（当然包括登录信息），并在用户再次访问该站点时浏览器通过HTTP协议将本地cookie内容发送给服务器，从而完成验证，或继续上一步操作。 

![img](https://github.com/astaxie/build-web-application-with-golang/raw/master/zh/images/6.1.cookie2.png?raw=true) 

cookie是有时间限制的，根据生命期不同分成两种：会话cookie和持久cookie；

如果不设置过期时间，则表示这个cookie的生命周期为从创建到浏览器关闭为止，只要关闭浏览器窗口，cookie就消失了。这种生命期为浏览会话期的cookie被称为会话cookie。会话cookie一般不保存在硬盘上而是保存在内存里。

如果设置了过期时间(setMaxAge(60*60*24))，浏览器就会把cookie保存到硬盘上，关闭后再次打开浏览器，这些cookie依然有效直到超过设定的过期时间。存储在硬盘上的cookie可以在不同的浏览器进程间共享，比如两个IE窗口。而对于保存在内存的cookie，不同的浏览器有不同的处理方式。 

### GO设置cookie

Go语言通过`net/http`包中的`SetCookie`来设置：

```go
http.SetCookie(w ResponseWriter, cookie *Cookie)
```

w表示需要写入的response，cookie是一个struct，让我们来看一下cookie对象是怎么样的

```go
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

我们来看一个例子，如何设置cookie

```go
expiration := time.Now()
expiration = expiration.AddDate(1, 0, 0)
cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
http.SetCookie(w, &cookie)
```

### Go读取cookie

```go
cookie, _ := r.Cookie("username")
fmt.Fprint(w, cookie)
```

还有另外一种读取方式

```go
for _, cookie := range r.Cookies() {
	fmt.Fprint(w, cookie.Name)
}
```

可以看到通过request获取cookie非常方便。

## session

**session**，简而言之就是在服务器上保存用户操作的历史信息。服务器使用session id来标识session，session id由服务器负责产生，保证随机性与唯一性，相当于一个随机密钥，避免在握手或传输中暴露用户真实密码。但该方式下，仍然需要将发送请求的客户端与session进行对应，所以可以借助cookie机制来获取客户端的标识（即session id），也可以通过GET方式将id提交给服务器。 

![img](https://github.com/astaxie/build-web-application-with-golang/raw/master/zh/images/6.1.session.png?raw=true) 

session，中文经常翻译为会话，其本来的含义是指有始有终的一系列动作/消息，比如打电话是从拿起电话拨号到挂断电话这中间的一系列过程可以称之为一个session。然而当session一词与网络协议相关联时，它又往往隐含了“面向连接”和/或“保持状态”这样两个含义。

session在Web开发环境下的语义又有了新的扩展，它的含义是指一类用来在客户端与服务器端之间保持状态的解决方案。有时候Session也用来指这种解决方案的存储结构。

session机制是一种服务器端的机制，服务器使用一种类似于散列表的结构(也可能就是使用散列表)来保存信息。但程序需要为某个客户端的请求创建一个session的时候，服务器首先检查这个客户端的请求里是否包含了一个session标识－称为session id，如果已经包含一个session id则说明以前已经为此客户创建过session，服务器就按照session id把这个session检索出来使用(如果检索不到，可能会新建一个，这种情况可能出现在服务端已经删除了该用户对应的session对象，但用户人为地在请求的URL后面附加上一个JSESSION的参数)。如果客户请求不包含session id，则为此客户创建一个session并且同时生成一个与此session相关联的session id，这个session id将在本次响应中返回给客户端保存。 

### session创建过程：

session的基本原理是由服务器为每个会话维护一份信息数据，客户端和服务端依靠一个全局唯一的标识来访问这份数据，以达到交互的目的。当用户访问Web应用时，服务端程序会随需要创建session，这个过程可以概括为三个步骤：

- 生成全局唯一标识符（sessionid）；
- 开辟数据存储空间。一般会在内存中创建相应的数据结构，但这种情况下，系统一旦掉电，所有的会话数据就会丢失，如果是电子商务类网站，这将造成严重的后果。所以为了解决这类问题，你可以将会话数据写到文件里或存储在数据库中，当然这样会增加I/O开销，但是它可以实现某种程度的session持久化，也更有利于session的共享；
- 将session的全局唯一标示符发送给客户端。

以上三个步骤中，最关键的是如何发送这个session的唯一标识这一步上。考虑到HTTP协议的定义，数据无非可以放到请求行、头域或Body里，所以一般来说会有两种常用的方式：cookie和URL重写。

1. Cookie 服务端通过设置Set-cookie头就可以将session的标识符传送到客户端，而客户端此后的每一次请求都会带上这个标识符，另外一般包含session信息的cookie会将失效时间设置为0(会话cookie)，即浏览器进程有效时间。至于浏览器怎么处理这个0，每个浏览器都有自己的方案，但差别都不会太大(一般体现在新建浏览器窗口的时候)；
2. URL重写 所谓URL重写，就是在返回给用户的页面里的所有的URL后面追加session标识符，这样用户在收到响应之后，无论点击响应页面里的哪个链接或提交表单，都会自动带上session标识符，从而就实现了会话的保持。虽然这种做法比较麻烦，但是，如果客户端禁用了cookie的话，此种方案将会是首选。

主要模块设计逻辑如下图：主要分为三个模块（全局管理，底层存储结构，存储模块）

![7](E:\starlight\learning\images\7.png)

### session管理设计

我们知道session管理涉及到如下几个因素

- 全局session管理器
- 保证sessionid 的全局唯一性
- 为每个客户关联一个session
- session 的存储(可以存储到内存、文件、数据库等)
- session 过期处理

session全局管理器：

```go
package session

import (
	"fmt"
	"crypto/rand"
	"io"
	"encoding/base64"
	"net/url"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	cookieName string	// private cookiename
	lock sync.Mutex		// protects session
	provider Provider
	maxLifeTime int64
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

var providers = make(map[string]Provider)
// 该方法获得一个全局管理员，根据providerName和cookieName
func NewSessionManager(providerName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q (forgotten import?", providerName)
	}
	return &Manager{cookieName: cookieName, provider: provider, maxLifeTime: maxLifeTime}, nil
}
// 该方法注册一个provider，name即为ProviderName
func Register(name string, provider *Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	providers[name] = *provider
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err !=nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
// 初始化session
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid:=manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie:=http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return session
}

func (manager *Manager) SessionDestory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestory(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/",
			HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func(){manager.GC()})
}
```

### session底层存储结构

我们知道session是保存在服务器端的数据，它可以以任何的方式存储，比如存储在内存、数据库或者文件中。因此我们抽象出一个Provider接口，用以表征session管理器底层存储结构。 

- SessionInit函数实现Session的初始化，操作成功则返回此新的Session变量
- SessionRead函数返回sid所代表的Session变量，如果不存在，那么将以sid为参数调用SessionInit函数创建并返回一个新的Session变量
- SessionDestroy函数用来销毁sid对应的Session变量
- SessionGC根据maxLifeTime来删除过期的数据
- SessionUpdate更新session的最近访问时间

provider底层存储结构：

```go
// 该模块主要提供对session的管理，与manager的相应函数对应
package session

import (
	"container/list"
	"time"
	"sync"
)
// sessions数据结构用于存储session，根据sessionId来得到相应的session
// list存储session列表，用于GC()销毁处理
type Provider struct {
	lock sync.Mutex
	sessions map[string]*list.Element
	list *list.List
}
// 生成一个session
func (provider *Provider) SessionInit(sid string) (Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := provider.list.PushBack(newsess)
	provider.sessions[sid] = element
	return newsess, nil
}
// 获取session，若不存在，则调用SessionInit创建session
func (provider *Provider) SessionRead(sid string) (Session, error)  {
	if element, ok := provider.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := provider.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}
// 销毁session
func (provider *Provider) SessionDestory(sid string) error {
	if element, ok := provider.sessions[sid]; ok {
		delete(provider.sessions, sid)
		provider.list.Remove(element)
		return nil
	}
	return nil
}
// 根据maxLifeTime销毁session
func (provider *Provider) SessionGC(maxLifeTime int64) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	for {
		element := provider.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			provider.list.Remove(element)
			delete(provider.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}
// 更新session的最近访问时间
func (provider *Provider) SessionUpdate(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	if element, ok := provider.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		provider.list.MoveToFront(element)
		return nil
	}
	return nil
}
```

### 存储模块：

```go
// 该模块主要是对于session的相关信息
package session

import (
	"container/list"
	"time"
)
// value中存储的是session内对应的key-value键值对
type SessionStore struct {
	sid string
	timeAccessed time.Time
	value map[interface{}]interface{}
}
// pder为一个Provider对象指针，其内部list初始化为空，默认没有session
var pder = &Provider{list: list.New()}
// 设置session中的值
func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}
// 根据key或者session中的值
func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}
// 根据key删除相对应的value
func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}
// 或者sessionId
func (st *SessionStore) SessionID() string {
	return st.sid
}
// init初始化函数，Register注册一个名为“memory”的Provider对象pder，其内部sessions为空
func init(){
	pder.sessions = make(map[string]*list.Element, 0)
	Register("memory", pder)
}
```

### main函数中使用方法：

```go
// 在初始化函数中，调用session应用包中的NewSessionManager函数，获取一个providerName=“memory”，cookieName=“goSeesionId”，maxLifeTime=3600的管理员
// 并用goroutine协程来执行GC()
import "session"
var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewSessionManager("memory", "goSessionid", 3600)
	go globalSessions.GC()
}
// 由于导入session的过程中已经调用了memory中的init函数，所以名为“memory”的provider已经注册成功，所以上面的函数成功返回一个管理员对象。
```

session测试：

```go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"session"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, err := template.ParseFiles("login.gtpl")
		if err != nil {
			fmt.Println(err)
			return
		}
		t.Execute(w, nil)
 
	} else {
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
	}
}

// 模板引擎渲染必须将属性字段首字母大写！！！！
type Data struct {
	Username string
	Password string
}

func login2(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login2.gtpl")
		//w.Header().Set("Content-Type", "text/html")
		if sess.Get("username") != nil {
			t.Execute(w, Data{
				Username: sess.Get("username").([]string)[0],
				Password: "123",
			})
		} else {
			t.Execute(w, nil)
		}
	} else {
		// r.Form 数值类型为 url.values -> map[string][]string		
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/login2", 302)
	}
}
 
func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		globalSessions.SessionDestory(w, r)
		sess = globalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}
func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/login2", login2)
	http.HandleFunc("/count", count)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Listen and server", err)
	}
}
```

效果：

`login2`界面中输入用户名密码，点击登陆，会在当前页面显示出相应信息，之后重新打开浏览器，打开`login2`界面，相应信息还在；`count`界面是计数功能，刷新一次计数器加1，关闭页面再打开，会在之前数字的基础上加1（相当于刷新了一次页面）。
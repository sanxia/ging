package ging

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * ging web framework
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IHttpRouter interface {
		Route() *gin.Engine
	}

	bootstrapFunc func(IApp)

	commandLine struct {
		App      string
		Host     string
		Ports    []int
		IsSsl    bool
		IsDebug  bool
		IsTest   bool
		IsOnline bool
	}

	serverInfo struct {
		Index int
		Addr  string
		Now   time.Time
	}

	serverStatus struct {
		stopC chan bool
	}
)

const (
	DEBUG   string = "debug"
	RELEASE string = "release"
)

var (
	bootstrapFuncs []bootstrapFunc
	engingStatus   *serverStatus
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialization ging
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func init() {
	fmt.Printf("%v ging engine init\n", time.Now())

	bootstrapFuncs = make([]bootstrapFunc, 0)

	engingStatus = &serverStatus{
		stopC: make(chan bool),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * bootstrap
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Bootstrap(args ...bootstrapFunc) {
	fmt.Printf("%v ging engine bootstrap\n", time.Now())

	for _, _bootstrapFunc := range args {
		bootstrapFuncs = append(bootstrapFuncs, _bootstrapFunc)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * launch an app
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Start() {
	fmt.Printf("%v ging engine start\n", time.Now())

	//parse cmd
	cmdLine, err := parseCommandLine()
	if err != nil {
		panic(err)
	}

	//active current app
	currentApp := GetApp(cmdLine.App)
	if currentApp == nil {
		panic(errors.New(fmt.Sprintf("app name: %s is not found", cmdLine.App)))
	}

	activingApp(currentApp)

	// command line highest priority
	commandAppSetting(currentApp.GetSetting(), cmdLine)

	bootstrapApp(currentApp)

	commandAppSetting(currentApp.GetSetting(), cmdLine)

	// task
	for _, task := range currentApp.GetTasks() {
		go task.Run(currentApp)
	}

	serverInfoChan := make(chan serverInfo)

	// start http
	currentSetting := currentApp.GetSetting()
	for index, port := range currentSetting.Server.Ports {
		host := "127.0.0.1"
		if len(currentSetting.Server.Host) > 0 {
			host = currentSetting.Server.Host
		}

		if port == 0 {
			port = 19811
		}

		addr := fmt.Sprintf("%s:%d", host, port)

		go httpServe(index, addr, currentApp.GetRouter(), serverInfoChan)
	}

	index := 0
	for server := range serverInfoChan {
		index++

		info := fmt.Sprintf("%v ging engine server %02d on %s Success", server.Now, server.Index, server.Addr)
		log.Println(info)

		if index == len(currentSetting.Server.Ports) {
			close(serverInfoChan)
		}
	}

	pending := fmt.Sprintf("%v ging engine server Running ...", time.Now())
	log.Println(pending)

	for {
		select {
		case <-engingStatus.stopC:
			pending = fmt.Sprintf("%v ging engine Exit ...", time.Now())
			log.Println(pending)
			return
		case <-time.After(500 * time.Millisecond):
			time.Sleep(5 * time.Second)
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetToken(ctx *gin.Context) IToken {
	var currentToken IToken

	if ctx != nil {
		if tokenIdentity, isOk := ctx.Get(TOKEN_IDENTITY); tokenIdentity != nil && isOk {
			if userToken, isOk := tokenIdentity.(*Token); isOk {
				currentToken = userToken
			}
		}
	}

	return currentToken
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * is ajax request
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func IsAjax(ctx *gin.Context) bool {
	var isAjax bool

	xRequestHeader := ctx.Request.Header.Get("x-requested-with")
	if xRequestHeader != "" {
		if strings.ToLower(xRequestHeader) == "xmlhttprequest" {
			isAjax = true
		}
	} else {
		contentType := ctx.Request.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			isAjax = true
		}
	}

	log.Printf("ging engine IsAjax RequestHeader: %#v", ctx.Request.Header)

	return isAjax
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * parse app command line
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func parseCommandLine() (*commandLine, error) {
	fmt.Printf("%v ging engine parse command line\n", time.Now())

	appNameFlag := flag.String("app", "", "请输入要启动的app名称")
	appHostFlag := flag.String("host", "127.0.0.1", "请输入要绑定的主机ip")
	appPortFlag := flag.String("port", "19811", "请输入要绑定的主机端口号")
	sslFlag := flag.String("ssl", "false", "请输入是否SSL模式")
	debugFlag := flag.String("debug", "true", "请输入是否调试模式")
	testFlag := flag.String("test", "false", "请输入是否测试模式")
	onlineFlag := flag.String("online", "false", "请输入是否线上")

	flag.Parse()

	if len(*appNameFlag) == 0 {
		return nil, errors.New("sorry, the app name is empty")
	}

	cmdLine := new(commandLine)
	cmdLine.App = *appNameFlag
	cmdLine.Host = *appHostFlag

	appPort := *appPortFlag
	if len(appPort) > 0 {
		cmdLine.Ports = glib.StringToIntSlice(appPort)
	}

	cmdLine.IsSsl = glib.StringToBool(*sslFlag)
	cmdLine.IsDebug = glib.StringToBool(*debugFlag)
	cmdLine.IsTest = glib.StringToBool(*testFlag)
	cmdLine.IsOnline = glib.StringToBool(*onlineFlag)

	return cmdLine, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * reset app setting
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func commandAppSetting(appSetting *Setting, cmdLine *commandLine) {
	appSetting.Server.Host = cmdLine.Host
	appSetting.Server.Ports = cmdLine.Ports

	appSetting.Domain.IsSsl = cmdLine.IsSsl
	appSetting.Domain.IsDebug = cmdLine.IsDebug
	appSetting.Domain.IsTest = cmdLine.IsTest
	appSetting.Domain.IsOnline = cmdLine.IsOnline
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * activing current app
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func activingApp(currentApp IApp) {
	fmt.Printf("%v ging engine activing app\n", time.Now())

	for _, _currentApp := range apps {
		if _currentApp.GetName() == currentApp.GetName() {
			if currentApp, isOk := currentApp.(*app); isOk {
				currentApp.IsActiving = true
			}
			break
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * bootstrap current app
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func bootstrapApp(currentApp IApp) {
	fmt.Printf("%v ging engine bootstrap app\n", time.Now())

	for _, _bootstrapFunc := range bootstrapFuncs {
		_bootstrapFunc(currentApp)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * http service
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func httpServe(index int, addr string, httpRouter IHttpRouter, serverInfoChan chan<- serverInfo) {
	log.Printf("%v ging engine start serve %d\n", time.Now(), index)

	routeHandler := httpRouter.Route()

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        routeHandler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverInfoChan <- serverInfo{
		Index: index,
		Addr:  addr,
		Now:   time.Now(),
	}

	httpServer.ListenAndServe()
}

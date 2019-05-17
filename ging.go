package ging

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

import (
	"github.com/gin-gonic/gin"
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

	ServerStatus struct {
		Index  int
		Addr   string
		Now    time.Time
		Status chan ServerStatus
	}
)

var (
	serverStatus *ServerStatus
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * ging 初始化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func init() {
	fmt.Printf("%v ging init\n", time.Now())

	serverStatus = &ServerStatus{
		Status: make(chan ServerStatus),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 启动App
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Start(app IApp) {
	fmt.Printf("%v ging start\n", time.Now())

	host := app.GetSettings().Server.Host
	ports := app.GetSettings().Server.Ports
	for index, port := range ports {
		if len(host) == 0 {
			host = "127.0.0.1"
		}

		if port == 0 {
			port = 80
		}

		addr := fmt.Sprintf("%s:%d", host, port)

		go httpServe(index, addr, app.GetRouter())
	}

	index := 0
	for server := range serverStatus.Status {
		index++

		info := fmt.Sprintf("%v ging server %02d on %s Success", server.Now, server.Index, server.Addr)
		log.Println(info)

		if index == len(ports) {
			close(serverStatus.Status)
		}
	}

	for {
		pending := fmt.Sprintf("%v ging server Running ...", time.Now())
		log.Println(pending)

		time.Sleep(12 * time.Hour)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 提供Http服务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func httpServe(index int, addr string, httpRouter IHttpRouter) {
	log.Printf("%v ging start serve\n", time.Now())

	routeHandler := httpRouter.Route()

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        routeHandler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverStatus.Status <- ServerStatus{
		Index: index,
		Addr:  addr,
		Now:   time.Now(),
	}

	httpServer.ListenAndServe()
}

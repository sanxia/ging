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
 * ging
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

type (
	IHttpRouter interface {
		Route() *gin.Engine
	}

	ServerOption struct {
		Host  string
		Ports []int
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

func init() {
	fmt.Printf("%v ging init\n", time.Now())
	serverStatus = &ServerStatus{
		Status: make(chan ServerStatus),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 启动服务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Start(serverOption *ServerOption, router IHttpRouter) {
	fmt.Printf("%v ging Start\n", time.Now())
	//Routines
	for index, port := range serverOption.Ports {
		addr := fmt.Sprintf("%s:%d", serverOption.Host, port)

		//Http服务
		go startServe(index, addr, router)
	}

	index := 0
	for server := range serverStatus.Status {
		index++

		info := fmt.Sprintf("%v ging Server %02d on %s Success", server.Now, server.Index, server.Addr)
		log.Println(info)

		if index == len(serverOption.Ports) {
			close(serverStatus.Status)
		}
	}

	for true {
		pending := fmt.Sprintf("%v ging Serve Running ...", time.Now())
		log.Println(pending)
		time.Sleep(60 * time.Minute)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 提供Http服务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func startServe(index int, addr string, httpRouter IHttpRouter) {
	log.Printf("%v ging start serve\n", time.Now())
	routeHandler := httpRouter.Route()
	server := &http.Server{
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

	server.ListenAndServe()
}

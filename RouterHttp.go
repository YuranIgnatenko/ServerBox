package ServerBox

import (
	"fmt"
	"net/http"
	"time"
)

type RouterHttp struct {
	PORT         string   // ":8001"
	ListHandlers []string // ["/home" "/about" "/ajax_admin"]
	CntrHandlers int      // 0...n ;count additions handlers
}

// use port default: "8001"
func NewRouterHttp(port string) *RouterHttp {
	router := RouterHttp{}
	router.PORT = ":" + port
	router.ListHandlers = make([]string, 0)
	router.CntrHandlers = 0
	return &router
}

// get list handler addresses -> ["/home", "/about" ...]
func (router *RouterHttp) GetListHandlers() []string {
	return router.ListHandlers
}

// set directoru to mode public 
func (router *RouterHttp) SetDirectoryToPublic(directory string) {
	http.Handle("/"+directory+"/", http.StripPrefix("/"+directory+"/", http.FileServer(http.Dir("./"+directory+""))))
}

// ("/home", homeHandler)
func (router *RouterHttp) AddHandlerHttp(path string, function func(w http.ResponseWriter, r *http.Request)) {
	router.CntrHandlers++
	router.ListHandlers = append(router.ListHandlers, path)
	http.HandleFunc(path, function)
}

// running server; information in console output - true/false
func (router *RouterHttp) Listen(info bool) {
	if info {
		fmt.Printf("\n[Server running] [port%v] [time:%v]\n", router.PORT, time.Now().Format(time.RFC1123))
	}
	http.ListenAndServe(router.PORT, nil)
}

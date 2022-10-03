package ServerBox

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// use port default: "8001"
type RouterMux struct {
	PORT         string // ":8001"
	Mux          *http.ServeMux
	CntrHandlers int // count additions handlers
}

// use port default: "8001"
func NewRouterMux(port string) *RouterMux {
	router := RouterMux{}
	router.PORT = ":" + port
	router.CntrHandlers = 0
	router.Mux = http.NewServeMux()
	return &router
}

// get list handler addresses -> ["/home", "/about" ...]
func (router *RouterMux) GetListHandlers() []string {
	list := make([]string, 0)
	res := S("%#v", router.Mux)
	r := strings.Split(res, "m:map[string]")
	c := 0
	c2 := 0
	for _, v := range r {
		for _, val := range strings.Split(v, ":(http.HandlerFunc)") {
			c++
			if c < 3 {
				continue
			}
			for _, value := range strings.Split(val, "pattern") {
				c2++
				if c2%2 != 0 {
					continue
				}
				path := strings.Split(value, "}")[0]
				path = path[2 : len(path)-1]
				list = append(list, path)
			}
		}
	}
	return list
}

// ("/home", "index_home.html")
func (router *RouterMux) AddHandlerHtmlPage(path, filename string) {
	router.CntrHandlers++
	f := func(w http.ResponseWriter, r *http.Request) {
		var temp = template.Must(template.ParseFiles(filename))
		temp.Execute(w, nil)
	}
	router.Mux.HandleFunc(path, f)
}

// ("/home", homeHandler)
func (router *RouterMux) AddHandlerHttp(path string, function func(w http.ResponseWriter, r *http.Request)) {
	router.CntrHandlers++
	http.HandleFunc(path, function)

}

// ("/home", homeHandlerFunc)
func (router *RouterMux) AddHandler(path string, function func(w http.ResponseWriter, r *http.Request)) {
	router.CntrHandlers++
	router.Mux.HandleFunc(path, function)
}

// running server; information in console output - true/false
func (router *RouterMux) Listen(info bool, muxOn bool) {
	if info {
		fmt.Printf("Server running port%v time:%v\n", router.PORT, time.Now().Format(time.RFC1123))
	}
	http.ListenAndServe(router.PORT, router.Mux)
}

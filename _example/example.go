package main

import (
	. "WebGui"
	"net/http"
	"time"
)

/*
template panel links control
for addition in every page
*/
func template_links(w http.ResponseWriter) {
	F(w, "<a href=\"/js\">js Page</a>")
	F(w, "<a href=\"/jsfile\">jsfile Page</a>")
	F(w, "<a href=\"/img\">img Page</a>")
	F(w, "<a href=\"/time\">time Page</a>")

}

/*
page time
*/
func timeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	OnWhile(w, 1) // on loop while referesh for update time in page
	t := time.Now().Format(time.RFC3339)
	if er := AddLine(w, t, 1); er != nil {
		P("%v", er)
	}
	template_links(w)
}

/*
page runner javascript string
*/
func jsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	RunScriptJS(w, `
	function myAlert() {
    	alert("12345...");
      	return;
  	};
	myAlert();`)
	template_links(w)
}

/*
page runner javascript from file
*/
func jsfileHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	RunFileJS(w, "script.js")
	template_links(w)
}

func main() {
	/*
		init port
		init structure Router
		addition handlers
		view supported hsndlers
		run Listener
	*/
	PORT := "8001"
	wb := NewRouter(PORT)
	wb.AddHandler("/time", timeHandler)
	wb.AddHandler("/js", jsHandler)
	wb.AddHandler("/jsfile", jsfileHandler)
	wb.AddHandlerHtmlPage("/img", "index_img.html")
	P("%v\n", wb.GetListHandlers())
	wb.Listen(true)
}

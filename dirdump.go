package main

import (
	"io"
	"http"
	"os"
	"flag"
)

var (
	password = flag.String("pwd", "", "Password to upload/download files")
	// IP:pwdok
	loggedIn  = ""
	rootDir   = ""
	filesDir  = ""
	tmplDir   = ""
	staticDir = ""
	test      = true
)

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		println(err.String())
		os.Exit(1)
	}
	rootDir = dir
	filesDir = dir + "/files"
	tmplDir = dir + "/templates"
	staticDir = dir + "/static"

	http.Handle("/", http.HandlerFunc(dispatch))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		println(err.String())
	}
}

func dispatch(w http.ResponseWriter, req *http.Request) {
	pageName := ""

	if req.URL.Path == "/" {
		pageName = "index"
	} else {
		pageName = req.URL.Path[1:]
	}

	var err os.Error
	p := new(page)
	ok := true


	if ok {
		p, err = loadPage(req, pageName)
	} else {
		p, err = loadPage(req, "login")
	}
	if err == nil {
		err = p.serve(w)
	}

	if err != nil {
		error(w, err)
	}
}

func error(w http.ResponseWriter, err os.Error) {
	io.WriteString(w, "<b>Error:</b><br />"+err.String())
}

package main

import (
	"io"
	"io/ioutil"
	"http"
	"os"
	"template"
	"strconv"
	"mime"
	"path"
	//	"time"
)

type page struct {
	filename string
	filedir  string
	template string
	html     map[string]interface{}
}

func (p *page) serve(w http.ResponseWriter) os.Error {
	// render template
	if p.template != "" {
		t, err := template.ParseFile(tmplDir+"/"+p.template, nil)
		if err != nil {
			return err
		} else {
			err = t.Execute(w, p.html)
			if err != nil {
				println(err)
			}
		}
		return nil

		// serve file	
	} else if p.filename != "" {
		os.Chdir(p.filedir)

		file, err := os.Open(p.filename)
		if err != nil {
			return err
		}
		fs, err := file.Stat()
		if err != nil {
			return err
		}

		ext := path.Ext(p.filename)
		ctype := mime.TypeByExtension(ext)
		if ctype == "" {
			return os.NewError("Failed to determine file type for " + p.filename)
		}
		w.Header().Set("Content-Type", ctype)
		w.Header().Set("Content-Length", strconv.Itoa64(fs.Size))
		w.Header().Set("Content-Disposition", "inline; filename="+p.filename)
		io.Copy(w, file)
		return nil
	}
	return os.NewError("nothing to display")

}

func loadPage(req *http.Request, name string) (*page, os.Error) {
	pages := make(map[string]func(*http.Request) (*page, os.Error))
	// map page name to page function
	pages["upload"] = upload
	pages["index"] = index
	pages["login"] = login
	pages["download"] = download
	pages["static"] = static

	_, ok := pages[name]
	if ok {
		return pages[name](req)
	}
	return nil, os.NewError("page not found: " + name)
}

// pages ---------------------------------------------------------

func upload(req *http.Request) (*page, os.Error) {
	var p page
	p.template = "upload.html"
	p.html = make(map[string]interface{})
	p.html["Body"] = ""

	// show upload form if no post
	if req.Method != "POST" {
		return &p, nil
	}
	f, fh, err := req.FormFile("upfile")
	if err != nil {
		return nil, err
	} else {
		file, err := os.Create(filesDir + "/" + fh.Filename)

		if err != nil {
			return nil, err
		} else {

			io.Copy(file, f)
			file.Close()
			p.html["Body"] = fh.Filename + " uploaded!"
		}
		f.Close()
	}

	return &p, nil
}

func login(req *http.Request) (*page, os.Error) {
	var p page
	p.template = "login.html"
	p.html = make(map[string]interface{})
	if req.Method != "POST" {
		return &p, nil
	}
	pwd := req.FormValue("password")
	if pwd == *password {
		p.html["Body"] = "Logged in!"
		// TODO: save login
	} else {
		p.html["Body"] = "Login failed!"
	}

	return &p, nil
}

func index(req *http.Request) (*page, os.Error) {
	var p page
	p.template = "index.html"
	p.html = make(map[string]interface{})
	p.html["Title"] = "Files to download"

	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		return nil, err
	}

	// Make slice of filenames
	fl := []string{}
	for _, file := range files {
		if file.Name[0:1] != "." {
			fl = append(fl, file.Name)
		}
	}

	p.html["FileList"] = fl

	return &p, nil
}

func download(req *http.Request) (*page, os.Error) {
	var p page
	filename := ""
	query, _ := http.ParseQuery(req.URL.RawQuery)
	q, ok := query["file"]
	if ok {
		filename = q[0]
	} else {
		return nil, os.NewError("failed to parse query string")
	}
	if filename == "style.css" {
		p.filedir = tmplDir
	} else {
		p.filedir = filesDir
	}
	p.filename = filename
	return &p, nil
}

func static(req *http.Request) (*page, os.Error) {
	var p page
	filename := ""
	query, _ := http.ParseQuery(req.URL.RawQuery)
	q, ok := query["file"]
	if ok {
		filename = q[0]
	} else {
		return nil, os.NewError("Failed to parse the query string")
	}
	p.filedir = staticDir

	p.filename = filename
	return &p, nil
}

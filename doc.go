package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"

	"github.com/russross/blackfriday"
)

const indexFile = "README.md"

var dir = http.Dir(".")

type MarkdownServer struct {
	fs http.Handler
}

func (m *MarkdownServer) handleMarkdown(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(string(dir), r.URL.Path, indexFile)

	_, err := os.Stat(filename)
	if err == nil {
		r.URL.Path = r.URL.Path + indexFile
	}

	if !strings.HasSuffix(r.URL.Path, ".md") {
		m.fs.ServeHTTP(w, r)
		return
	}

	rec := httptest.NewRecorder()
	m.fs.ServeHTTP(rec, r)

	unsafe := blackfriday.MarkdownCommon(rec.Body.Bytes())

	tmpl, err := template.New("file").Parse(fileHTML)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data := struct {
		Output template.HTML
	}{
		Output: template.HTML(string(unsafe)),
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (m *MarkdownServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handleMarkdown(w, r)

}

func NewMarkdownServer(fs http.Handler) http.Handler {
	return &MarkdownServer{fs}
}

func main() {
	fs := http.FileServer(dir)

	http.Handle("/", NewMarkdownServer(fs))
	http.ListenAndServe(":8080", nil)
}

const fileHTML = `<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="/doc.css">
</head>
<body>
<div>
<form method="get" action="/search">
<input type="text" name="q" placeholder="Search documentation">
</form>
</div>
<div>
{{.Output}}
</div>
</body>
</html>`

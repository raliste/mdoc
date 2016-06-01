package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/russross/blackfriday"
)

const indexFile = "README.md"

var dir = flag.String("dir", ".", "")

type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Dir struct {
	Path       string        `json:"path"`
	Name       string        `json:"name"`
	IsDir      bool          `json:"is_dir"`
	IsMarkdown bool          `json:"is_markdown"`
	Size       int64         `json:"size,omitempty"`
	Content    template.HTML `json:"content,omitempty"`
	Children   []Dir         `json:"children,omitempty"`
}

func dirList(w http.ResponseWriter, name string, f http.File) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		return
	}
	sort.Sort(byName(dirs))

	d := Dir{
		Name:     name,
		IsDir:    true,
		Children: []Dir{},
	}

	for _, dd := range dirs {
		d.Children = append(d.Children, Dir{
			Name:       dd.Name(),
			IsDir:      dd.IsDir(),
			Size:       dd.Size(),
			IsMarkdown: strings.HasSuffix(dd.Name(), ".md"),
		})
	}

	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%s", out)

}

type MarkdownServer struct {
	dir http.FileSystem
}

func NewMarkdownServer(dir string) http.Handler {
	return &MarkdownServer{http.Dir(dir)}
}

// TODO factorization
// TODO error handling
func (m *MarkdownServer) handleMarkdown(w http.ResponseWriter, r *http.Request) {
	upath := strings.TrimPrefix(r.URL.Path, "/-/")

	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}

	f, err := m.dir.Open(path.Clean(upath))
	if err != nil {
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		return
	}

	name := d.Name()

	if d.IsDir() {
		index := strings.TrimSuffix(upath, "/") + "/README.md"
		ff, err := m.dir.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				name = dd.Name()
				f = ff
				d = dd
			}
		}

	}

	if d.IsDir() {
		dirList(w, name, f)
		return
	}

	if !strings.HasSuffix(name, ".md") {
		return
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	unsafe := blackfriday.MarkdownCommon(content)

	ff := Dir{
		Name:       name,
		IsDir:      false,
		IsMarkdown: true,
		Size:       d.Size(),
		Content:    template.HTML(unsafe),
	}

	out, err := json.MarshalIndent(ff, "", "  ")
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%s", out)
}

func (m *MarkdownServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handleMarkdown(w, r)
}

func main() {
	flag.Parse()

	http.Handle("/-/", NewMarkdownServer(*dir))
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("").Parse(explorerHTML)
		tmpl.Execute(w, nil)
	})

	http.ListenAndServe(":8080", nil)
}

const explorerHTML = `<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="/static/doc.css">
<link rel="stylesheet" href="https://code.getmdl.io/1.1.3/material.indigo-pink.min.css">
<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto">

<script defer src="https://code.getmdl.io/1.1.3/material.min.js"></script>
<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js?autoload=false" defer="defer"></script>
</head>
<body>
<div class="mdl-layout mdl-js-layout mdl-layout--fixed-drawer">
  <div class="mdl-layout__drawer">
    <span class="mdl-layout-title">Docs</span>
    <nav class="mdl-navigation">
	  <a href="#" class="mdl-navigation__link" onclick="update('..'); return false;">Back</a>
	  <div id="nav">
	  </div>
    </nav>
  </div>
  <main class="mdl-layout__content">
    <div class="page-content">
      <div id="output"></div>
	</div>
  </main>
</div>
<script src="/static/explorer.js"></script>
</body>
</html>`

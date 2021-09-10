package common

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type localFileSystem struct {
	fs    http.Handler
	dir   string
	index string
}

func StaticServe(urlPrefix string) gin.HandlerFunc {
	fileserver := LocalFile(urlPrefix)
	return func(c *gin.Context) {
		fileserver.ServeHTTP(c.Writer, c.Request)
	}
}

func LocalFile(dir string) *localFileSystem {
	return &localFileSystem{
		fs:    http.FileServer(http.Dir(dir)),
		dir:   dir,
		index: path.Join(dir, "index.html"),
	}
}

func (l localFileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if l.exists(r.URL.Path) {
		l.fs.ServeHTTP(w, r)
	} else {
		if strings.Contains(r.URL.Path, "/fail") || strings.Contains(r.URL.Path, "/success") {
			http.ServeFile(w, r, l.index)
		}
	}
}

func (l *localFileSystem) exists(filepath string) bool {
	name := path.Join(l.dir, filepath)
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

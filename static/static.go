package static

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

const INDEX = "index.html"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {

		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			ext := path.Ext(c.Request.URL.Path)
			usebr := true
			if ext != "" {
				switch ext {
				case ".html":
					c.Header("content-type", "text/html")
					break
				case ".js":
					c.Header("content-type", "text/javascript")
					break
				case ".png":
					c.Header("content-type", "image/png")
					break
				default:
					usebr = false
				}
			} else { //直接进主页
				c.Request.URL.Path = c.Request.URL.Path + "index.html" //虽然会本来就自带跳转,但是那样的话就没办法用br了
				c.Header("content-type", "text/html")
			}
			if usebr {
				c.Header("Content-Encoding", "br")  //声明压缩格式,否则会被当作二进制文件下载
				c.Header("Vary", "Accept-Encoding") //客户端使用缓存
				c.Request.URL.Path = c.Request.URL.Path + ".br"
			} else {
				fmt.Println(c.Request.URL.Path)
			}
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

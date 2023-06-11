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
			usebr(c)
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func usebr(c *gin.Context) {
	ext := path.Ext(c.Request.URL.Path)
	supportbr := strings.Contains(c.GetHeader("Accept-Encoding"), "br")
	if supportbr {
		if ext != "" {
			switch ext {
			case ".html":
				setbr(c, "text/html")
			case ".js":
				setbr(c, "text/javascript")
			case ".png":
				setbr(c, "image/png")
			default:
				fmt.Println(ext)
			}
		} else { //直接进主页
			c.Request.URL.Path = c.Request.URL.Path + "index.html" //虽然会本来就自带跳转,但是那样的话就没办法用br了
			setbr(c, "text/html")
		}
	}
}

func setbr(c *gin.Context, mime string) {
	c.Header("content-type", mime)     //压缩的文件必须显式指明mime type,否则会被当作二进制文件
	c.Header("Content-Encoding", "br") //声明压缩格式,否则会被当作二进制文件下载
	// c.Header("Vary", "Accept-Encoding") //客户端使用缓存,开发阶段先去掉
	c.Request.URL.Path = c.Request.URL.Path + ".br"
}

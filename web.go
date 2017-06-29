package swizzle

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var t = template.New("index")

func init() {
	t.Parse(`
<html>
	<h1>swizzle</h1>
	<h2>hello! select a function below:</h2>
	<ul>
		<li><a href="/ping">/ping</a>
		<li><a href="/echo?message=hello world">/echo</a>
		<li><a href="/pgsql">/pgsql</a>
		<li><a href="/redis">/redis</a>
		<li><a href="/status">/status</a>
		<li><a href="/s3">/s3</a>
	</ul>
</html>
	`)
}

// RunWeb start the web server
func RunWeb(config *Config) {
	r := gin.Default()
	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{})
	})

	// also add service endpoints
	addServiceEndpoints("web", config, r)

	// Listen on all IPs on the configured port
	r.Run(fmt.Sprintf(":%d", config.Port))
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/templates"
	"net/http"
)

func WebRoutes(r *core.Engine) {
	r.HTMLRender = templates.NewDefaultPongo2Render()

	r.GET("/aaa", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test/test.html", templates.Pongo2Ctx{
			"test": "test",
			"arr":  []string{"aaa", "bbb", "ccc"},
		})
	})
}

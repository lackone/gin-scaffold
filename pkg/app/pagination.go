package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/convert"
)

func GetPage(c *gin.Context) int {
	page := convert.Str(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}
	return page
}

func GetPageSize(c *gin.Context) int {
	engine := c.MustGet("engine").(*core.Engine)
	container := engine.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)

	pageSize := convert.Str(c.Query("page_size")).MustInt()
	if pageSize <= 0 {
		pageSize = configService.GetInt("app.default_page_size")
	}
	if pageSize > configService.GetInt("app.max_page_size") {
		pageSize = configService.GetInt("app.max_page_size")
	}
	return pageSize
}

func GetOffset(page, pageSize int) int {
	offset := 0
	if page > 0 {
		offset = (page - 1) * pageSize
	}
	return offset
}

package test

import (
	"github.com/gin-gonic/gin"
	"github.com/lackone/gin-scaffold/app/models"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/internal/provider/gorm"
)

type Test struct {
}

// TestSwag godoc
// @Summary 测试swag
// @Schemes
// @Description 测试swag
// @Tags TestSwag
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /test_swag [get]
func (t *Test) TestSwag(c *gin.Context) {
	c.JSON(200, "ok")
}

// TestSwag2 godoc
// @Summary 测试swag2
// @Schemes
// @Description 测试swag2
// @Tags TestSwag2
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /test_swag2 [get]
func (t *Test) TestSwag2(c *gin.Context) {
	c.JSON(200, "ok")
}

// TestORM godoc
// @Summary 测试orm
// @Schemes
// @Description 测试orm
// @Tags TestORM
// @Accept json
// @Produce json
// @Success 200 {string} ok
// @Router /test_orm [get]
func (t *Test) TestORM(c *gin.Context) {
	engine := c.MustGet("engine").(*core.Engine)
	container := engine.GetContainer()
	log := container.MustMake(contract.LogKey).(contract.Log)
	orm := container.MustMake(contract.GORMKey).(contract.IGORM)
	db, err := orm.GetDB(gorm.WithConfigPath("database.default"))
	if err != nil {
		log.Error(c, err.Error(), nil)
		c.Abort()
		return
	}
	db.WithContext(c)

	db.AutoMigrate(&models.User{})

	//增加记录
	u := &models.User{Name: "test", Pwd: "test"}
	err = db.Create(u).Error
	log.Info(c, "insert user", map[string]interface{}{
		"id":  u.Id,
		"err": err,
	})
	db.Create(&models.User{Name: "aaa", Pwd: "aaa"})
	db.Create(&models.User{Name: "bbb", Pwd: "bbb"})
	db.Create(&models.User{Name: "ccc", Pwd: "ccc"})

	//更新记录
	u.Name = "aaa"
	err = db.Save(u).Error
	log.Info(c, "update user", map[string]interface{}{
		"id":  u.Id,
		"err": err,
	})

	//查询数据
	q := &models.User{Id: 2}
	err = db.First(q).Error
	log.Info(c, "query user", map[string]interface{}{
		"id":  q.Id,
		"err": err,
	})

	//删除数据
	d := &models.User{Id: 3}
	err = db.Delete(d).Error
	log.Info(c, "delete user", map[string]interface{}{
		"id":  d.Id,
		"err": err,
	})

	c.JSON(200, "ok")
}

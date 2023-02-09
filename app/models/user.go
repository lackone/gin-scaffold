package models

type User struct {
	Id      int    `gorm:"primaryKey;auto_increment;column:id;type:int(10);not null;comment:ID"`
	Name    string `gorm:"column:name;type:varchar(32);not null;default:'';comment:名称"`
	Pwd     string `gorm:"column:pwd;type:varchar(32);not null;default:'';comment:密码"`
	Created int    `gorm:"autoCreateTime;column:created;type:int(10);not null;comment:创建时间"`
	Updated int    `gorm:"autoUpdateTime;column:updated;type:int(10);not null;comment:更新时间"`
}

package models

import (
	"gorm.io/gorm"
)

type User struct{
	gorm.Model //用模型本身的id
	Name 	string	`gorm:"unique"`
	Passwd	[]byte	//记得加盐
	Email	string	`gorm:"unique"`
}
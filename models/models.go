package models

import (
	"gorm.io/gorm"
)

type user struct{
	gorm.Model //用模型本身的id
	name 	string	`gorm:"unique"`
	passwd	[]byte	//记得加盐
	email	string	`gorm:"unique"`
}
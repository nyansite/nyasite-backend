package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        //用模型本身的id
	Name	string `gorm:"unique"`
	Passwd	[]byte //TODO 记得加盐
	Email	string `gorm:"unique"`
	Level	uint	//4位权限4位等级,所以满级15
}


//这个要重构,先摸了
type VideoPreviewRequire struct {
	gorm.Model
	CoverFile    string
	VideoFile    string
	Title        string
	Pass         uint
	Introduction string
}


type Video struct {
	gorm.Model
	VideoLink 	string
	CoverLink   string		//封面也用磁力链接
	Title       string
	Profile  	string		//芝士简介
	Comment		[]Comment	`gorm:"ForeignKey:Vid"`	//评论
	Tag			[]Tag		`gorm:"ForeignKey:Tid"`
	// Views         	uint	//这是什么
}

type Tag struct{
	ID		uint 	`gorm:"primarykey"`
	Vid		uint	`gorm:"index"`//对应的视频的id
	Tid		uint	`gorm:"index"`//避免tag文本被多次存储
}

type TagText struct{//tag的文本,其他地方有一个切片存储
	gorm.Model
	Text 	string	`gorm:"unique"`
}

type Comment struct {//不知道为什么基于gorm.Model会出问题
	gorm.Model
	Vid			uint	`gorm:"index"`	//所属的视频的id,帖子先摸了
	Cid			uint	`gorm:"index"`	//楼中楼上一层的ID,不是楼中楼应该为0
	Text		string
	/*
	文本类型
	0:	字符串
	1:	markdown
	2:	bbcode
	3:	reStructuredText
	*/
	Type 		uint8	
	Author  	uint
	Comment		[]Comment	`gorm:"ForeignKey:Cid"`
}

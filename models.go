package main

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusRepeatUserName   int = 601
	StatusRepeatEmail      int = 602
	StatusUserNameNotExist int = 611
	StatusPasswordError    int = 612
	StatusAlreadyLogin     int = 613
	StatusRepeatTag        int = 621
)

type User struct {
	gorm.Model        //用模型本身的id
	Name       string `gorm:"unique"`
	Passwd     []byte //TODO 记得加盐
	Email      string `gorm:"unique"`
	Level      uint   //4位权限4位等级,所以满级15
}

type Video struct {
	gorm.Model
	// VideoLink string 	//ipfs files 有文件名,可以指向uid,所以不需要这个了
	// ImgLink string
	Title       string
	Description string         //芝士简介
	CommentP    []VideoComment `gorm:"ForeignKey:Vid"` //评论
	Tag         []Tag          `gorm:"ForeignKey:Tid"`
	Views       uint           //这是播放量
}

type Tag struct {
	ID  uint `gorm:"primarykey"`
	Vid uint `gorm:"index"` //对应的视频的id
	Tid uint `gorm:"index"` //避免tag文本被多次存储
}

type TagText struct { //tag的文本,其他地方有一个切片存储
	ID   uint   `gorm:"primarykey"`
	Text string `gorm:"unique"`
}

type VideoComment struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Vid       uint `gorm:"index"` //所属页面的id
	Cid  uint `gorm:"index"`	//楼中楼上一层的id
	Text      string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type    uint8
	Author  uint
	Comment []VideoComment `gorm:"ForeignKey:Cid"`
}

// 论坛部分
// 不需要楼中楼,直接引用
type MainForum struct {
	gorm.Model
	ID     uint `gorm:"primarykey"`
	Title  string
	Author uint            `gorm:"index"` //发起人
	Views  uint            //阅读量
	UnitP  []UtilForum `gorm:"ForeignKey:Pid"`
}

type UtilForum struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	// Count		uint								//楼层
	Pid  uint `gorm:"index:UtilForum"` //所属论坛的id,楼中楼为0(大概)
	Cid  uint `gorm:"index:UtilForum"` //楼中楼上一层的ID,不是楼中楼应该为0
	Text string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type     uint8
	Author   uint
	CommentP []CommentPage `gorm:"ForeignKey:Pid"`
}

type CommentPage struct { //一页16个
	ID uint `gorm:"primarykey"`
	// Count   uint      //页数
	Comment []Comment `gorm:"ForeignKey:Pid"`
	Uid     uint      `gorm:"index"` //所属的视频/论坛的id
}
type Comment struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Pid       uint `gorm:"index:Comment"` //所属页面的id
	Cid       uint `gorm:"index:Comment"` //楼中楼上一层的ID
	Text      string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type   uint8
	Author uint
}

// 这个要重构,先摸了
type VideoPreviewRequire struct {
	gorm.Model
	CoverFile string
	VideoFile string
	Title     string
	Up        uint
	Pass      uint
	Profile   string
}

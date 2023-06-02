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
	Passwd     []byte
	Email      string `gorm:"unique"`
	Level      uint   //4位权限4位等级,所以满级15
}

type Video struct { //获取视频和获取评论分开
	gorm.Model
	// VideoLink string 	//ipfs files 有文件名,可以指向uid,所以不需要这个了
	// ImgLink string
	Title       string
	Description string         //芝士简介
	CommentP    []VideoComment `gorm:"ForeignKey:Vid"`   //评论
	Tag         []uint         `gorm:"index;type:bytes"` //tag的id
	likes       uint           //芝士点赞数量
	Views       uint           //这是播放量
}

type Tag struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Text      string `gorm:"unique"`
}

type VideoComment struct {
	gorm.Model
	Vid  uint `gorm:"index"` //所属页面的id
	Text string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type    uint8
	Author  uint
	likes   uint                //芝士点赞数量
	Comment []VideoCommentReply `gorm:"ForeignKey:Cid"`
}
type VideoCommentReply struct { //楼中楼的回复.......
	gorm.Model
	Cid  uint `gorm:"index"` //楼中楼上一层的id
	Text string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type    uint8
	Author  uint
	likes   uint                //芝士点赞数量
	Comment []VideoCommentReply `gorm:"ForeignKey:Cid"`
}

// 论坛部分
// 不需要楼中楼,直接引用
type MainForum struct { //获取视频和获取评论分开
	gorm.Model
	Title  string
	UnitP  []UnitForum `gorm:"ForeignKey:Mid"` //评论
	Views  uint        //这是播放量
	Author uint
}

type UnitForum struct {
	gorm.Model
	Mid  uint `gorm:"index"` //所属页面的id
	Cid  uint
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
	likes    uint      //芝士点赞数量
	CommentP []Comment `gorm:"ForeignKey:Cid"`
}
type Comment struct { //楼中楼的回复.......
	gorm.Model
	Cid  uint `gorm:"index"` //楼中楼上一层的id
	Text string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type   uint8
	Author uint
	likes  uint //芝士点赞数量
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

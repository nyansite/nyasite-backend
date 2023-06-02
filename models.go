package main

import (
	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int  //使用时间戳而非time.time(字符串)
	UpdatedAt int
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	StatusRepeatUserName   int = 601
	StatusRepeatEmail      int = 602
	StatusUserNameNotExist int = 611
	StatusPasswordError    int = 612
	StatusAlreadyLogin     int = 613
	StatusRepeatTag        int = 621
)

type User struct {
	Model         //用模型本身的id
	Name   string `gorm:"unique"`
	Passwd []byte
	Email  string `gorm:"unique"`
	Level  uint   //4位权限4位等级,所以满级15
}

type Video struct { //获取视频和获取评论分开
	Model
	// VideoLink string 	//ipfs files 有文件名,可以指向uid,所以不需要这个了
	// ImgLink string
	Title       string         `gorm:"default:芝士标题"`
	Description string         `gorm:"default:简介不见惹"`    //芝士简介
	CommentP    []VideoComment `gorm:"ForeignKey:Vid"`   //评论
	Tag         []uint         `gorm:"index;type:bytes"` //tag的id
	likes       uint           `gorm:"default:0"`        //芝士点赞数量
	Views       uint           `gorm:"default:0"`        //这是播放量
}

type Tag struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt int
	Text      string `gorm:"unique"`
}

type VideoComment struct {
	Model
	Vid  uint `gorm:"index"` //所属页面的id
	Text string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type    uint8               `gorm:"default:0"`
	Author  uint                //发表评论的用户
	likes   uint                //芝士点赞数量
	Comment []VideoCommentReply `gorm:"ForeignKey:Cid"`
}
type VideoCommentReply struct { //楼中楼的回复.......
	Model
	Cid  uint `gorm:"index"` //楼中楼上一层的id
	Text string
	/*
		文本类型
		0:	字符串
		1:	markdown
		2:	bbcode
		3:	reStructuredText
	*/
	Type    uint8 `gorm:"default:0"`
	Author  uint
	likes   uint                //芝士点赞数量
	Comment []VideoCommentReply `gorm:"ForeignKey:Cid"`
}

// 论坛部分
// 不需要楼中楼,直接引用
type MainForum struct {
	Model
	ID     uint `gorm:"primarykey"`
	Title  string
	Author uint        `gorm:"index"` //发起人
	Views  uint        //阅读量
	UnitP  []UtilForum `gorm:"ForeignKey:Pid"`
}

type UtilForum struct {
	ID uint `gorm:"primarykey"`

	// Count		uint				//楼层
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
	Model
	Pid  uint `gorm:"index:Comment"` //所属页面的id
	Cid  uint `gorm:"index:Comment"` //楼中楼上一层的ID
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
}

// 这个要重构,先摸了
type VideoPreviewRequire struct {
	Model
	CoverFile string
	VideoFile string
	Title     string
	Up        uint
	Pass      uint
	Profile   string
}

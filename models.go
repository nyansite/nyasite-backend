package main

import (
	_ "xorm.io/xorm"
)

type Model struct {
	ID        uint //xorm自动主键
	CreatedAt int  `xorm:"created"` //使用时间戳而非time.time(字符串)
	UpdatedAt int  `xorm:"updated"`
	DeletedAt int `xorm:"deleted"`
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
	Name   string `xorm:"unique"`
	Passwd []byte
	Email  string `xorm:"unique"`
	Level  uint8  `xorm:"default:0"` //4位权限4位等级,所以满级15(要不了这么多)
}

type Video struct { //获取视频和获取评论分开
	Model
	Title       string         `xorm:"default:芝士标题"`
	Description string         `xorm:"default:简介不见惹"`
	CommentP    []VideoComment `gorm:"ForeignKey:Vid"`   //评论	
	Tag         []uint         `xorm:"index"` //tag的id
	likes       uint           `xorm:"default:0"`        //芝士点赞数量
	Views       uint           `xorm:"default:0"`        //这是播放量
	Author      uint           `xorm:"index"`            //作者/上传者
}

// Tags切片存一样的数据
type Tag struct {
	ID        uint //xorm自动主键
	CreatedAt int  `xorm:"created"`
	Text      string `gorm:"unique"`
}

type VideoComment struct {
	Model
	Vid     uint `xorm:"index"` //所属页面的id
	Text    string
	IsMD    bool                `xorm:"default:false"` //t:markdown,f:str
	Author  uint                `xorm:"index"`         //发表评论的用户
	likes   uint                `xorm:"default:0"`     //芝士点赞数量
	Comment []VideoCommentReply `gorm:"ForeignKey:Cid"`
}

type VideoCommentReply struct { //楼中楼的回复.......
	Model
	Cid     uint `xorm:"index"` //楼中楼上一层的id,自动生成
	Text    string
	IsMD    bool                `xorm:"default:false"` //t:markdown,f:str
	Author  uint                `xorm:"index"`
	likes   uint                `xorm:"default:0"`      //芝士点赞数量
	Comment []VideoCommentReply `gorm:"ForeignKey:Cid"` //楼中楼中楼...
}

// 论坛部分
// 不需要楼中楼,直接引用

type Forum struct { //获取视频和获取评论分开
	Model
	Title   string
	Comment []ForumComment `gorm:"ForeignKey:Mid"` //评论
	Views   uint           `xorm:"default:0"`      //这是播放量
	Author  uint           `xorm:"index"`
}

type ForumComment struct {
	Model
	Mid    uint `xorm:"index"` //所属页面的id
	Text   string
	IsMD   bool `xorm:"default:false"` //t:markdown,f:str
	Author uint `xorm:"index"`
	// Emoji  []uint 	//先摸了
}

// 正在转码压制中的视频
type VideoTranscode struct {
	UUID      string
	Author    uint
	CreatedAt int
}

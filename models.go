package main

import (
	_ "xorm.io/xorm"
)

const (
	StatusRepeatUserName   int = 601
	StatusRepeatEmail      int = 602
	StatusUserNameNotExist int = 611
	StatusPasswordError    int = 612
	StatusAlreadyLogin     int = 613
	StatusRepeatTag        int = 621
)

type Model struct {
	Id        int64 //xorm自动主键
	CreatedAt int   `xorm:"created"` //使用时间戳而非time.time(字符串)
	UpdatedAt int   `xorm:"updated"`
	DeletedAt int   `xorm:"deleted"` //用模型本身的id
}
type User struct {
	Model  `xorm:"extends"`
	Name   string `xorm:"unique"`
	Passwd []byte
	Email  string `xorm:"unique"`
	Level  uint8  `xorm:"default 0"` //4位权限4位等级,所以满级15(要不了这么多)
}

type Video struct { //获取视频和获取评论分开
	Model       `xorm:"extends"`
	Title       string `xorm:"default '芝士标题'"`
	Description string `xorm:"default '简介不见惹'"`
	Tag         []uint `xorm:"index"`     //tag的id
	likes       uint   `xorm:"default 0"` //芝士点赞数量
	Views       uint   `xorm:"default 0"` //这是播放量
	Author      uint   `xorm:"index"`     //作者/上传者
}

type Tag struct {
	Id        int64  //xorm自动主键
	CreatedAt int    `xorm:"created"`
	Text      string `xorm:"unique"`
}

type VideoComment struct {
	Model  `xorm:"extends"`
	Vid    uint `xorm:"index"` //所属页面的id
	Text   string
	Author uint `xorm:"index"`         //发表评论的用户
	likes  uint `xorm:"default 0"`     //芝士点赞数量
}

type VideoCommentReply struct { //楼中楼的回复.......
	Model  `xorm:"extends"`
	Cid    uint `xorm:"index"` //楼中楼上一层的id,自动生成
	Text   string
	Author uint `xorm:"index"`
	likes  uint `xorm:"default 0"` //芝士点赞数量
}

// 论坛部分
// 不需要楼中楼,直接引用

type Forum struct { //获取视频和获取评论分开
	Model   `xorm:"extends"`
	Title   string
	Comment []ForumComment `xorm:"extends"`   //评论
	Views   uint           `xorm:"default 0"` //这是播放量
	Author  uint           `xorm:"index"`
}

type ForumComment struct {
	Model  `xorm:"extends"`
	Mid    uint64 `xorm:"index"` //所属论坛的id
	Text   string
	Author uint64   `xorm:"index"`
	Emoji  []uint //先摸了
}

// 正在转码压制中的视频
type VideoTranscode struct {
	UUID      string
	Author    uint
	CreatedAt int
}

// Session的密钥
type SessionSecret struct {
	CreatedAt      int64 `xorm:"created unique pk"` //没错芝士主键
	Authentication []byte
	Encryption     []byte
}

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
	Title       string `xorm:"default '芝士标题'"`
	Description string `xorm:"default '简介不见惹'"`
	Tag         []uint `xorm:"index"`     //tag的id
	likes       uint   `xorm:"default 0"` //芝士点赞数量
	Views       uint   `xorm:"default 0"` //这是播放量
	Author      uint   `xorm:"index"`     //作者/上传者
	Model       `xorm:"extends"`
}

type Tag struct {
	Id        int64  //xorm自动主键
	CreatedAt int    `xorm:"created"`
	Text      string `xorm:"unique"`
}

type VideoComment struct {
	Vid    uint `xorm:"index"` //所属页面的id
	Text   string
	Author uint `xorm:"index"`     //发表评论的用户
	likes  uint `xorm:"default 0"` //芝士点赞数量
	Model  `xorm:"extends"`
}

type VideoCommentReply struct { //楼中楼的回复.......
	Cid    uint `xorm:"index"` //楼中楼上一层的id,自动生成
	Text   string
	Author uint `xorm:"index"`
	likes  uint `xorm:"default 0"` //芝士点赞数量
	Model  `xorm:"extends"`
}

// 论坛部分
// 不需要楼中楼,直接引用

type Forum struct { //获取视频和获取评论分开
	Title  string
	Views  uint `xorm:"default 0"` //这是浏览量
	Author uint `xorm:"index"`
	Kind   uint8
	//0:官方通知区;1:用户反馈区;2:结束的用户反馈区;3:Thread贴;4:完结的Thread贴;5:资源贴
	Model `xorm:"extends"`
}

type ForumComment struct {
	Mid         uint `xorm:"index"` //所属论坛的id
	Text        string
	Author      uint `xorm:"index"`
	Like        uint //uint8只有255，可能不太够用
	Dislike     uint
	Smile       uint
	Celebration uint
	Confused    uint
	Heart       uint
	Rocket      uint
	Eyes        uint
	Model       `xorm:"extends"`
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

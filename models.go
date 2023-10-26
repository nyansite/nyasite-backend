package main

/*
重要!!!!!!!!!!!!!!
psql无uint
*/
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
	Avatar string `xorm:"TEXT"`
	Level  uint8  `xorm:"TINYINT default 0"` //4位权限4位等级,所以满级15(要不了这么多)
}

// 标签部分
type TagModel struct {
	Id    int64  //xorm自动主键
	Text  string `xorm:"unique"`
	Times int
}

type Tag struct {
	Id  int64 //xorm自动主键
	Tid int   `xorm:"SMALLINT"`
	Pid int
}

// 视频部分
type Video struct { //获取视频和获取评论分开
	CoverPath   string
	Title       string `xorm:"default '芝士标题'"`
	Description string `xorm:"default '简介不见惹'"`
	Likes       int    `xorm:"default 0"` //芝士点赞数量
	Views       int    `xorm:"default 0"` //这是播放量
	Author      int    `xorm:"index"`     //作者/上传者
	Model       `xorm:"extends"`
}

type VideoComment struct {
	Vid    int `xorm:"index"` //所属页面的id
	Text   string
	Author int `xorm:"index"`     //发表评论的用户
	Likes  int `xorm:"default 0"` //芝士点赞数量
	Model  `xorm:"extends"`
}

type VideoCommentReply struct { //楼中楼的回复.......
	Cid    int `xorm:"index"` //楼中楼上一层的id,自动生成
	Text   string
	Author int `xorm:"index"`
	Likes  int `xorm:"default 0"` //芝士点赞数量
	Model  `xorm:"extends"`
}

type VideoLikeRecord struct {
	Author int
	Vid    int
}

// 论坛部分
// 不需要楼中楼,直接引用

type Forum struct { //获取视频和获取评论分开
	Title  string
	Views  int `xorm:"default 0"` //这是浏览量
	Author int `xorm:"index"`
	Kind   int
	//0:官方通知区;1:用户反馈区;2:关闭的用户反馈区;3:Thread贴;4:完结的Thread贴;5:资源贴
	Model `xorm:"extends"`
}

type ForumComment struct {
	Mid         int `xorm:"index"` //所属论坛的id
	Text        string
	Author      int `xorm:"index"`
	Choose      int `xorm:"-"`
	Like        int //uint8只有255，可能不太够用
	Dislike     int
	Smile       int
	Celebration int
	Confused    int
	Heart       int
	Rocket      int
	Eyes        int
	Model       `xorm:"extends"`
}

type EmojiRecord struct {
	Author int
	Uid    int
	Emoji  int
}

// 正在转码压制中的视频
type VideoTranscode struct {
	UUID      string
	Author    int
	CreatedAt int
}

// Session的密钥
type SessionSecret struct {
	CreatedAt      int64 `xorm:"created unique pk"` //没错芝士主键
	Authentication []byte
	Encryption     []byte
}

//搜索部分

type SearchFourmReturn struct {
	Id    int64
	Title string
	Text  string
	Kind  int
}

type SearchVideoReturn struct {
	Id        int64
	CoverPath string
	Title     string
	Views     int
}

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

// 用户部分
type User struct {
	Name   string `xorm:"unique"`
	Passwd []byte
	Email  string `xorm:"unique"`
	Avatar string
	Level  uint8 `xorm:"TINYINT default 0"` //4位权限4位等级,所以满级15(要不了这么多)
	Model  `xorm:"extends"`
}

type UserDataShow struct {
	Name   string
	Avatar string
	Id     int64
}

// 标签部分
type TagModel struct {
	Id    int64  //xorm自动主键
	Kind  uint8  `xorm:"SMALLINT"`
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
	VideoUid    string //视频路径,封面路径
	CoverPath   string
	Title       string
	Description string `xorm:"TEXT"`
	Likes       int    `xorm:"default 1"` //芝士点赞数量
	Views       int    `xorm:"default 1"` //这是播放量
	Author      int    `xorm:"index"`     //作者/上传者
	Model       `xorm:"extends"`
}

type VideoNeedToCheck struct {
	VideoUid    string
	CoverPath   string
	Title       string
	Description string `xorm:"TEXT"`
	Tags        []uint8
	Author      int `xorm:"index"`
	Model       `xorm:"extends"`
}

// 视频站评论部分
type VideoComment struct {
	Vid         int                 `xorm:"index"` //所属视频的id
	Text        string              `xorm:"TEXT"`
	Author      int                 `xorm:"index"`
	Choose      int8                `xorm:"-"`
	CRdisplay   []VideoCommentReply `xorm:"-"`         //CR = CommentReply
	Like        int                 `xorm:"default 1"` //uint8只有255，可能不太够用
	Dislike     int                 `xorm:"default 1"`
	Smile       int                 `xorm:"default 1"`
	Celebration int                 `xorm:"default 1"`
	Confused    int                 `xorm:"default 1"`
	Heart       int                 `xorm:"default 1"`
	Rocket      int                 `xorm:"default 1"`
	Eyes        int                 `xorm:"default 1"`
	Model       `xorm:"extends"`
}

type VideoCommentReply struct { //楼中楼的回复.......
	Cid    int    `xorm:"index"` //楼中楼上一层的id,自动生成
	Text   string `xorm:"TEXT"`
	Author int    `xorm:"index"`
	Likes  int    `xorm:"default 1"` //芝士点赞数量
	Like_c bool   `xorm:"-"`         //判断是否点赞
	Model  `xorm:"extends"`
}

type VideoCommentEmojiRecord struct {
	Author int
	Cid    int
	Emoji  int8 `xorm:"TINYINT"`
}

type VideoCommentReplyLikeRecord struct {
	Author int
	Crid   int //CR = CommentReply
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

type SearchVideoReturn struct {
	Id        int64
	CoverPath string
	Title     string
	Views     int
}

// 消息部分
type MessageVideoEmojiComment struct {
	Vid     int64
	Reciver int
	Cliker  int
	Kind    int8
}

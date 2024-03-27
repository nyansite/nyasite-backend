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

type ModelLight struct {
	Id        int64 //xorm自动主键
	UpdatedAt int   `xorm:"updated"`
}

// 用户部分
type User struct {
	Name   string `xorm:"unique"`
	Passwd []byte
	Email  string `xorm:"unique"`
	Avatar string
	Level  uint8 `xorm:"TINYINT default 0"` //4位权限4位等级,所以满级15(要不了这么多)
	LTCM   int   //LastTimeCheckMessage
	LTC    int   //LastTimeClockIn
	Model  `xorm:"extends"`
}

type UserDataShow struct {
	Name   string
	Avatar string
	Id     int64
}

// 社团
type ApplyCircleVote struct {
	IsAgree bool
	Acid    int
	Author  int
}

type ApplyCircle struct {
	Name       string `xorm:"unique"`
	Avatar     string
	Descrption string `xorm:"TEXT"` //markdown
	ApplyText  string `xorm:"TEXT"`
	Stauts     bool   //false:审核中 true:驳回
	Kinds      int16  `xorm:"SMALLINT"` //1st bite: video,2nd bite: image,3rd bite: music.
	Applicant  int
	ModelLight `xorm:"extends"`
}

type VoteOfApplyCircle struct {
	Reviewer int
	Agree    bool
	Acid     int
}

type Circle struct {
	Name       string `xorm:"unique"`
	Avatar     string
	Descrption string `xorm:"TEXT"`
	Kinds      int16  `xorm:"SMALLINT"` //1:video 2:music 4:image
	Model      `xorm:"extends"`
}

type MemberOfCircle struct {
	Uid        int   //User.Id
	Cid        int   //Circle.Id
	Permission uint8 `xorm:"TINYINT"` //0:Subscribe,1:Staff,2:Creator,3:Maintainer,4:Owner
	UpdatedAt  int   `xorm:"updated"`
}

type CircleDataShow struct {
	Name     string
	Avatar   string
	Relation int8
	Id       int64
}

type UserDataShowWithPermission struct {
	Name       string
	Avatar     string
	Permission uint8
	Id         int64
}

// 标签部分
type TagModel struct {
	Id       int64 //xorm自动主键
	ParentId int
	Kind     uint8  `xorm:"SMALLINT"`
	Text     string `xorm:"unique"`
	Times    int
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
	Upid        int    //上传用户
	Caution     int
	Model       `xorm:"extends"`
}

type VideoNeedToCheck struct {
	OriginalId  int
	VideoUid    string
	CoverPath   string
	Title       string
	Description string `xorm:"TEXT"`
	Tags        []int
	Author      int `xorm:"index"`
	Upid        int //上传用户
	Stauts      bool
	Reason      string
	ModelLight  `xorm:"extends"`
}

type VideoLikeRecord struct {
	Author int
	Vid    int
}

type VideoPlayedRecord struct {
	Uid      int
	Vid      int
	LastPlay int
}

type VideoReturn struct {
	Id         int64
	CoverPath  string
	Title      string
	Author     SearchCircleReturn
	Views      int
	Likes      int
	SelfUpload bool
	CreatedAt  int
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

// 弹幕部分
type VideoBullet struct {
	Author int
	Vid    int
	Text   string
	Time   float64
	Color  string
	Top    bool
	Bottom bool
	Force  bool
	Model  `xorm:"extends"`
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

type SearchCircleReturn struct {
	Id   int64
	Name string
}

// 消息部分

type Invitation struct {
	Inviter   int
	Invitee   int
	Circle    int
	Kind      uint8
	Stauts    bool
	CreatedAt int `xorm:"created"`
}

type Discharge struct {
	WhoDischarge     int
	WhoBeDischargeed int
	Circle           int
	CreatedAt        int `xorm:"created"`
}

type CircleAffairMessage struct {
	Kind uint8
	// 0:join 1.discharge 2.quit
	//(3.invite staff 4.invite creator 5.invite maintianer) self 6.invite staff 7.invite creator 8.invite maintianer
	//9.reject circle
	SenderId    int
	SenderName  string
	ReciverId   int
	ReciverName string
	CircleId    int
	CircleName  string
	Time        int
	Id          int
}

type CheckMessage struct {
	Kind uint8
	//0:pass video 1.reject video 3.pass circle 4.reject circle
	Name   string
	Reason string
	DBId   int //Actual id in database
	Image  string
	Time   int
	Id     int //Key for rendering list
}

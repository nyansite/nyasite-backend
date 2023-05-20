package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        //用模型本身的id
	Name       string `gorm:"unique"`
	Passwd     []byte //TODO 记得加盐
	Email      string `gorm:"unique"`
	Level      uint   //4位权限4位等级,所以满级15
}

// 这个要重构,先摸了
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
	VideoLink string
	CoverLink string //封面也用磁力链接
	Title     string
	Profile   string    //芝士简介
	Comment   []Comment `gorm:"ForeignKey:Vid"` //评论
	Tag       []Tag     `gorm:"ForeignKey:Tid"`
	Views     uint      //这是播放量
}

type Tag struct {
	ID  uint `gorm:"primarykey"`
	Vid uint `gorm:"index"` //对应的视频的id
	Tid uint `gorm:"index"` //避免tag文本被多次存储
}

type TagText struct { //tag的文本,其他地方有一个切片存储
	gorm.Model
	Text string `gorm:"unique"`
}

type Tag struct {
	ID  uint `gorm:"primarykey"`
	Vid uint `gorm:"index"` //对应的视频的id
	Tid uint `gorm:"index"` //避免tag文本被多次存储
}

//论坛部分

type MainPost struct {
	gorm.Model
	Title       string `json:"title"`
	User_p      string `json:"user_p"` //发起人
	Views       uint   `json:"views"`
	Likes       uint   `json:"likes"`
	Video_p     string `json:"video_p"` //如果帖子是视频的评论区存储视频id，如果不是存储"independent"
	ContentShow string `json:"contentshow"`
}

type UnitPost struct {
	gorm.Model
	MainPost_p string `json:"mainpost_p"`
	Content    string `json:"content"` //以md形式储存
	User_p     string `json:"user_p"`
}

// 一个mainpost下面挂着unitpost
type Comment struct {
	gorm.Model
	Post_p  string `json:"post_p"`
	Content string `json:"content"`
	User_p  string `json:"user_p"`
}

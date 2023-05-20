package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        //用模型本身的id
	Name       string `gorm:"unique"`
	Passwd     string //TODO 记得加盐
	Email      string `gorm:"unique"`
}

//视频部分

type VideoRequireReview struct {
	gorm.Model
	CoverFile    string `json:"cover`
	VideoFile    string `json:"videofile"`
	UpId_p       string `json:"up_p"`
	Title        string `json:"title"`
	Pass         uint   `json:"pass"`
	Introduction string `json:"introduction"`
}

type Video struct {
	gorm.Model
	VideoIpfsSite string `json:"VideoIpfsSite"`
	CoverFile     string `json:"coverfile"`
	UpId_p        string `json:"up_p"`
	Title         string `json:"title"`
	Introduction  string `json:"introduction"`
	Views         uint   `json:"views"`
}

//论坛部分

type MainPost struct {
	gorm.Model
	Title       string `json:"title"`
	User_p      string `json:"user_p"` //发起人
	Views       uint   `json:"views"`
	Video_p     string `json:"video_p"` //如果帖子是视频的评论区存储视频id，如果不是存储"independent"
	ContentShow string `json:"contentshow"`
}

type UnitPost struct {
	gorm.Model
	MainPost_p string `json:"mainpost_p"`
	Content    string `json:"content"` //以md形式储存
	User_p     string `json:"user_p"`
	Level      uint   `json:"level"` //楼层
}

// 一个mainpost下面挂着unitpost

type Comment struct {
	gorm.Model
	Post_p  string `json:"post_p"`
	Content string `json:"content"`
	User_p  string `json:"user_p"`
}

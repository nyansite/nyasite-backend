package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        //用模型本身的id
	Name	string `gorm:"unique"`
	Passwd	[]byte //TODO 记得加盐
	Email	string `gorm:"unique"`
	Level	uint	//4位权限4位等级,所以满级15
}

//视频部分

type VideoPreviewRequire struct {
	gorm.Model
	CoverFile    string
	VideoFile    string
	Title        string
	Pass         uint
	Introduction string
}

//说好的不用ipfs呢??????
type Video struct {
	gorm.Model
	VideoLink 	string
	CoverLink   string		//封面也用磁力链接
	Title       string
	Profile  	string		//芝士简介
	Comment		[]Comment	//评论	
	// Views         	uint	//这是什么
}

//论坛部分
//TDOD 这块要重构
// type MainPost struct {
// 	gorm.Model
// 	Title       string `json:"title"`
// 	User_p      string `json:"user_p"` //发起人
// 	Views       uint   `json:"views"`
// 	Video_p     string `json:"video_p"` //如果帖子是视频的评论区存储视频id，如果不是存储"independent"
// 	ContentShow string `json:"contentshow"`
// }

// type UnitPost struct {
// 	gorm.Model
// 	MainPost_p string `json:"mainpost_p"`
// 	Content    string `json:"content"` //以md形式储存
// 	User_p     string `json:"user_p"`
// 	Level      uint   `json:"level"` //楼层
// }

// 一个mainpost下面挂着unitpost


type Comment struct {
	gorm.Model
	Text	string
	Type 	uint8	//0为字符串,1为markdown
	Author  uint
}

# cute_site

正在开发的后端

## api

芝士api

### api/user_status

获取已登录用户自身的信息

### api/user_status/:id

根据id获取其他用户信息

### api/register

注册;传入表单;属性是username/passwd/email

### api/login

登录;传入表单;属性是username/passwd,用户名也有可能是邮箱

### api/all_forum/:page

查看所有帖子;传入伪静态链接;属性是page,20个帖子一page,page决定查看论坛的区间

### api/browse_unitforum/:mid/:page

查看单个帖子下跟帖;传入伪静态链接;属性是mid,page,20个跟帖一page,page决定查看论坛的区间

### api/browse_forum/:board/:page

根据分区查看帖子;传入伪静态链接;属性是board,page,board决定板块,20个跟帖一page,page决定查看论坛的区间

| board | 含意      |
| ----- | ------- |
| 0     | 官方通知区   |
| 1     | 用户反馈区   |
| 2     | Thread贴 |
| 3     | 资源贴     |

## uapi

这是不安全的api,必须登录

### Video

#### uapi/upload_video

上传视频;等级限制为;传入表单;属性是Title,Description(简介),file(视频文件),cover(封面文件)

#### upi/add_video_comment

提交视频评论;等级限制为;传入表单;属性是vid(评论视频的vid),Text

### Forum

#### uapi/add_mainforum

发帖;无等级限制;传入表单;属性是title,text,type(决定发送分区)

| type | 含意        |
| ---- | --------- |
| 0    | 用户反馈区     |
| 1    | 创作Thread区 |
| 2    | 资源站       |
| 3    | 灌水区       |

#### uapi/add_unitforum

跟帖;无等级限制;传入表单;属性是uid(被跟的帖子的id),text

#### uapi/add_emoji

添加表情评论;无等级限制;传入表单;属性是uid(被跟的帖子的id),emoji(添加表情的编号)(详见表情符号)

#### uapi/finish_forum

完结帖子;需要作者本人;传入表单;属性是mid(完结的帖子)

## 自定义的状态码

### 6xx

| 状态码 | 含意        |
| --- | --------- |
| 601 | 用户名重复     |
| 602 | 邮箱重复      |
| 611 | 用户名或邮箱不存在 |
| 612 | 密码错误      |
| 613 | 重复登录      |
| 621 | tag重复     |

成功返回标准的200

## 权限

0-15

| 等级    | 含意           |
| ----- | ------------ |
| 0     | 什么都不能干的未答题用户 |
| 1     | 答完题的用户       |
| 10-14 | 普通管理员        |
| 15    | 超级管理员\开发者    |

## ipfs files(目录)

| 路径     | 作用  | 备注                           |
| ------ | --- | ---------------------------- |
| /img   | 图片  | {%videoid}.webp, webp+brotli |
| /video | 视频  | {%videoid}.m3u8, av1+brotli  |

<!-- |/temporary|临时文件|
|/temporary/video/{%date}/{%uuid}|未审核的视频| -->

## session

| 键        | 说明                                |
| -------- | --------------------------------- |
| is_login | 登录为true                           |
| userid   |                                   |
| level    | 虽然存储的是字符串但是请当作uint8使用awa 4b权限4b经验 |

## 帖子类型（forum.kind)

| value | 含意         |
| ----- | ---------- |
| 0     | 官方通知区      |
| 1     | 用户反馈区      |
| 2     | 关闭的用户反馈区   |
| 3     | Thread贴    |
| 4     | 完结的Thread贴 |
| 5     | 资源帖        |

## 表情编号

因为xorm的特性，表情记录从1开始，故每次调用都要-1

| emoji | 含意             |
| ----- | -------------- |
| 0     | Like 👍        |
| 1     | Dislike 👎     |
| 2     | Smile 😄       |
| 3     | Celebration 🎉 |
| 4     | Confused 😕    |
| 5     | Heart ❤️       |
| 6     | Rocket 🚀      |
| 7     | Eyes 👀        |
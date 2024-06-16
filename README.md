# cute_site

正在开发的后端

## api/

### 用户

#### api/user_status

GET方法

获取已登录用户自身的信息

无提交

返回：

| Key    | Explain    |
| ------ | ---------- |
| name   | 用户名称       |
| userid | 用户id       |
| mail   | 用户邮箱       |
| level  | 用户等级（详见等级） |
| avatar | 用户头像链接     |

#### api/user_status/:id

GET方法

根据id获取其他用户信息

以伪静态链接形式传入查询用户的id（:id)

返回：

| Key    | Explain |
| ------ | ------- |
| name   | 用户名称    |
| level  | 用户等级    |
| avatar | 用户头像链接  |

#### api/register

POST方法

注册

传入表单：

| Key      | Explain |
| -------- | ------- |
| username | 用户名     |
| passwd   | 密码      |
| email    | 电子邮件地址  |

无返回

#### api/login

POST方法

登录

传入表单：

| Key      | Explain |
| -------- | ------- |
| username | 用户名或者邮箱 |
| passwd   | 密码      |

创建session cookie

#### api/refresh

GET方法

更新cookie的过期日期

#### api/

## video

##### author 子结构体

创作视频的社团

| Key      | Explain   |
| -------- | --------- |
| Name     | 社团名称      |
| Avatar   | 社团头像      |
| Relation | 访问者和社团的关系 |
| Id       | 社团id      |

##### videoReturn 子结构体

   返回的视频数据

| Key        | Explain             |
| ---------- | ------------------- |
| Id         | 视频id                |
| CoverPath  | 封面链接                |
| Title      | 标题                  |
| Author     | 创作的社团（一个author子结构体） |
| Views      | 播放量                 |
| Likes      | 点赞量                 |
| Marks      | 收藏量                 |
| SelfUpload | 是否为自己上传             |
| CreatedAt  | 视频的上传时间             |

#### api/get_video/:id

GET方法

获取视频信息

以伪静态链接形式传入查询视频的id（:id)

返回：

| Key         | Explain             |
| ----------- | ------------------- |
| title       | 视频的标题               |
| videoPath   | 视频的链接               |
| author      | 创作的社团（一个author子结构体） |
| creatTime   | 上传时间                |
| description | 简介                  |
| views       | 播放量                 |
| likes       | 点赞量                 |
| isLiked     | 是否已经点赞              |
| marks       | 收藏量                 |
| isMarked    | 是否已经收藏              |

#### api/upload_video

POST方法

上传视频

传入表单：

| Key         | Explain                 |
| ----------- | ----------------------- |
| author      | 创作社团的id                 |
| title       | 视频标题                    |
| description | 视频简介                    |
| cover       | 封面链接                    |
| tags        | 视频标签对应tagModel的id的array |

## 注册和登录的错误代码

### 注册

| 错误字符串            | 含义     |
| ---------------- | ------ |
| NameUsed         | 重复用户名  |
| EmailAddressUsed | 重复邮箱地址 |

### 登录

对登录而言200表示登录成功

| 错误代码 | 含义           |
| ---- | ------------ |
| 400  | 用户名或密码为空     |
| 401  | 用户名不正确或密码不存在 |
| 403  | 已经登录了        |

## 等级

0-15

| 等级    | 含意           |
| ----- | ------------ |
| 0     | 什么都不能干的未答题用户 |
| 1     | 答完题的用户       |
| 10-14 | 普通管理员        |
| 15    | 超级管理员\开发者    |

## session

| 键      | 说明                          |
| ------ | --------------------------- |
| userid |                             |
| level  | 4bite权限4bite经验              |
| pwd-8  | 密码的最后8位,校验用,别问我为什么6位密码有8位数据 |

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
| 1     | Like 👍        |
| 2     | Dislike 👎     |
| 3     | Smile 😄       |
| 4     | Celebration 🎉 |
| 5     | Confused 😕    |
| 6     | Heart ❤️       |
| 7     | Rocket 🚀      |
| 8     | Eyes 👀        |
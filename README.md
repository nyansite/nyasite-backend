# cute_site

正在开发的后端

## api

芝士api

#### api/user_status

获取已登录用户自身的信息，同时刷新cookie

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

根据id获取其他用户信息

以伪静态链接形式传入查询用户的id（:id)

返回：

| Key    | Explain |
| ------ | ------- |
| name   | 用户名称    |
| level  | 用户等级    |
| avatar | 用户头像链接  |

#### api/register

注册;传入表单 

提交：

| Key      | Explain |
| -------- | ------- |
| username | 用户名     |
| passwd   | 密码      |
| email    | 电子邮件地址  |

无返回

#### api/login

登录;传入表单

提交：

| Key      | Explain |
| -------- | ------- |
| username | 用户名或者邮箱 |
| passwd   | 密码      |

## video

#### api/get_video/:id

获取视频信息

以伪静态链接形式传入查询视频的id（:id)

返回：

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

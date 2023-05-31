# cute_site

正在开发的后端

## api

芝士api

### api/user_status

获取已登录用户自身的信息

### api/user_status/:id

根据id获取其他用户信息

### api/register

注册,传入表单,属性是username/passwd/email

### api/login

登录,传入表单,属性是username/passwd,用户名也有可能是邮箱

## 自定义的状态码

### 6xx

|状态码|含意|
|---|---|
|601|用户名重复|
|602|邮箱重复|
|611|用户名或邮箱不存在|
|612|密码错误|
|613|重复登录|
|621|tag重复|

成功返回标准的200

## 权限

0-15

|等级|含意|
|---|---|
|0|什么都不能干的未答题用户
|1|答完题的用户
|10-14|普通管理员|
|15|超级管理员\开发者|

## ipfs files(目录)

|路径|作用|备注
|---|---|---|
|/img|图片|{%videoid}.webp, webp+brotli|
|/video|视频|{%videoid}.m3u8, av1+brotli|
|/temporary|临时文件|
|/temporary/video/{%date}/{%uuid}|未审核的视频|

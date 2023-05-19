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

## 自定义的状态码

### 6xx

|状态码|含意|
|---|---|
|601|创建成功|
|602|用户名重复|
|603|邮箱重复|
|611|登录成功|
|612|用户名或邮箱不存在|
|613|密码错误|

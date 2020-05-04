# 快速上手
### myway-mock是什么，能干什么
#### 是什么？
一个跨平台的应用程序，用于模拟http服务端接口。
不需要懂服务应用的开发，只需要知道接口约定，就可以获得一个模拟的服务接口。
它能做参数校验，根据不同的参数返回不同的报文，能模拟耗时，让你不用做任何设计开发上的改动， 
就可以很好的完成开发、测试、演示等工作。
#### 用于哪些场景
* 环境复杂难于创建（比如复杂的测试环境）
* 接口服务不易获取（复杂的获取流程）
* 接口服务使用方很多，或频繁发版不稳定（几率性的获取失败）
* 前后端分离（前后依赖是并行任务）  
* 由于时间资源缘故对应的接口服务不可用
* 基于契约编程，项目依赖的服务接口的进度明显落后于本项目
* 对外联调接口演示
* 性能测试中模拟外部或无测试环境的服务接口

### 下载安装
#### 下载地址

#### 安装方法
将对应的压缩包解压到想安装的地方。
* windows 直接启动即可
* linux、mac系统 请给myway-mock运行权限
### 运行
执行myway-mock命令，模拟服务就启动了，压缩包带有demo api，只要启动就可以看到如下输出。
默认使用8080，9090端口，如果程序直接退出，说明端口被占用。你可以修改配置文件config.yaml
进行端口修改。  
可以在浏览器直接访问 localhost:8080/helloworld?name=mike,进行检测看模拟服务是否运行正常。

### 第一个模拟API
新建myfirst.yaml文件，并放在apis目录下。  
文件内容如下：
```yaml
name: myfirst  
url: /first 
description: 我的第一个api
# 这个api支持GET和POST请求方式
methods:  
  - GET  
  - POST
# 不设置模拟耗时
delay:  
# 请求参数定义
request:
# 请求格式是json形式的    
  style: json  
# 当参数错误的时候的返回，当然也可以不定义
  error: |  
    {"success":false,"errcode":"50001","errmessage":"参数{{ ._Parameter}},校验错误 {{._Msg}}"}  
# 请求参数定义
  items:  
    - {paramter: name ,type: String ,length: 30,policy: Must,validator: }  
# 返回值定义
response:  
#默认返回结果编号  
  default: success  
#触发条件，根据条件返回不同的数据
  triggers:  
    - data: error 
# 当name==jone的时候返回error对应的结果。可以多个条件，当有多个的时候表示同时满足 
      match:  
        - {paramter: name ,value: jone} 



```
第二步，完成data数据文件编写，同样是yaml格式。在前面我们定义了两个数据success和error。  
我们在datas目录下新见first这个目录，并在first目录下新增 success.yaml文件和error.yaml文件。
文件内容如下：  
sucess.yaml    
```yaml
# 数据编号，在接口内唯一就可以
code: success
# 描述
desc: 成功的返回结果
# 返回的数据的格式
style: json
# 数据体，支持变量方式 ，例如用{{.name}} 表示引用参数name的值
data: |
  {"success":true,"errcode":"0000","errmessage":"成功"}
``` 
error.yaml  
```yaml
code: error
desc: 错误的返回结果
style: json
data: |
  {"success":false,"errcode":"10001","errmessage":"不允许通过，请配合！"}
```
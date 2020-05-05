# 提高篇
## 一个API服务的构成
我们访问一个服务端提供的接口一般形式为：[地址]:[端口]\[接口地址]  
然后我们传递相应的参数，可能是 name=x1&name2=x2形式，也可能是json、xml格式，当然还有其它形式，
只是这几种形式比较常用。  
服务端会校验参数，参数校验通过后会经过一定的逻辑运算等处理，最后饭后处理的结果给请求者，可能是json格式或xml格式。  
## 端口设置
第一步 myway-mock 首先支持端口的设置，达到和目标服务一致的效果。
该配置在configs.yaml中的port中，port设置的形式是数组，按yaml的标准如下：
```yaml
# 这表示设置了多个端口
port:
  - 8080
  - 9090
------------------------
# 当然也可以只用一个
port:
  - 8080
```
什么情况下会用到多个端口呢？其实也可以通过部署多套myway-mock来达到这个效果，部署多套需要同步多套api的定义，
这样就需要保证api的定义和用例不要出现人为的偏差。而且很多情况下mock服务的定义都是一样的，如果你同作者一样需要一套mock来支持多套环境的话，这个特性是能帮你省下不少事情的。
#### 假如不设置端口会怎么样
如果port没有定义或没有设置，myway-mock会使用默认的80端口来工作。  
因此强烈建议指定端口
#### 端口冲突会如何
如果你定义的端口，已经被其它应用服务占用了，那么myway-mock将会启动失败，会直接退出。  
所以当你发现myway-mock启动失败的时候，很大可能是端口被占用了。
## 接口的定义
接口与接口之间区分，除了接口名称，然后就是接口的访问的url了。而我们访问一个http接口时候，是通过不同的http方法访问的。
因此最基本的描述就是：接口名、接口访问url地址、接口允许的http方法。  
其中http方法常用的就是GET和POST了。
```yaml
# 接口名
name: hello
# 接口访问的url地址
url: /helloworld
# 描述，用于阐述接口的功能
description: hello world
# 接口支持的http方法。取值都是大写的，允许的有 GET、POST、PUT、DELETE
methods:
  - GET
  - POST
```
#### 服务请求
服务请求描述的是我们访问http接口提供的请求参数，包括了请求参数的形式和参数本身。常用的有
跟随在地址后的参数列表,紧跟着一个？号。例如：localhost\helloworld?name=test&sex=1。  
或者是使用post方法提交的表单格式，以及现在比较流行的json格式，当然还有xml格式的请求。  
而这些都在request下定义，style 表示的是请求的参数形式，取值有json、xml、url-form；  
既然有参数了，那一定会发生参数不和约定的情况，可以通过error来定义在参数校验不通过的时候，接口返回什么样的格式的报文。
其中有两个预定义的变量可以使用{{ ._Parameter}} 校验不通过的参数名，{{._Msg}} 对应的错误信息。  
然后通过items定义需要关心的参数，每个参数的属性如下：  
   * paramter 参数名
   * type 参数类型
   * length 参数最大长度
   * policy 策略，支持两种 1、Must 必填 2、Option 可选
   * validator 校验器。现在支持的有regex 正则表达式， dict 字典
   * expr 校验表达式，这些是提供给校验器使用的。一般用于存放正则表达式（如果使用regex校验器的话),以及 字典名称
```yaml
# 请求定义
request:
  # 请求参数的形式，json格式
  style: json
  # 请求参数有误的时候返回的报文模板
  error: |
    {"success":false,"errcode":"50001","errmessage":"参数{{ ._Parameter}},校验错误 {{._Msg}}"}
  # 请求参数定义，可以多个
  items:
    - {paramter: name ,type: String ,length: 30,policy: Must,validator: }
    - {paramter: world ,type: Sex ,length: 30,policy: Option,validator: }
```
#### 服务返回报文
服务访问报文通过response标签来定义的，由两部分构成。default 默认的返回值(报文编号)，triggers 触发定义。  
其中trigger定义由两部分构成，1、报文编号 2、触发的匹配条件，可以多个，条件之间的关系是取逻辑并，也就是说多个条件同时匹配才满足触发条件。  
trigger的意义就是当条件同时满足时候，触发定义的接口返回指定的报文。那么当没有一个trgger被触发的时候，就返回default指明的报文。  
示例如下：
```yaml
response:
  #默认返回结果编号
  default: success
  #触发条件，根据条件返回不同的数据
  triggers:
    # 报文编号
    - data: success
    # 匹配的条件定义，参数 name 取值为mike
      match:
        - {paramter: name ,value: mike }
    # 这个trigger定义的是当参数name 取值为jone，返回error对应的报文。
    - data: error
      match:
        - {paramter: name ,value: jone}

```
通过定义多个trigger，被模拟的api就可以根据请求参数的不同取值，返回不同的报文，来模拟不同的场景。

#### 报文定义
为了防止报文过多，造成api文件臃肿。报文的定义在单独的文件中，并放在datas目录下，并放在api同名的目录下，这样不同的api报文的编号可以相同，不必为了命名而伤透脑筋了。  
如上例的hello接口，我们需要在datas下建立hello目录，并把报文放在这个目录下。接上面的例子，我们定义了success和error两个报文。  
报文的组成由四部分构成
* code 报文编号，与文件同名
* desc 报文描述，用于说明报文含义的，防止忘却了
* style 报文格式说明。取值：json、xml 
* data 报文数据，支持对参数的引用，引用时候使用 {{.参数名}} 方式来引用  
例如 success.yaml
```yaml
code: success
desc: 成功的返回结果
style: json
data: |
  {"success":true,"errcode":"0000","errmessage":"成功"}
```
error.yaml
```yaml
code: error
desc: 错误的返回结果
style: json
data: |
  {"success":false,"errcode":"10001","errmessage":"入口被关闭"}
```

#### 如何模拟耗时
有时候在特定场景下，我们需要模拟的接口，不要太快返回，起码要显得真实点，于是就需要模拟耗时。比如在做性能压测的时候，有些服务不能被压，但我们有一些耗时数据。这时候我们可以让myway-mock模拟该api的时候，表现的在努力工作的样子，耗费些时间。  
通过在api的定义文件中设置delay属性可以来完成这个诉求。
这个属性支持的值为数组，每个成员的取值必须是整型，单位是毫秒。myway-mock会随机在delay定义的数组中取出一个值来进行延时处理。
例如：
```yaml
name: hello
url: /helloworld
description: hello world
methods:
  - GET
  - POST
# 模拟耗时，单位毫秒，以下定义了多个耗时
delay:
  - 100
  - 150
  - 230
  - 400
```
如果不知道是在使用mock的api，是不是会以为在访问真实的服务接口。

## 预定义的类型
* String 字符串
* Date 日期类型，允许的格式 2020-05-01 
* Time 时间类型，允许的格式 19:30:23
* DateTime 日期时间类型，允许的格式 2020-05-01 19:30:23
* Dict 字典类型，表示取值必须是字典中允许的值
* Int 整型，整数类型
* Numric 数字型字符串，字符中只允许出现数字
* Bool 布尔类型
  - {name: String,label: 字符串,length: ,option: ,validator: }
  - {name: Date,label: 日期 ,length: 10,option: "[1-9]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])" ,validator: regex}
  - {name: Time,label: 时间 ,length: 8,option: "(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d" ,validator: regex}
  - {name: DateTime,label: 日期时间 ,length: 19,option: "[1-9]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\\s+(20|21|22|23|[0-1]\\d):[0-5]\\d:[0-5]\\d" ,validator: regex}
  - {name: Dict,label: 字典类型,length: ,option: ,validator: dict }
  - {name: Numric,label: 数字字符串,length: ,option: "[0-9]\\d*" ,validator: regex }
  - {name: Sex,label: 字典类型,length: ,option: sex ,validator: dict }
  - {name: Int,label: 整型 ,length: 10, option: "[-]?[0-9]\\d+", validator: regex}
  - {name: Bool,label: 布尔类型,length:5 ,option: bool ,validator: dict }
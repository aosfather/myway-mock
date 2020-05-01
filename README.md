# myway-mock
# 吾道-模拟API服务
”吾道“之模拟服务，该应用能实现模拟基于http协议的web服务，可用于模拟外部接口和内部未实现的接口服务。适用于前后端分离以及服务依赖之间的解耦。
## version 1.0.0
 基于“吾道”(myway)的技术，实现对api服务的模拟。  
 “吾道” 三件套之一，三件套包括：myway-gateway，myway-mock、myway-datamanager。
 myway-mock 是一个简单的高效基于服务端的mock解决方案，对于现有的代码逻辑不用因为服务是模拟的而需要特殊处理。  
 在实际开发中不会存在所谓的测试代码和正式代码，避免由于频繁的修改代码造成的质量问题。  
  我相信myway-mock会成为实际开发中不可多得的工具，
 在实际很多场景中都会派得上用场。  
 例如：
 * 某个服务由于多种原因只有生产环境，没有开发及测试环境
 * 某个场景重现很复杂或不好重现的时候
 * 依赖的某个服务由于资源等等原因，无法提供联调环境的时候
 * 前后端分离，前端在内部做页面应用的时候，想要比较真实的交互调试的时候
 * 搭建独立的稳定的单元验证调试环境，避免互相发版影响
 * 性能压测
 * 等等
 
### 特性列表
 * 支持http请求方式
    * GET请求
    * POST请求
 * 支持参数形式
   * url-form
   * json
 * 支持参数校验
   * 必填
   * 可选
 * 内置常用的类型
   * 日期 YYYY-MM-DD(年月日)
   * 日期时间 YYYY-MM-DD hh:mm:ss(年月日时分秒)
   * 整形
   * 数字字符串格式
   * 手机号
   * 
 * 支持根据输入参数值条件返回不同的数据
 * 返回的数据支持变量，可以引用请求参数和内置参数
 * 支持指定端口启动
 * 支持热加载api服务定义，更新和新增定义不用重启
 * 支持yaml格式的api服务定义
 ### api服务定义模型说明
 #### 组成
 一个需要被模拟的api由如下几部分构成
 * url api访问的地址描述
 * api说明
 * 支持的访问方式(GET、POST、PUT、DEL)
 * 参数形式(url-form、json、xml)
 * 参数定义列表
   * 参数名称
   * 参数类型
   * 是否必填
   * 额外的校验规则，例如正则表达式校验
   * 参数的长度要求
  * 默认的返回数据
  * 条件匹配定义(match)
    * 触发条件
    * 返回数据ID 
 #####  额外的数据定义
 * 数据ID
 * 数据类型(text、json、html、xml)
 * 内容
 #### 特别说明 
 出于对性能和资源的占用的考虑，额外的数据和api的定义单独设计。
 api的定义存放于apis的目录下，数据存在datas目录下。  
 其中在datas目录使用api的名称作为子目录名，每个数据使用数据ID的名称命名。
 ##### 举例如下
 假设定义一个api叫“xxxapi”，其中定义了两条结果数据，一个叫 data1，一个叫data2
 ，那么数据文件和定义文件的存放位置如下所示。
 * myway-mock所在目录
   * apis
     * xxxapi.yaml
   * datas
     * xxxapi
       * data1.yaml
       * data2.yaml 
 #### api定义的yaml样例模板
 详情见 example目录。其中 api.yaml为api的定义样例，data.yaml为返回数据的样例 
 ## 特性计划
 * 支持多种http请求方式
   * put
   * delete
 * 支持 xml 格式的请求
 * 支持模拟耗时
 * 支持加解密
   * 非对称
   * 对称
        
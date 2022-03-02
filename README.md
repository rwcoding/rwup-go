## 基于 toml 描述的文档系统同步客户端

> 下划线开头的文件或文件夹会被忽略  
> 请使用英文目录名和文件名  

### 配置文件
+ 默认文件名 `.docg`
+ 可在命令中指定 docg -conf=/etc/.docg
+ 默认从当前命令执行目录中查找
+ 格式
```ini
# 配置示例

# 文档系统分配的账号密码
username = admin
password = admin

# 工程标识
project = pms

# 根目录，可指定绝对目录
root = docs 

# 文档系统接口地址
url = http://localhost:80

# 解析的文件扩展名
parse = md,toml
```

### TOML配置编写，参见`sample`示例文件夹
> `_dirs.toml`  
> 为您的目录定义中文名称，否则默认使用目录名
> 
> `_titles.toml`  
> 为某些文件定义文档标题，如 markdown 文件，否则使用目录名

### TOML接口描述示例
```toml

title = "配置列表" # 文档标题
route = "/api/config/list" # 路由或唯一指令

doc_desc      = "" # 文档事项描述
request_desc  = "" # 请求事项描述
response_desc = "" # 响应事项描述

[request] # 请求参数列表
page      = ["int", "页码", "required|min=1", "如: 1"]
page_size = ["int", "每页数量"]
key       = ["string", "搜索关键词"]

# 数组第1项：数据类型
# 数组第2项：描述
# 数组第3项：规则，竖线分割，完整如：required|ignore|min=1|max=10|order=100
#          required 参数必需
#          ignore 参数不必需
#          min 最小，字符串代表长度，数字代表大小
#          max 最大，字符串代表长度，数字代表大小
#          order 排序，所有的字段排序默认100，越小越靠前
# 数组第4项：参数样例

[response] # 响应参数列表
datas   = ["array[item]", "数组数据"]
count   = ["int", "总数"]
user    = ["user", "对象数据示例"]

[component.item] # 组件,混入datas array[item]
id         = ["int", "配置ID"]
name       = ["string", "名称"]
k          = ["string", "键"]
v          = ["string", "值"]
data_type  = ["string", "数据类型"]
created_at = ["string", "创建时间"]

[component.user] # 组件,混入对象 user
username = ["string", "用户"]
name     = ["string", "名字"]

```
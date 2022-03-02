## 请求-方法

```HTTP POST```

## 请求-路由 (Route) 

``` 在 URL 参数中 ```

## 请求-HTTP头

| 参数           | 描述                  |
|--------------|---------------------|
| X-Time       | UNIX时间戳             |
| X-Sign       | 签名                  |
| X-Token      | 用户标识，登录后获取，登录前为空字符串 |
| Content-Type | application/json    |

## 请求-数据

```JSON``` 格式的字符串

## 请求-签名方法

> md5(Route + X-Time + X-Token + Body + TokenKey)  
> 不包含 '+'  
> 结果为小写

## 响应数据

```JSON``` 格式的字符串

```json
{
    "code": 10000,
    "msg": "",
    "data": {}
}
```

- `code` 错误码
    - `10000` 请求成功
    - `99999` 请求失败
    - `11000` 需要登录，遇到该错误码，引导用户前去登录
- `msg` 请求结果描述，通常用作错误信息提示
- `data` 响应业务数据，并不是每个接口都有该字段，不需要业务数据的接口没有该字段

## 示例

```javascript
const token = '4e66ccf7c8c8dd1491fffcdaca1cf88f' // 登录后获取，登录前为随机32字符串
const key = '47a004d84228e9fa599cd945b4ace0b4' // 登录后获取，登录前为随机32字符串

const route = '/api/open/login'
const time = Math.round(new Date().getTime() / 1000)
const body = JSON.stringify({username: 'admin', password: 'admin'})
const sign = md5(route + time + token + body + key)
```

```text
POST /api/open/login

Context-Type: application/json
X-Time: 1646113502 
X-Token: 4e66ccf7c8c8dd1491fffcdaca1cf88f
X-Sign: 7cb8415906f3083328573f77458938d1 

{"username":"admin",password:"admin"}
```






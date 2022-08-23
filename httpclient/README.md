## 简述

方便使用的`http`库

## 功能

- 链式调用
- 超时设置
- 代理
- 缓存

## 使用

### 基本使用

```go
client := zerohttpclient.NewClient()

ctx := client.Get("https://www.keylala.cn")
result, _ := ctx.ToString()
fmt.Println(result)
```

### 添加参数

```go
client := zerohttpclient.NewClient()

ctx := client.WithParams(map[string]interface{}{
    "id": "123456",
    "name": "zero",
}).Get("https://www.keylala.cn")


ctx = client.WithBody(map[string]interface{}{
    "id": "123456",
    "name": "zero",
}).Post("https://www.keylala.cn")


// 指定 Content-Type 格式，默认 application/x-www-form-urlencoded
ctx = client.WithBody(map[string]interface{}{
    "id": "123456",
    "name": "zero",
}).WithContextTypeJSON().Post("https://www.keylala.cn")
```

### 设置超时

默认连接超时时间 2 秒，每次调用超时时间 5 秒

```go
client := zerohttpclient.NewClient().WithDialTimeout(time.Second * time.Duration(1)).WithTimeout(time.Second * time.Duration(2))
ctx := client.Get("https://www.keylala.cn")
```

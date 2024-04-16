## 优雅的重启和关闭服务器

- 优雅的重启需要获取原进程监听套接字的文件描述符，所以不能直接使用`http.server`的`ListenAndServe`和`ListenAndServeTLS`
- 只要将`http.Server`替换为`github.com/zerogo-hub/zero-helper/graceful`即可
- 优雅的关闭服务器:

  - 服务器会关闭监听，执行用户注册的清理函数，等待未处理完毕的连接继续处理
  - 服务器每隔 `500` 毫秒判断一次所有连接是否处理完毕，如果处理完毕，则关闭服务器
  - 如果直到超时，都未处理完，直接关闭服务器。如果没有额外设置，默认超时时间为 `5` 秒

- 优雅的重启服务器:
  - 新进程"继承"旧进程的 `os.Stdin, os.Stdout, os.Stderr` 三个文件描述符以及 `监听套接字` 的文件描述符
  - 旧进程优雅的关闭，监听新连接的工作由新进程接手

## 消息处理

- ctrl + c: 关闭
- kill: 关闭
- ctrl + \: 重启

## 使用测试示例

- 在测试目录下执行`go run main.go`，生成`main`文件
- 运行`./main`运行程序
- 根据终端输出打开网页
- 修改`main.go`代码，并执行`go run main.go`，生成新的`main`文件
- 网页中有提示形如`kill -USR2 70386 to restart`的提示，在终端运行
- 刷新网页，可以发现修改部分已经出现在网页上了

## 使用默认信号示例

见 [examples/default/main.go](./examples/default/main.go)

## 用户自定义清理函数示例

见 [examples/userdefined/main.go](./examples/userdefined/main.go)

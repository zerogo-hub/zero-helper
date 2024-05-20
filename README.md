# zero-helper

# 模块

- actor: Actor 模型
- array: 切片操作，如切片转字符串，`[1,2,3] -> 1+2+3`
- bloom: 封装布隆过滤器
- buffer
  - circle: 环形缓存区
- bytes: `[]byte`相关
- cache: 封装`redis`
- codec: 编码与解码器
- compress: 压缩与解压
- config: 读取配置表
- crypto: 加密与解密
- database: 封装`mysql`
- email: 发送邮件
- entity: `cache-aside`，封装`gorm`和`bigcache`
- file: 文件相关
- graceful: 优雅的重启和关闭服务器
- human: 身份证验证
- ip: IP 字符串类型地址与整型的转换
- json: 替换默认的 `json` 解析库
- jwt: 封装 `jwt`
- locker: 分布式锁
- logger: 日志相关
- os: 系统相关
- random: 随机数
  - choice: 权重随机
  - shortsf: 46 位，workID [1,8]
  - snowflake: 64 位，workID [0,1023]
  - uuid: 单机版
- reflect: 封装 `reflect`
- regexp: 一些正则表达式
- time: 时间相关
- timer: 定时器，时间轮
- words: 字数统计

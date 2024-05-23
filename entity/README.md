# 简述

适合读多写少的场景

# 特点

- 多级数据：内存缓存(big-cache)，远端缓存(redis)，数据库(mysql)，自定义获取(未实现)
- 数据一致性：所有内存缓存通过`pub-sub`保障一致性，更新内存缓存时也同时更新远端缓存

# 策略

可以自由组合数据获取策略

- 内存缓存 + 远端缓存: 本地查看排行榜，本地短暂缓存排行榜数据
- 内存缓存 + 数据库: 用于本地直接操作数据库时，缓存中无数据时从数据库获取
- 内存缓存 + 自定义：用于专用数据库服务器操作数据时，缓存中无数据时使用自定义函数从数据库服务器拉取数据

# 目录

- cache: 缓存实现
  - bigcache: https://github.com/allegro/bigcache
  - freecache: https://github.com/coocood/freecache.git
  - redis: github.com/gomodule/redigo
- db: 数据库实现
  - gorm: gorm.io/gor
- example: 示例
- entity: 实体
- result: 查询结果封装

# 替换

## 替换内存缓存

## 替换数据库

## 自定义获取

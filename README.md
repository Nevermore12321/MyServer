# MyServer
Keep going. Keep knowing.

# 框架结构
- Gin
    - go get -u github.com/gin-gonic/gin
- Gin-Swagger
    - go get -u github.com/swaggo/swag/cmd/swag
    - go get -u github.com/swaggo/gin-swagger
    - 每次修改了 接口， 要 swag init 
- viper
    - go get -u github.com/spf13/viper

- logger： zap + Lumberjack（归档，写入文件）
    - go get -u go.uber.org/zap
    - go get -u github.com/natefinch/lumberjack
    
- gorm : 操作 mysql
    - go get -u github.com/jinzhu/gorm

## gin-swagger 使用

1. swag init
2. go run main.go
3. http://IP:Port/swagger/index.html

## Gin 常用

### 参数，表单的获取
查询参数 | Form 表单 | 说明  
--- | --- | ---
Query | PostForm | 获取key对应的值，不存在为空字符串
GetQuery | GetPostForm | 多返回一个key是否存在的结果
QueryArray | PostFormArray | 获取key对应的数组，不存在返回一个空数组
GetQueryArray | GetPostFormArray | 多返回一个key是否存在的结果
QueryMap | PostFormMap | 获取key对应的map，不存在返回空map
GetQueryMap | GetPostFormMap | 多返回一个key是否存在的结果
DefaultQuery | DefaultPostForm | key不存在的话，可以指定返回的默认值




https://github.com/ReadRou/gin_project
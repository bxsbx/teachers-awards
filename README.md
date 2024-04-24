使用swag命令时先安装swag命令工具

go get github.com/swaggo/swag/cmd/swag

本项目采用三层架构，controller、service、dao三层，
middleware存放gin的中间件，
global存放全局变量、常量、方法等，
client放第三方接口调用，
common为公共包，主要存放第三方库的接入、自定义error、自定义util等，
generate为代码模板生成器
model存放输入输出数据结构模型
router为方法路由
migration存放sql迁移文件

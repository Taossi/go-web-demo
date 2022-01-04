## **golang-web demo**

一个golang实现的web后端demo，比较基础，用于自己熟悉入门golang web开发。


主要使用了MVC设计模式，实现了简单的用户注册、登陆功能，并且利用jwt中间件来保护路由。 用户数据存放在MySQL中。


所用框架：gin gorm viper

gin: web框架

gorm：golang-orm，用于连接MySQL数据库处理

viper：用于配置文件，本项目中viper使用yml文件配置数据库。

jwt: JSON Web Token,用于检验token保护路由，详细可见https://www.liwenzhou.com/posts/Go/jwt_in_gin/

目前没有实现前端部分，需要使用postman测试API功能，主要测试内容为：注册、登陆功能；token中间件是否有保护作用等。

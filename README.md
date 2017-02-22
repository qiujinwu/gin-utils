1. gin有多种路径，各个插件不统一
1. session使用了gorilla的context模块，性能较差
1. gin内置的模板连标准的go模板都支持不全，引入pongo2，语法和Django类似，学习简单，功能强大
1. 无简单合适的认证模块，自己写
1. flash和crsf没必要和session一样存到后台session数据库中
1. 数据校验模块会自动设置http code，傻逼做法

# Comic2Video (comic to video)
A simple service that can download manga and supply by video

## TODO

- [x] 18comic适配
- [x] 18comic服务适配
- [x] URL 智能识别
- [x] 任务队列和等待系统
- [ ] 合并转换函数调用
- [ ] 通过数据库和哈希池管理视频文件？
- [ ] 服务读取环境变量设置
- [ ] 失败任务冷却


## 架构
ADAPTER 爬取 + 下载 （下载和爬取受到平台影响，难以抽象）-> 下载完成的资源列表
CONVERTER 转码成视频，入文件池

ADAPTER 用结构体+接口的方式管理

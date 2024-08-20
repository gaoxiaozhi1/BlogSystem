## 综合博客系统

项目描述：该项目是一款基于gin+es实现的多功能博客系统，主要用于用户（注册，登录，实现了多种登录
方式），文章（发布，检索，评论，点赞，发布日历），匿名群聊等。

技术栈： Go(gin,gorm), DB(MySQL,Redis), MQ(Elasticsearch)

缓存与存储：使用Redis实现高频信息缓存和七牛云进行对象存储。

用户认证：使用JWT保存登录信息，支持邮箱登录和绑定，实现验证码发送。

搜索功能：使用Elasticsearch实现全站内容搜索，支持关键词搜索，标题搜索等。

用户交互：实现文章搜索，支持根据排序搜索，标签搜索，标题搜索等，使用WebSocket实现聊天室功
能
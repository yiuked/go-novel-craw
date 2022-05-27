# go-novel

## 概述

`go-novel`是一款通过`go`实现的多协程小说采集器，通过指定特定的`yaml`文件来采集不同的网站， 、`go-novel`支持大部分小说网站采集， 只需要根据根据目标网站的格式调整`yaml`文件中的采集规则则可。

`go-novel`采用分段采集的方式，将采集过程分为：书籍ID采集 =》 书籍详情采集 =》 图片采集 =》 章节列表采集 =》 章节内容采集，来实现任务间的解藕。同时通过`sqlite`来实现单个阶段的渐进式采集。
高效的避免采集过程出现数据异常对整体采集进度的影响性。

`go-novel` 不仅支持小说的全量采集，还支持小说的更新采集，可以灵活的指定区间段的内容进行定向更新。

## 安装

```
go install github.com/yiuked/go-novel
```

## 命令介绍

```
$go-novel help
NAME:
   go-novel - A new cli application

USAGE:
   src.exe [global options] command [command options] [arguments...]

COMMANDS:
   api      start web api service 启动一个 web api 服务,用于展示已采集的数据
   bookid   craw books ID to local 采集书籍ID到本地
   chapter  craw books chapter to local 采集章节内容到本地
   content  craw books content to local,save to HTML file 采集书籍详情并另存了HTML文件
   cover    craw books cover image to local 采集图片并保存到本地
   detail   craw books detail to local 采集书籍详情介绍
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --data_dir value, -c value   set（sqlite、cover、HTML）storage path,default current 数据保存地址
   --goroutine value, -g value  set goroutine craw limited (default: 20) 同时工作的协程总数
   --source value, -s value     set need craw source yaml file (default: "./yaml/tudu.yaml") 指定采集源文件
   --help, -h                   show help (default: false)

```

## 采集示例

1. 采集书籍ID

```
go-novel bookid
```

通过添加`--page`参数,可以配置采集多少页,此处的页数为每个分类采集多少页

2. 采集书籍内容

```
go-novel detail
```
如果添加`--update`参数,则从已采集详情中的书籍进行更新

3. 采集书籍封面

```
go-novel cover
```

4. 采集章节

```
go-novel chapter
```

5. 采集内容

```
go-novel content
```
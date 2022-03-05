# 聊天室

## 使用：

#### 服务端:

>启动:

在/chat_server/src目录下执行

`go run main.go`

或者windows环境下双击 *main.exe* 运行

>GM指令：

直接在服务端终端执行指令：

`/stats [username]` 打印玩家的信息

`/popular [roomid]` 打印最近10分钟发送频率最高的词

#### 客户端:
>启动:

在/chat_client/src目录下执行

`go run main.go`

或者windows环境下双击 *main.exe* 运行

>操作：

直接在客户端终端执行指令：

`Login:[username]`  登录用户

`JoinRoom:[roomid]`  加入房间

`SendMessage:[message]` 发送消息

## 算法：

#### 脏词过滤

算法：前缀树匹配，具体实现参考：

/chat_server/src/common/trie.go

>性能指标：

以提供的脏词库，16核E5-2680,61个字符为基准，基准测试为15000ns/op

>优化建议：
    
如果是纯*号替换，直接匹配成功可直接跳过已匹配字符。
考虑到扩展，如果出现以其他字符替换，如：'s',存在替换后满足脏词的情况故未实现。

#### 获取最近10分钟频率最高的词
    
算法：滑动窗口，在房间里面维护一个消息记录的左指针`Left`，在发送消息/请求热词的时候移动到最新的10分钟时间点前，
再用一个`HotRecord map[string]int64`用来保存，左指针到当前最后一条消息每个词出现的频率。

## 三方库

`github.com/gogo/protobuf` 用于序列号请求参数
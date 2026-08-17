# Bug 是什么

多个请求同时导入和读取周排班时，缓存、内存存储和统计计数发生数据竞争。

# 如何触发

用 32 个 goroutine 并发导入排班并读取同一周，然后用 `go test -race` 运行并发用例。

# 错误信息

```text
WARNING: DATA RACE
fatal error: concurrent map writes
```

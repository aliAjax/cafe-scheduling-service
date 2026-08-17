# Bug 是什么

周视图只返回一个日期，排班、备注和标签还会跟着调用方或缓存读写一起变化。

# 如何触发

在 `internal/roster` 中构造两个排班，调用周视图构建，再修改输入切片和缓存读取结果。

# 错误信息

```text
unexpected week shape: ... Days:[]roster.Day{...}
view aliases input
```

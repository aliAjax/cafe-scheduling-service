# Bug 是什么

没有备注的排班在复制、周视图、缓存或 JSON 导出路径上会触发空指针崩溃。

# 如何触发

构造 `Note=nil`、`Tags=nil` 的排班，依次调用复制、周视图、缓存和导出测试。

# 错误信息

```text
panic: runtime error: invalid memory address or nil pointer dereference
```

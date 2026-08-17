# Bug 是什么

调用方已经取消后，排班导入、周视图构建和 JSON 导出仍继续执行并产生结果。

# 如何触发

先取消一个 context，再把它传给导入、周视图、缓存读取和导出方法。

# 错误信息

```text
import error = <nil>, want context canceled
```

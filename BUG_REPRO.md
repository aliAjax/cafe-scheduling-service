# Bug 是什么

排班校验错误经过周视图和导入服务后丢失了可判断的错误类型与字段信息。

# 如何触发

提交 `employeeId=0` 的排班，分别调用周视图构建和导入服务，再用 `errors.Is`、`errors.As` 检查返回值。

# 错误信息

```text
errors.Is failed for assignment 0: employeeId: must be positive
```

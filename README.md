# Cafe Scheduling API

面向咖啡店的小团队排班系统 Go 后端。店长维护员工、班次类型和每周营业时间，在周视图为员工安排班次；系统检测同一员工的重叠班次，并统计每人每周排班小时数。员工可查看自己的排班。

`internal/roster` 提供并发安全的周视图规划组件，负责排班导入、按周分组、缓存隔离、取消传播和 JSON 导出。组件不依赖数据库，可用于后台批量导入和后续调班审批。

## 技术栈

- Go 1.22
- MySQL 8.4
- `database/sql` + `github.com/go-sql-driver/mysql`
- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`
- Docker Compose

## 启动

```bash
cp .env.example .env
# 按需修改 .env
docker compose up -d --build
```

服务监听容器内 `8080`，宿主端口 `18092`。健康检查地址：

```bash
curl http://127.0.0.1:18092/healthz
```

## 演示账号

| 角色 | 用户名 | 密码 |
| --- | --- | --- |
| 店长 | `manager` | `manager123` |
| 员工 | `alice` | `employee123` |
| 员工 | `bob` | `employee123` |

首次启动会自动执行 `migrations/*.sql`，并创建以上演示账号、员工、班次类型和营业时间。

## 接口概览

所有业务接口均使用 `Authorization: Bearer <JWT>`。除登录和健康检查外，其余接口需要认证。

### 认证

- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`

### 员工管理（店长）

- `GET /api/v1/employees`
- `POST /api/v1/employees`
- `GET /api/v1/employees/{id}`
- `PUT /api/v1/employees/{id}`
- `DELETE /api/v1/employees/{id}`

### 班次类型（店长）

- `GET /api/v1/shift-types`
- `POST /api/v1/shift-types`
- `GET /api/v1/shift-types/{id}`
- `PUT /api/v1/shift-types/{id}`
- `DELETE /api/v1/shift-types/{id}`

### 营业时间（店长）

- `GET /api/v1/business-hours`
- `PUT /api/v1/business-hours`

### 排班

- `GET /api/v1/schedules?week_start=YYYY-MM-DD`（员工只能看到自己）
- `GET /api/v1/schedules/me?week_start=YYYY-MM-DD`
- `POST /api/v1/schedules`
- `PUT /api/v1/schedules/{id}`
- `DELETE /api/v1/schedules/{id}`

### 统计（店长）

- `GET /api/v1/statistics/weekly?week_start=YYYY-MM-DD`

## 示例

```bash
TOKEN=$(curl -s -X POST http://127.0.0.1:18092/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"manager","password":"manager123"}' | jq -r .token)

curl -s http://127.0.0.1:18092/api/v1/employees \
  -H "Authorization: Bearer $TOKEN"
```

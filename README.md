# API密钥管理系统使用文档

## 系统概述

这是一个基于Go语言实现的API密钥管理系统，提供API密钥的创建、验证、鉴权和有效期管理功能。系统使用SQLite数据库存储API密钥信息，通过Gin框架提供RESTful API接口。

## 功能特点

- API密钥的创建、读取、更新和删除
- API密钥有效期管理
- API密钥状态管理（激活/禁用）
- 支持通过查询参数方式使用API密钥（便于日常使用）

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/wuhenai/api_tool.git
cd api_tool

# 安装依赖
go mod tidy

# 编译
go build -o apikey-server ./cmd/server
```

### 运行

```bash
./apikey-server
```

服务将在 `http://localhost:8080` 启动。

## API接口说明

所有API请求都需要包含有效的API密钥作为查询参数：`?key=YOUR_API_KEY`

### 初始化API密钥

首次使用系统时，需要创建一个初始API密钥：

```bash
curl -X POST http://localhost:8080/api/init-key \
  -H "Content-Type: application/json" \
  -d '{"name":"初始密钥", "user_id":1, "expires_in_days":365}'
```

### 创建新的API密钥

```bash
curl -L -X POST "http://localhost:8080/api/keys?key=YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"新密钥", "expires_in_days":30}'
```

### 获取所有API密钥

```bash
curl -L -X GET "http://localhost:8080/api/keys?key=YOUR_API_KEY"
```

### 获取特定API密钥

```bash
curl -L -X GET "http://localhost:8080/api/keys/1?key=YOUR_API_KEY"
```

### 更新API密钥

```bash
curl -L -X PUT "http://localhost:8080/api/keys/1?key=YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"更新的密钥名称", "active":true, "expires_in_days":60}'
```

### 删除API密钥

```bash
curl -L -X DELETE "http://localhost:8080/api/keys/1?key=YOUR_API_KEY"
```

### 健康检查

```bash
curl -X GET http://localhost:8080/health
```

## API响应格式

### 成功响应

```json
// 获取API密钥列表
[
  {
    "id": 1,
    "key": "b520b25d-f1c9-4e7b-a4e4-7b2318acd736",
    "name": "初始密钥",
    "user_id": 1,
    "expires_at": "2026-06-10T12:22:09.114674+08:00",
    "active": true,
    "created_at": "2025-06-10T12:22:09.116954+08:00",
    "updated_at": "2025-06-10T12:22:09.116954+08:00"
  }
]

// 删除API密钥
{
  "message": "API密钥已成功删除"
}
```

### 错误响应

```json
{
  "error": "API密钥不存在或不属于当前用户"
}
```

## 数据存储

系统使用SQLite数据库存储API密钥信息，数据文件位于 `data/apikeys.db`。

## 环境变量

- `PORT`: 服务器端口号，默认为8080

## 注意事项

1. 初始API密钥只能创建一次，如果已存在则无法再次创建
2. API密钥会在过期时间后自动失效
3. 使用查询参数 `?key=YOUR_API_KEY` 方式验证API密钥，便于日常使用

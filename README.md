# WeAvatar

## 说明

这是 WeAvatar 的后端项目，基于 Goravel 框架开发。

## 依赖

- Nginx >= 1.18
- Go >= 1.20
- PostgreSQL >= 15
- libwebp-dev（用于处理WEBP图像）

## 部署

### 1. 安装依赖

```bash
sudo yum install libwebp-devel
```

### 2. 安装环境

自己解决

### 3. 配置 Hosts

需要配置 `hosts` 文件，将 `proxy.server` 指向反向代理服务器的 IP 地址。

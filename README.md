# WeAvatar

## 说明

这是 WeAvatar 的后端项目，使用 AGPLv3 协议开源。

WeAvatar 是超越 Gravatar 的新一代头像服务，不仅支持用户上传头像，同时也能从 Gravatar、QQ 上获取头像，同时支持 AI 自动化审核。

## 依赖

- AlmaLinux / RockyLinux 9
- OpenResty >= 1.21
- Go >= 1.21
- PostgreSQL >= 16
- vips >= 8.15

## 部署

### 1. 安装依赖

先导入 [Remi](https://blog.remirepo.net/pages/Config-en) 源和 [RPM Fusion](https://rpmfusion.org/Configuration) 源再运行命令。

```bash
dnf install -y vips vips-heif vips-tools
dnf install -y libheif-freeworld
vipsthumbnail -v
```

确保 vipsthumbnail 命令能正确输出版本号，否则 WeAvatar 主程序无法生成头像。

### 2. 安装环境

在耗子 Linux 面板安装对应的环境。

### 3. 配置 Hosts

需要配置 `hosts` 文件，将 `proxy.server` 指向反向代理服务器的 IP 地址。

规则请参考 [https://git.haozi.net/opensource/mirror-conf](https://git.haozi.net/opensource/mirror-conf)。

### 4. 配置 .env

```bash
cp .env.example .env
```

自行修改 `.env` 文件中的配置。

### 5. 初始化数据库

```bash
./weavatar artisan migrate
./weavatar artisan hash:make
./weavatar artisan hash:insert
```

其中 `hash:make` 步骤可选，用于生成 QQ 邮箱的 Hash 表，纯 csv 约占用 150 GB，导入后约占用 450 GB。

### 6. 启动项目

```bash
./weavatar
```

建议通过 `supervisor` 管理进程。

### 7. 配置 Nginx

伪静态规则如下

```nginx
set_real_ip_from 0.0.0.0/0;
real_ip_header X-Forwarded-For;

location / {
    try_files $uri $uri/ /index.html;
}

location /api/
{
    proxy_pass http://127.0.0.1:3002/;
    proxy_set_header Host weavatar.com;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header REMOTE-HOST $remote_addr;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Accept-Encoding "";
    proxy_cache off;
    proxy_no_cache 1;
    proxy_cache_bypass 1;
}

location /avatar {
    rewrite ^/avatar(.*)$ /api/avatar$1 last;
}
```

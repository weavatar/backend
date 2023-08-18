# WeAvatar

## 说明

这是 WeAvatar 的后端项目，基于 Goravel 框架开发。

## 依赖

- Debian 12
- OpenResty >= 1.21
- Go >= 1.21
- TiDB >= 7.1
- libvips-dev >= 8.10

## 部署

### 1. 安装依赖

```bash
apt install -y libvips-dev
```

### 2. 安装环境

部署 TiDB 集群，参考 [TiDB 部署文档](https://docs.pingcap.com/zh/tidb/stable/production-deployment-using-tiup)。

### 3. 配置 Hosts

需要配置 `hosts` 文件，将 `proxy.server` 指向反向代理服务器的 IP 地址。

规则请参考 [https://jihulab.com/haozi-team/mirror-conf](https://jihulab.com/haozi-team/mirror-conf)。

### 4. 配置 .env

```bash
cp .env.example .env
```

自行修改 `.env` 文件中的配置。

### 5. 初始化数据库

```bash
./weavatar artisan migrate
./weavatar artisan hash:generate
./weavatar artisan hash:table
```

其中 `hash:generate` 步骤可选，用于生成 QQ 邮箱的 Hash 表，约占用 150 GB。

导入步骤所使用的 TiDB-Lightning 配置文件参考如下

```yaml
[ lightning ]
  # 日志
  level = "info"
  file = "tidb-lightning.log"

  [ tikv-importer ]
  # 选择使用的导入模式
  backend = "local"
  # 设置排序的键值对的临时存放地址，目标路径需要是一个空目录
  sorted-kv-dir = "/tmp/kvtmp"

  [ mydumper ]
  data-source-dir = "/www/weavatar/hash"
  strict-format = true

  [ mydumper.csv ]
  # 字段分隔符，支持一个或多个字符，默认值为 ','。如果数据中可能有逗号，建议源文件导出时分隔符使用非常见组合字符例如'|+|'。
  separator = ','
  # 引用定界符，设置为空表示字符串未加引号。
  delimiter = '"'
  # 行尾定界字符，支持一个或多个字符。设置为空（默认值）表示 "\n"（换行）和 "\r\n" （回车+换行），均表示行尾。
  terminator = ""
  # CSV 文件是否包含表头。
  # 如果为 true，首行将会被跳过，且基于首行映射目标表的列。
  header = false
  # CSV 是否包含 NULL。
  # 如果为 true，CSV 文件的任何列都不能解析为 NULL。
  not-null = true
  # 如果 `not-null` 为 false（即 CSV 可以包含 NULL），
  # 为以下值的字段将会被解析为 NULL。
  null = '\N'
  # 是否解析字段内的反斜线转义符。
  backslash-escape = true
  # 是否移除以分隔符结束的行。
  trim-last-separator = false

  [ tidb ]
  # 目标集群的信息
  host = "127.0.0.1"
  port = 4000
  user = "root"
  password = ""
  status-port = 10080
  # 集群 pd 的地址
  pd-addr = "127.0.0.1:2379"
```

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
}

location /avatar {
    rewrite ^/avatar(.*)$ /api/avatar$1 last;
}
```

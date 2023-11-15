## 手动运行
本示例教程将从源代码编译本程序，并运行程序到`/www/wwwroot/autobuytgvip/`目录

### 准备工作
- 至少1个telegram机器人
- 域名
- mysql
- redis
- nginx

### 导入数据库
创建数据库，并导入数据库脚本`doc/btp_saas.sql`

### 编译源代码
```shell
cd ~
# 拉取代码
git clone https://github.com/buyaobilian1/autobuytgvip-saas.git

# 创建运行目录
mkdir -p /www/wwwroot/autobuytgvip/etc

# 编译售卖端
cd autobuytgvip-saas/btp-saas
go mod tidy
go build
cp btp-saas /www/wwwroot/autobuytgvip/
cp etc/btp-saas.yaml.example /www/wwwroot/autobuytgvip/etc/btp-saas.yaml

# 编译代理端
cd ../btp-agent
go mod tidy
go build
cp btp-agent /www/wwwroot/autobuytgvip/
cp etc/btp-agent.yaml.example /www/wwwroot/autobuytgvip/etc/btp-agent.yaml

```

### 修改配置文件
btp-saas.yaml
```yaml
Name: btp-saas
Host: 0.0.0.0
Port: 8001

# mysql配置
MysqlConf:
  Host: 127.0.0.1:3306
  User: user
  Pass: pass
  Db: db

# redis配置
RedisConf:
  Host: 127.0.0.1:6379
  Type: node
  Pass: pass
  Tls: true

# epusdt 配置
PayConf:
  BaseApi: http://127.0.0.1:8001
  ApiToken: apiToken
  NotifyUrl: http://127.0.0.1:8001/pay/notify/%s

AppConf:
  # 开会员订单过期时间(分)
  OrderExpireMinute: 9
  # 充值订单过期时间(分)
  RechargeExpireMinute: 9
  # 代理地址
  ProxyUrl: socks5://127.0.0.1:1080
  # TON助记词
  TonSeed: birth pattern then forest walnut then phrase walnut fan pumpkin pattern then cluster blossom verify then forest velvet pond fiction pattern collect then then
  # fragment hash
  Hash: f8e8690959d5b64095
  # fragment cookie
  Cookie: stel_ssid=052c1dad9226b29029_9517489419833914253; stel_dt=-480; stel_ton_token=PlrMkcQGmSKJU12o2SzqAEgOE-gpIWxSbo3EyeNShSyvI4pQkYURLOTuXY-v6DZYTTtehJ9SmGOZq0vNEEwcsKSTgJ84UZ3yf23vGuH-WD21Ib067kBGBrrpsxUT3OtMUD-ZFRUYdMSvu1NUlTh9-JhJt-4QUZmWnSt2E_nzYuLa_CprRTQ

```

btp-agent.yaml
```yaml
Name: btp-agent
Host: 0.0.0.0
Port: 8002

# mysql配置
MysqlConf:
  Host: 127.0.0.1:3306
  User: btp_saas
  Pass: 123456
  Db: btp_saas

AppConf:
  # 代理管家机器人TOKEN
  BotToken: 1234567890:sdfasdfasdfasfasdfasdfasdfasdfasdfasdf
  # 代理地址
  ProxyUrl: socks5://127.0.0.1:7890
  # 机器人回调地址，saas程序的地址
  WebhookUrl: http://127.0.0.1:8001/webhook/%d
  # webhook密钥
  WebhookSecret: AbAb123456
```

### nginx配置
新建网站，并代理到btp-saas运行的端口

### 运行
```shell
cd /www/wwwroot/autobuytgvip/
nohup btp-saas &
nohup btp-agent &
```

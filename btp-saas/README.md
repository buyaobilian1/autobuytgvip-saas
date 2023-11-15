
## generate project
```shell
mkdir btp-saas && cd btp-saas
goctl api -o btp-saas.api
goctl api go --api btp-saas.api --dir .
```

## generate gorm code using gentool
[gentool](https://github.com/go-gorm/gen/tree/master/tools/gentool)
```shell
# install gentool
go install gorm.io/gen/tools/gentool@latest
# generate models 
gentool -dsn "user:pass@(ip:port)/db?charset=utf8mb4&parseTime=True&loc=Local" -tables "user, order"
gentool -dsn "user:pass@(ip:port)/db?charset=utf8mb4&parseTime=True&loc=Local" -onlyModel -fieldNullable



```

## build linux
```shell
set GOARCH=amd64
set GOOS=linux
go build
```

## 存在的问题
1. TON余额不足提醒
2. 购买成功后给代理结算费用
3. 
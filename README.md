# tx_base_api
使用golang+echo+grom+configor+dep框架搭建的程序

# 正常运行命令
go main.go

# 热加载框架 realize 
使用框架  https://github.com/tockins/realize
go get github.com/tockins/realize
自动热加载(自动重启)使用以下命令启动
realize run 

# 配置文件在 configs文件夹下 
使用框架 https://github.com/jinzhu/configor
配置结构体文件 config.go
默认使用development配置文件运行，切换运行环境命令使用 CONFIGOR_ENV=production go run main.go

# 使用dep框架管理第三方依赖包 
使用框架 https://github.com/golang/dep
go get -u github.com/golang/dep/cmd/dep

增加一个依赖
$ dep ensure -add github.com/foo/bar
或者手动修改 Gopkg.toml配置文件
然后执行
$ dep ensure

# orm框架使用 
使用框架 https://github.com/jinzhu/configor

# GoCleanArchitecture
# Introduction 
&emsp;&emsp;本篇用Golang以clean architecture的方式開發一個會員系統的後端，利用Restful API跟gRPC證實clean architecture在外部資料變動下核心邏輯還是穩健的，有使用到的特性與套件請看Features。

# Clean Architecture
Clean architecture架構可以帶來以下好處:
* 獨立於框架
* 可測試
* 獨立於 UI
* 獨立於 Database
* 獨立於任何外部代理

# Features    
* RESTful API
* gRPC
* Gin
* Gorm
* jwt-go
* go-swagger
* golang-migrate
* Docker 

# Enviroment
* Golang 1.17+
* Mysql 8.0  

# Build
1.將.env與config的default去掉並設定完成<br>
2.執行docker-compose即可建置完成
<div><pre>docker compose up
</pre></div>

# Command
Generate migrate file
<div><pre>migrate create -ext sql -dir database/migrations {file_name}
</pre></div>

Generate protor file<br>
如果要用到proto的其他included就必須要將proto path設定到指定路徑<br>
ex: 想要將gRPC的Function傳入設置為空的話就必須要引用google.protobuf.Empty
<div><pre>protoc *.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path={proto_included_path}/include/ --proto_path=.
</pre></div>

Generate swagger file<br>
<div><pre>swagger generate spec -o ./swagger.yml</pre></div>

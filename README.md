
1.配置 `configs/config.yaml`文件
   
   ```yaml 本地
   mysql:
   -  name: "default"
      username: "root"
      password: "123456"
      database: "db_ginadmin"
      host: "127.0.0.1"
      port: 3306
      max_open_conn: 50
      max_idle_conn: 20
   redis:
      addr: "localhost:6379"
      db: 0
      password: ""
   session:
      session_name: "gosession_id"
   base:
      host: 0.0.0.0
      port: 20011
      log_media: "redis"
   ```

2. 運行 `go run .\cmd\ginadmin`访问地址 http://localhost:端口地址/admin/login。默认账户：admin  密码：111111


### :small_blue_diamond:<a name="docker-compose">构建开发环境</a>

1. 替换conf目录下的配置项
   
   ```yaml
   mysql:
   -  name: "default"
      username: "docker"
      password: "123456"
      database: "docker_mysql"
      host: "localmysql"
      port: 3306
      max_open_conn: 50
      max_idle_conn: 20
   redis:
    -  addr: "localredis:6379"
      db: 0
      password: "ginadmin"
   session:
    -  session_name: "gosession_id"
   base:
    -  host: 0.0.0.0
      port: 20010
      log_media: "redis"
   ```

2. 执行命令 `docker-compose up -d`

3. 进入到容器中 `docker exec -it ginadmin-web bash`

4. 下载扩展 `go mod tidy`

5. 运行项目 ` go run ./cmd/ginadmin/ run`  访问地址 `http://localhost:20010/admin/login`



### :small_blue_diamond:<a name="分页">分页</a>

 使用 `pkg/paginater/paginater.go` 里面的 `PageOperation` 进行分页
   
   ```go
   adminDb := models.Db.Table("admin_users").Select("nickname","username").Where("uid != ?", 1)
   adminUserData := paginater.PageOperation(c, adminDb, 1, &adminUserList)
   ```


### :small_blue_diamond:<a name="日志">日志</a>

1. 系统日志
   
   设置路由中间件来收集系统日志和错误日志，设置 `internal/router/default.go` 文件

2. 自定义日志
   
   使用 loggers.LogInfo()` 方法记录日志  `eGame-demo-back-office-api/pkg/loggers`
   
   ```golang
   loggers.LogInfo("admin", "this is a info message", map[string]string{
           "user_info": "this is a user info",
   })
   ```

3. 切换存储介质
   
   在配置文件中修改 `log_media` 参数默认file文件存储可选redis存储

### :small_blue_diamond:<a name="数据库">数据库</a>

1. models下定义的文件均需要实现 `TableName() string`  方法，并将实现该结构体的指针写入到 `GetModels` 方法中
   
   ```go
   func GetModels() []interface{} {
       return []interface{}{
           &AdminUsers{},
           &AdminGroup{},
       }
   }
   ```

2. model需要继承 BaseModle 并且实现 TableName 方法，如果需要初始化填充数据的话，需要实现 FillData() 方法，并将数据填充需要执行的代码写到函数体里。详情参照 AdminUsers

### :small_blue_diamond:<a name="定时任务">定时任务</a>

- 在 `pkg/cron/cron.go`  添加定时执行任务

### :small_blue_diamond:<a name="配置文件">配置文件</a>

1. 现在 `configs/config.go` 添加配置项的 struct 类型，例如
   
   ```go
   type AppConf struct {
       BaseConf `yaml:"base"`
   }
   type BaseConf struct {
       Port string `yaml:"port"`
   }
   ```

2. 在 `configs/config.yaml` 添加配置信息
   
   ```
   base:
      port: 20011
   ```

3. 在代码中调用配置文件的信息
   
   ```go
   configs.App.BaseConf.Port
   ```

### :small_blue_diamond:<a name="线上部署">线上部署</a>

- 使用 `go build .\cmd\ginadmin`  生成二进制文件
- 打包静态资源部署 `go build -tags=embed .\cmd\ginadmin` 

### :small_blue_diamond:<a name="命令行操作">命令行操作</a>

*  运行程序命令
```
PS F:\ginadmin> go run .\cmd\ginadmin\ run -h
Run app

Usage:
  ginadmin run [flags]

Flags:
  -c, --config path string   config path
  -h, --help                 help for run
  -m, --mode string          dev or release (default "dev")
```
* 数据表迁移命令
```
PS F:\ginadmin> go run .\cmd\ginadmin\ db migrate -h
DB Migrate

Usage:
  ginadmin db migrate [-t table] [flags]

Flags:
  -c, --config path string   config path
  -h, --help                 help for migrate
  -t, --table string         input a table name
```

*  数据填充命令
```
PS F:\ginadmin> go run .\cmd\ginadmin\ db seed -h   
DB Seed

Usage:
  ginadmin db seed [-t table] [flags]

Flags:
  -c, --config path string   config path
  -h, --help                 help for seed
  -t, --table string         input a table name
```
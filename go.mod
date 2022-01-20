module github.com/swhc/utils

go 1.14

require (
	github.com/FZambia/sentinel v1.1.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gomodule/redigo v1.8.8
	github.com/shima-park/agollo v1.2.12
	github.com/spf13/viper v1.10.1
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.8.2
	go.uber.org/zap v1.20.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	xorm.io/xorm v1.2.5
)

replace github.com/spf13/viper => github.com/spf13/viper v1.7.1

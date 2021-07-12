module github.com/gxravel/bus-routes

go 1.16

<<<<<<< HEAD
require github.com/streadway/amqp v1.0.0
=======
require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi v1.5.4
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-swagger/go-swagger v0.27.0
	github.com/golangci/golangci-lint v1.41.1
	github.com/gxravel/bus-routes/pkg/rmq v0.0.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lopezator/migrator v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.8.1
	github.com/streadway/amqp v1.0.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
)

replace github.com/gxravel/bus-routes/pkg/rmq v0.0.0 => ../bus-routes/pkg/rmq
>>>>>>> 301fe6e (feat(amqp): support amqp)

// Package classification Member system API.
//
// The purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Schemes: http
// Host: localhost
// Base path: /api/v1
// Version: 0.0.1
// Consumes: application/json
// Produces: application/json
// Security: Bearer
// securityDefinitions:
//    Bearer:
//      type: apiKey
//      name: Authorization
//      in: header
// swagger:meta
package main

import (
	httpDelivery "GoCleanArchitecture/delivery/http"
	httpMiddlewares "GoCleanArchitecture/delivery/http/middlewares"
	repo "GoCleanArchitecture/repository"
	usecase "GoCleanArchitecture/usecase"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
	goMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbName := viper.GetString("database.name")
	jwtKey := viper.GetString("key.jwtKey")

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(goMySQL.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("get db failed:", err)
		return
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Hour)

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		fmt.Println("driver error: " + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	if err := m.Up(); err != nil {
		fmt.Println("migrate up error: " + err.Error())
	}

	redisDB := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	tokenRepository := repo.NewTokenRepository(db, redisDB)
	tokenUsecase := usecase.NewTokenUsecase(jwtKey, tokenRepository)

	userRepo := repo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// http
	router := gin.Default()
	httpDelivery.NewTokenHandler(router, tokenUsecase)
	authMiddleware := httpMiddlewares.NewAuthMiddlewares(tokenUsecase)
	httpDelivery.NewUserHandler(router, userUsecase, tokenUsecase, authMiddleware)

	router.Run(":8080")

	// grpc
	// grpcAuthMiddleware := grpcMiddlewares.NewAuthMiddlewares(tokenUsecase)
	// s := grpc.NewServer(grpc.UnaryInterceptor(grpcAuthMiddleware.UnaryServerInterceptor))

	// grpcDelivery.NewTokenHandler(s, tokenUsecase)
	// grpcDelivery.NewUserHandler(s, userUsecase, tokenUsecase)

	// lis, err := net.Listen("tcp", ":8081")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}

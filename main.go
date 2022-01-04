package main

import (
	grpcDelivery "GoCleanArchitecture/delivery/grpc"
	grpcMiddlewares "GoCleanArchitecture/delivery/grpc/middlewares"
	repo "GoCleanArchitecture/repository"
	usecase "GoCleanArchitecture/usecase"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile("config.yaml")
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	tokenRepository := repo.NewTokenRepository(db)
	tokenUsecase := usecase.NewTokenUsecase(jwtKey, tokenRepository)

	userRepo := repo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// http
	// router := gin.Default()
	// delivery.NewTokenHandler(router, tokenUsecase)
	// authMiddleware := middlewares.NewAuthMiddlewares(tokenUsecase)
	// delivery.NewUserHandler(router, userUsecase, tokenUsecase, authMiddleware)

	// router.Run(":8080")

	// grpc
	grpcAuthMiddleware := grpcMiddlewares.NewAuthMiddlewares(tokenUsecase)
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcAuthMiddleware.UnaryServerInterceptor))

	grpcDelivery.NewTokenHandler(s, tokenUsecase)
	grpcDelivery.NewUserHandler(s, userUsecase, tokenUsecase)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

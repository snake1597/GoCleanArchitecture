package main

import (
	grpcDelivery "GoCleanArchitecture/delivery/grpc"
	grpcMiddlewares "GoCleanArchitecture/delivery/grpc/middlewares"
	delivery "GoCleanArchitecture/delivery/http"
	middlewares "GoCleanArchitecture/delivery/http/middlewares"
	"GoCleanArchitecture/entities"
	repo "GoCleanArchitecture/repository"
	usecase "GoCleanArchitecture/usecase"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
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

	db.AutoMigrate(&entities.User{})

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxLifetime(time.Hour)

	router := gin.Default()
	jwtKey := viper.GetString("key.jwtKey")

	tokenRepository := repo.NewTokenRepository(db)
	tokenUsecase := usecase.NewTokenUsecase(jwtKey, tokenRepository)
	delivery.NewTokenHandler(router, tokenUsecase)

	authMiddleware := middlewares.NewAuthMiddlewares(tokenUsecase)

	userRepo := repo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	delivery.NewUserHandler(router, userUsecase, tokenUsecase, authMiddleware)

	// grpc
	grpcAuthMiddleware := grpcMiddlewares.NewAuthMiddlewares(tokenUsecase)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpcAuthMiddleware.UnaryServerInterceptor))
	//s := grpc.NewServer()
	grpcDelivery.NewUserHandler(s, userUsecase, tokenUsecase)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	router.Run(":8080")
}

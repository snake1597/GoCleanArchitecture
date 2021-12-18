package main

import (
	delivery "GoCleanArchitecture/delivery/http"
	middlewares "GoCleanArchitecture/delivery/http/middlewares"
	"GoCleanArchitecture/entities"
	repo "GoCleanArchitecture/repository"
	usecase "GoCleanArchitecture/usecase"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

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

	router.Run(":8080")
}

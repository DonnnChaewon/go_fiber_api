package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/go_fiber_api?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fail to connect database:", err)
	}

	// Auto migrate models
	db.AutoMigrate(&User{}, &Movie{})

	return db
}
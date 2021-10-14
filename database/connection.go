package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jwtproject/models"
)

var DB *gorm.DB

func Connect(){
	connection,err:=gorm.Open(mysql.Open("root:00@/jwttoken"),&gorm.Config{})
	if err != nil {
		panic("couldn't connect to db")
	}
	DB=connection
	//creating the table
	connection.AutoMigrate(&models.User{})
}

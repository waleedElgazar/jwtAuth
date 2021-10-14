package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"jwtproject/database"
	"jwtproject/models"
	"strconv"
	"time"
)
var secretKey="helloWaleed"

func Register(ctx *fiber.Ctx)error{
	var data map[string]string
	err:=ctx.BodyParser(&data)
	if err != nil {
		fmt.Println("error",err)
	}
	password,err:=bcrypt.GenerateFromPassword([]byte(data["password"]),14)
	user := models.User{
		Name: data["name"],
		Email: data["email"],
		Password: password,
	}
	database.DB.Create(&user)
	return ctx.JSON(user)
}

func Login(ctx *fiber.Ctx)error{
	var data map[string]string
	err:=ctx.BodyParser(&data)
	if err != nil {
		fmt.Println("error",err)
	}
	var user models.User
	//get first row where email equal what u want
	database.DB.Where("email=?",data["email"]).First(&user)
	//that mean the user hadn't been founded
	if user.Id==0 {
		ctx.Status(fiber.StatusNotFound)
		ctx.JSON(fiber.Map{
			"message":"user not found",
		})
		return err
	}
	//to compare the password but the stored password had been encrypted
	err=bcrypt.CompareHashAndPassword(user.Password,[]byte(data["password"]))
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		ctx.JSON(fiber.Map{
			"message":"incorrect password",
		})
		return err
	}
	claims:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Minute*20).Unix(),
	})
	token,err:=claims.SignedString([]byte(secretKey))
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		ctx.JSON(fiber.Map{
			"message":"couldn't login",
		})
		return err
	}
	cookie:=fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Minute*20),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return ctx.JSON(fiber.Map{
		"message":"login",
	})
}

func User(ctx *fiber.Ctx)error{
	//must run login function to get the token first
	cookie:=ctx.Cookies("jwt")
	token,err:=jwt.ParseWithClaims(cookie,&jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey),nil
	})
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.JSON(fiber.Map{
			"message":"Unauthorized",
		})
	}
	claims:=token.Claims.(*jwt.StandardClaims)
	var user models.User
	database.DB.Where("id=?",claims.Issuer).First(&user)
	return ctx.JSON(user)


}

func LogOut(ctx *fiber.Ctx)error{
	//first we will remove the cookie we will create another cookie and set the expired time in the past
	cookie:=fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return ctx.JSON(fiber.Map{
		"message":"success",
	})
}
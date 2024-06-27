package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Parva-Parmar/GO-Auth/database"
	"github.com/Parva-Parmar/GO-Auth/helpers"
	helper "github.com/Parva-Parmar/GO-Auth/helpers"
	"github.com/Parva-Parmar/GO-Auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword()

func VerifyPassord()

func Signup() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var user models.User 

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return
		}

		count,err := UserCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking for email"})
		}

		count,err = UserCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured  while checking for the Phone"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"this email of phone number already exists"})
		}
	}
}

func Login()

func GetUsers()

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")

		if err := helper.MatchUserTypetoUid(c,userId);err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var user models.User
		err := UserCollection.FindOne(ctx,bson.M{"user_id":userId}).Decode(&user)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOk,user)
	}
}
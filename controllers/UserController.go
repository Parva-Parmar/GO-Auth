package controllers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Parva-Parmar/GO-Auth/database"
	helper "github.com/Parva-Parmar/GO-Auth/helpers"
	"github.com/Parva-Parmar/GO-Auth/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/check.v1"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword(password string) string{
	bcrypt.GenerateFromPassword([]byte(password),14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassord(userPassword string,providedPassword string)(bool,string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword),[]byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
}

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


		password := HashPassword(*user.Password)
		user.Password = &password
		count,err := UserCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured  while checking for the Phone"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"this email of phone number already exists"})
		}

		user.Created_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_at,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token,refershToken,_ := helper.GenerateAllTokens(*user.Email,*user.First_name,*user.Last_name,*user.User_type,*&user.User_id)
		user.Token = &token
		user.Refersh_token = &refershToken

		resultInsertionNumber, insertErr := UserCollection.InsertOne(ctx,user)
		if insertErr != nil{
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOk,resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"email of password is incorrect"})
			return
		}

		passwordIsValid,msg := VerifyPassord(*user.Password,*foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user not found"})
		}
		token,refreshToken,_ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *&foundUser.Last_name,*&foundUser.User_type,foundUser.User_id)
		helper.UpdateAllTokens(token,refreshToken,foundUser.User_id)
		UserCollection.FindOne(ctx,bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,foundUser)
	}
}

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
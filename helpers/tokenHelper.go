package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/Parva-Parmar/GO-Auth/database"
	"github.com/dgrijalva/jwt-go"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetailes struct{
	Email 			string
	First_name 		string
	Last_name 		string
	Uid 			string
	User_type 		string
	jwt.StandardClaims
}

var UserCollection *mongo.Collection = database.OpenCollection(database.Client,"user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string)(signedToken string,signedRefreshToken string){
	claims := &SignedDetailes{
		Email : email,
		First_name: firstName,
		Last_name: lastName,
		Uid : uid,
		User_type: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetailes{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt : time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SigendString([]byte(SECRET_KEY))
	refershToken,err := jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SigendString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token,refershToken,err
}
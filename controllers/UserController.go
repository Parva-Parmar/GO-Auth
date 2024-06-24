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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword()

func VerifyPassord()

func Signup()

func Login()

func GetUsers()

func GetUser()
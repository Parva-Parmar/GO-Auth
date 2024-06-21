package routes

import (
	controller "github.com/Parva-Parmar/GO-Auth/controlllers"
	"github.com/Parva-Parmar/GO-Auth/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incommingRoutes *gin.Engine){
	incommingRoutes.Use(middleware.Authenticate())
	incommingRoutes.GET("/users",controller.GetUsers())
	incommingRoutes.GET("/users/:user_id",controller.GetUsers())
}
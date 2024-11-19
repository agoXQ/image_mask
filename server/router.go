package server

import (
	"picmask/controllers"
	// "picmask/imageUtils"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	image := r.Group("/image")
	{
		image.POST("/mask/:id/:operate/:passwd", controllers.Img{}.GetImg)
		image.GET("/download/:id/:filename", controllers.Img{}.DownImg)
	}
	return r
}

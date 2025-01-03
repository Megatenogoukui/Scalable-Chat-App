package Routers

import (
	"Backend/Routers/Websocket"

	"github.com/gin-gonic/gin"
)

func MapRoutes(router *gin.Engine) {
	Websocket.SetupDropDownRoutes(router)
}

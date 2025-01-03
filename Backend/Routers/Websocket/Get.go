package Websocket

import (
	upgrade_to_websocket "Backend/Controllers"
	get_messages "Backend/Controllers/GetMessages"
	signup_login_user "Backend/Controllers/SignUpLoginUser"
	users "Backend/Controllers/Users"

	"github.com/gin-gonic/gin"
)

func SetupDropDownRoutes(router *gin.Engine) {
	router.GET("/websocket_connection", upgrade_to_websocket.Handle_socket_connection)
	router.GET("/get_messages", get_messages.GetMessage)
	router.POST("/sign_up", signup_login_user.SignUpUser)
	router.POST("/login", signup_login_user.LoginUser)
	router.GET("/GetAllUsers", users.GetAllUsers)

}

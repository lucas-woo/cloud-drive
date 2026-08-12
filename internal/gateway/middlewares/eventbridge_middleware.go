package middlewares

import "github.com/gin-gonic/gin"


type EventBridgeMiddleware struct {}

func (m *EventBridgeMiddleware)EventBridgeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func NewEventbridgeMiddleware() *EventBridgeMiddleware{
	return &EventBridgeMiddleware{}
}
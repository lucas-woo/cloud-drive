package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



type HealthHandler struct {

}

func (h *HealthHandler) Check(c *gin.Context) {
	c.Status(http.StatusOK)
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
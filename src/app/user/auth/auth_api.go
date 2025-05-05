package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MedodsTechTask/app/core"
	auth "github.com/MedodsTechTask/app/user/auth/share"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.GET(core.UserAuthHealthCheck, healthCheck)
	r.POST(core.UserAuthSignUpEmail, signupEmail)
}

// @Summary Проверка работоспособности сервиса
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Router /user/auth/health [get]
func healthCheck(c *gin.Context) {

}

// @Summary Регистрация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.QEmailSignup true "Данные регистрации"
// @Success 200 {object} auth.ZEmailSignup
// @Router /user/auth/signup/email [post]
func signupEmail(c *gin.Context) {
	var req auth.QEmailSignup

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

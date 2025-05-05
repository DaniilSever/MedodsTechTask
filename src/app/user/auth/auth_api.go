package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/MedodsTechTask/app/core"
)

func SetupRoutes(r *gin.RouterGroup) {
	r.GET(core.UserAuthHealthCheck, healthCheck)
	r.POST(core.UserAuthSignUpEmail, signupEmail)
	r.POST(core.UserAuthConfirmEmail)
	r.POST(core.UserAuthLoginEmail)
	r.POST(core.UserAuthRefreshToken)
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

}

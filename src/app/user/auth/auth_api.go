package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MedodsTechTask/app/core"
	share "github.com/MedodsTechTask/app/user/auth/share"
)

type API struct {
	uc *AuthUseCase
}

func NewAPI(uc *AuthUseCase) *API {
	return &API{uc: uc}
}

func (h *API) SetupRoutes(r *gin.RouterGroup) {
	r.GET(core.UserAuthHealthCheck, h.healthCheck)
	r.POST(core.UserAuthSignUpEmail, h.signupEmail)
}

// @Summary Проверка работоспособности сервиса
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Router /user/auth/health [get]
func (h *API) healthCheck(c *gin.Context) {

}

// @Summary Регистрация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body share.QEmailSignup true "Данные регистрации"
// @Success 200 {object} share.ZEmailSignup
// @Router /user/auth/signup/email [post]
func (h *API) signupEmail(c *gin.Context) {
	var req share.QEmailSignup

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.SignupEmail(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

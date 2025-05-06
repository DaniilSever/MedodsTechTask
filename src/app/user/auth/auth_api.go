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
	r.POST(core.UserAuthSignUpEmail, h.signupEmail)
	r.POST(core.UserAuthConfirmEmail, h.confirmEmail)
	r.POST(core.UserAuthLoginEmail, h.loginEmail)
	r.POST(core.UserAuthRefreshToken, h.refreshToken)
}

// @Summary Регистрация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body share.QEmailSignup true "Данные регистрации"
// @Success 200 {object} share.ZEmailSignup
// @Failure 400 {object} core.ZError
// @Failure 500 {object} core.ZError
// @Router /user/auth/signup/email [post]
func (h *API) signupEmail(c *gin.Context) {
	var req share.QEmailSignup

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.SignupEmail(c.Request.Context(), &req)
	if err != nil {
		switch err.Code {
		case 400:
			c.JSON(http.StatusBadRequest, err)
			return
		case 500:
			c.JSON(http.StatusInternalServerError, err)
			return
		}

	}

	c.JSON(http.StatusOK, res)
}

// @Summary Подтверждение регистрации
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body share.QConfirmEmail true "Данные подтверждения"
// @Success 200 {object} share.ZAccount
// @Failure 400 {object} core.ZError
// @Failure 404 {object} core.ZError
// @Failure 500 {object} core.ZError
// @Router /user/auth/confirm/email [post]
func (h *API) confirmEmail(c *gin.Context) {
	var req share.QConfirmEmail

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.ConfirmEmail(c.Request.Context(), &req)
	if err != nil {
		switch err.Code {
		case 400:
			c.JSON(http.StatusBadRequest, err)
			return
		case 404:
			c.JSON(http.StatusNotFound, err)
			return
		case 500:
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Вход в аккаунт через email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body share.QLoginEmail true "Данные аккаунта"
// @Success 200 {object} share.ZToken
// @Failure 400 {object} core.ZError
// @Failure 404 {object} core.ZError
// @Failure 500 {object} core.ZError
// @Router /user/auth/login/email [post]
func (h *API) loginEmail(c *gin.Context) {
	var req share.QLoginEmail

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.LoginEmail(c.Request.Context(), &req, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		switch err.Code {
		case 400:
			c.JSON(http.StatusBadRequest, err)
			return
		case 404:
			c.JSON(http.StatusNotFound, err)
			return
		case 500:
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Рефреш токена
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body share.QRefreshToken true "Токен"
// @Success 200 {object} share.ZToken
// @Failure 400 {object} core.ZError
// @Failure 404 {object} core.ZError
// @Failure 500 {object} core.ZError
// @Router /user/auth/refresh/token [post]
func (h *API) refreshToken(c *gin.Context) {
	var req share.QRefreshToken

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.RefreshToken(c.Request.Context(), &req, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		switch err.Code {
		case 400:
			c.JSON(http.StatusBadRequest, err)
			return
		case 404:
			c.JSON(http.StatusNotFound, err)
			return
		case 500:
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

package handler

import (
	"net/http"
	"tasklist/internal/models"

	"github.com/gin-gonic/gin"
	"tasklist/internal/service"
)

type AuthHandler struct {
	svc service.Auth
}

func NewAuthHandler(svc service.Auth) *AuthHandler { return &AuthHandler{svc: svc} }

func (h *AuthHandler) Register(api *gin.RouterGroup) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", h.RegisterUser)
		authGroup.POST("/login", h.Login)
	}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Creates a new user account with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  models.AuthRequest  true  "User credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Register(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{Token: token})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  models.AuthRequest  true  "User credentials"
// @Success      200  {object}  models.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.svc.Login(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}

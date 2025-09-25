package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "tasklist/docs"
	"tasklist/internal/service"
)

type Handler struct {
	AuthService service.Auth
	TaskService service.Task
	log         *logrus.Logger
}

func NewHandler(service *service.Service, log *logrus.Logger) *Handler {
	return &Handler{
		AuthService: service.AuthService,
		TaskService: service.TaskService,
		log:         log,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler := NewAuthHandler(h.AuthService)
	taskHandler := NewTaskHandler(h.TaskService, h.log)

	api := router.Group("/api")
	{
		authHandler.Register(api)
		taskHandler.Register(api)
	}

	return router
}

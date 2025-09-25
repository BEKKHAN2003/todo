package handler

import (
	"net/http"
	"strconv"
	"tasklist/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"tasklist/internal/models"
	"tasklist/internal/service"
)

type TaskHandler struct {
	svc service.Task
	log *logrus.Logger
}

func NewTaskHandler(svc service.Task, log *logrus.Logger) *TaskHandler {
	return &TaskHandler{svc: svc, log: log}
}

func (h *TaskHandler) Register(rg *gin.RouterGroup) {
	taskGroup := rg.Group("/tasks", middleware.JWTAuth())
	{
		taskGroup.POST("", h.Create)
		taskGroup.GET("", h.List)
		taskGroup.GET("/:id", h.GetByID)
		taskGroup.POST("/:id/complete", h.Complete)
		taskGroup.PUT("/:id", h.Update)
		taskGroup.DELETE("/:id", h.Delete)
	}
}

// Create godoc
// @Summary      Create task
// @Description  Create a new task item
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body  models.TaskRequest  true  "Task data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	userId := c.GetInt(models.UserCtxKey)
	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.Create(c, userId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "task created successfully"})
}

// Complete godoc
// @Summary      Complete task
// @Description  Complete task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /tasks/{id}/complete [post]
func (h *TaskHandler) Complete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = h.svc.Complete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "task completed successfully"})
}

// List godoc
// @Summary      Get tasks
// @Description  Get all tasks for authenticated user
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.Task
// @Failure      401  {object}  map[string]string
// @Router       /tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	userId := c.GetInt(models.UserCtxKey)
	tasks, err := h.svc.List(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetByID godoc
// @Summary      Get task by ID
// @Description  Get a single task by its ID
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	userId := c.GetInt(models.UserCtxKey)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.svc.GetByID(c, id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

// Update godoc
// @Summary      Update task
// @Description  Update task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int          true  "Task ID"
// @Param        input  body      models.TaskRequest  true  "Task data"
// @Success      200    {object}  models.TaskRequest
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	userId := c.GetInt(models.UserCtxKey)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(c, id, userId, req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

// Delete godoc
// @Summary      Delete task
// @Description  Delete a task by its ID
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.svc.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

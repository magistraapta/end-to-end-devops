package handlers

import (
	"code/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) Index(c *gin.Context) {
	users, err := h.getUsers()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to fetch users")
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"users": users,
	})
}

func (h *UserHandler) List(c *gin.Context) {
	h.renderUsersTable(c)
}

func (h *UserHandler) Create(c *gin.Context) {
	user := models.User{
		ID:       uuid.NewString(),
		Name:     c.PostForm("name"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.String(http.StatusBadRequest, "name, email, and password are required")
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "failed to create user")
		return
	}

	h.renderUsersTable(c)
}

func (h *UserHandler) EditForm(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.DB.First(&user, "id = ?", id).Error; err != nil {
		c.String(http.StatusNotFound, "user not found")
		return
	}

	c.HTML(http.StatusOK, "partials/edit_form.html", gin.H{
		"user": user,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	updates := map[string]any{
		"name":     c.PostForm("name"),
		"email":    c.PostForm("email"),
		"password": c.PostForm("password"),
	}

	if updates["name"] == "" || updates["email"] == "" || updates["password"] == "" {
		c.String(http.StatusBadRequest, "name, email, and password are required")
		return
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.String(http.StatusBadRequest, "failed to update user")
		return
	}

	h.renderUsersTable(c)
}

func (h *UserHandler) ClearEditForm(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		c.String(http.StatusBadRequest, "failed to delete user")
		return
	}

	h.renderUsersTable(c)
}

func (h *UserHandler) renderUsersTable(c *gin.Context) {
	users, err := h.getUsers()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to fetch users")
		return
	}

	c.HTML(http.StatusOK, "partials/users_table.html", gin.H{
		"users": users,
	})
}

func (h *UserHandler) getUsers() ([]models.User, error) {
	var users []models.User
	err := h.DB.Order("created_at desc").Find(&users).Error
	return users, err
}

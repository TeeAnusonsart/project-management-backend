package http

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.usecase.ListUsers()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []fiber.Map
	for _, u := range users {
		response = append(response, fiber.Map{
			"user_id":      u.ID,
			"username":     u.Username,
			"name":         u.Name,
			"email":        u.Email,
			"profile_path": u.ProfilePath,
		})
	}

	return c.JSON(fiber.Map{"data": response})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		Password    string `json:"password"`
		Email       string `json:"email"`
		ProfilePath string `json:"profile_path"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user := &domain.User{
		Username:    req.Username,
		Name:        req.Name,
		Password:    req.Password,
		Email:       req.Email,
		ProfilePath: req.ProfilePath,
	}

	if err := h.usecase.CreateUser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เพิ่มผู้ใช้งานเรียบร้อยแล้ว",
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	user, err := h.usecase.GetUser(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"user_id":      user.ID,
		"username":     user.Username,
		"name":         user.Name,
		"email":        user.Email,
		"profile_path": user.ProfilePath,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	if err := h.usecase.DeleteUser(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "ลบบัญชีผู้ใช้งานเรียบร้อยแล้ว",
	})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.ChangePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เปลี่ยนรหัสผ่านเรียบร้อยแล้ว",
	})
}

func (h *UserHandler) UploadProfile(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	file, err := c.FormFile("profile")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "profile file is required")
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	filename := fmt.Sprintf("user_%d_%s", id, file.Filename)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := h.usecase.UploadProfile(uint(id), savePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message":      "อัปโหลดรูปโปรไฟล์เรียบร้อยแล้ว",
		"profile_path": savePath,
	})
}

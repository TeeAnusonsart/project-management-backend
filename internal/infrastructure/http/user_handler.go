package http

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

// --- Request DTOs ---

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"` // Added validate:"email" to check email format
}

// --- Handler Methods ---

// [GET] /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.usecase.ListUsers()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve users")
	}
	response := make([]UserResponse, 0)
	for _, u := range users {
		response = append(response, ToUserResponse(u))
	}

	return sendResponse(c, fiber.StatusOK, "Users retrieved successfully", response)
}

// [POST] /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	// 1. Validate JSON payload
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	// 2. Structural validation (Required, Email, etc.)
	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	// 3. Map to Domain Model
	user := &domain.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  domain.RoleUser,
	}

	// 4. Save to Database
	if err := h.usecase.CreateUser(user); err != nil {
		// You can pass err.Error() if you want the frontend to see the exact error, 
		// or keep it generic to hide database details.
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendResponse(c, fiber.StatusCreated, "User created successfully", nil)
}

// [GET] /users/:email
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	email := c.Params("email")

	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email parameter is missing")
	}

	user, err := h.usecase.GetUser(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return sendResponse(c, fiber.StatusOK, "User retrieved successfully", ToUserResponse(user))
}

// [DELETE] /users/:email
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	email := c.Params("email")

	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email parameter is missing")
	}

	if err := h.usecase.DeleteUser(email); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user account")
	}

	return sendResponse(c, fiber.StatusOK, "User account deleted successfully", nil)
}

// --- Response Helpers & DTOs ---

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		Name:  u.Name,
		Email: u.Email,
		Role:  string(u.Role),
	}
}
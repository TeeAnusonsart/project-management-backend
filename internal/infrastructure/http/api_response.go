package http

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
    Status  string      `json:"status"`
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"` // omitempty จะซ่อนฟิลด์นี้ถ้าเป็น nil
}

// ฟังก์ชัน Helper สำหรับส่ง Response แบบรวดเร็ว
func sendResponse(c *fiber.Ctx, code int, message string, data interface{}) error {
    status := "success"
    if code >= 400 {
        status = "error"
    }
    
    return c.Status(code).JSON(APIResponse{
        Status:  status,
        Code:    code,
        Message: message,
        Data:    data,
    })
}
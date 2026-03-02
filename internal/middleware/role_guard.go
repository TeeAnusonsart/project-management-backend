package middleware

import (
    "github.com/gofiber/fiber/v2"
	"project-home-iot/internal/core/domain"
)

// สร้าง Interface สำหรับค้นหา Role (หรือจะส่ง DB เข้าไปตรงๆ ก็ได้)
type UserGetter interface {
    GetRoleByEmail(email string) (string, error)
}

func RoleGuard(userRepo domain.UserRepository, allowedRoles ...domain.UserRole) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // ดึง email ที่ FirebaseAuth ทำการ Verify และฝากไว้ใน Locals แล้ว
        email, ok := c.Locals("email").(string)
        if !ok || email == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Email not found"})
        }

        // Query หา User จาก DB
        user, err := userRepo.FindByEmail(email)
        if err != nil {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "User not found in system"})
        }

        // ตรวจสอบ Role
        for _, allowed := range allowedRoles {
            if allowed == user.Role {
                c.Locals("user_role", user.Role)
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "message": "คุณไม่มีสิทธิ์เข้าถึงส่วนนี้ (Admin Only)",
        })
    }
}
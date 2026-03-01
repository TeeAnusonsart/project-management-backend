package middleware

import (
    "context"
    "strings"

    "firebase.google.com/go/v4/auth"
    "github.com/gofiber/fiber/v2"
)

func FirebaseAuth(authClient *auth.Client) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        token, err := authClient.VerifyIDToken(context.Background(), tokenString)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
        }

        // 1. ดึง Email จาก Token
        email, ok := token.Claims["email"].(string)
        if !ok || email == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email not found in token"})
        }

        // 2. กฎเหล็ก: ต้องเช็คว่าอีเมลนี้ถูก Verify แล้วจริงๆ (OAuth มักจะ true เสมอ)
        emailVerified, ok := token.Claims["email_verified"].(bool)
        if !ok || !emailVerified {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Email is not verified"})
        }

        // 3. เก็บแค่ Email ลงใน Context ก็พอแล้ว
        c.Locals("email", email)

        return c.Next()
    }
}
package middleware

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func RejectDuplicateJSONKeys() fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Body()

		// ถ้าไม่ใช่ JSON → ข้าม
		if len(body) == 0 || !json.Valid(body) {
			return c.Next()
		}

		dec := json.NewDecoder(bytes.NewReader(body))

		// อ่าน token แรก
		t, err := dec.Token()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid JSON")
		}

		// เช็คว่าเป็น object ไหม
		if delim, ok := t.(json.Delim); !ok || delim != '{' {
			return c.Next()
		}

		keys := make(map[string]bool)

		for dec.More() {
			token, err := dec.Token()
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid JSON format")
			}

			key := token.(string)

			if keys[key] {
				return fiber.NewError(
					fiber.StatusBadRequest,
					"duplicate key: "+key,
				)
			}
			keys[key] = true

			// skip value
			var v interface{}
			if err := dec.Decode(&v); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid JSON value")
			}
		}

		return c.Next()
	}
}
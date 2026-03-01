package http

import (
	"fmt"
	"project-home-iot/internal/core/domain" // เรียกใช้ domain เพื่อเอาโครงสร้าง FieldError

	"github.com/go-playground/validator/v10"
)

// สร้าง instance ของ validator ไว้ใช้ร่วมกันทั้ง package
var validate = validator.New()

// validateStruct ตรวจสอบ struct และคืนค่ารายการ error กลับไป (ถ้ามี)
// หมายเหตุ: สมมติว่าใน domain.ErrorResponse ของคุณประกาศชนิดของ Errors เป็น []domain.FieldError
func validateStruct(s interface{}) []domain.FieldError {
	var validationErrors []domain.FieldError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var elem domain.FieldError
			elem.Field = err.Field()
			elem.Message = getErrorMessage(err)
			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

// getErrorMessage แปลง Tag ของ Validator เป็นข้อความภาษาไทยที่อ่านง่าย
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("ฟิลด์ %s จำเป็นต้องระบุ", err.Field())
	case "min":
		return fmt.Sprintf("ฟิลด์ %s ต้องมีความยาวอย่างน้อย %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("ฟิลด์ %s ต้องมีความยาวไม่เกิน %s", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("ฟิลด์ %s ต้องเป็นค่าใดค่าหนึ่งใน: %s", err.Field(), err.Param())
	}
	return err.Error() // กรณีที่ไม่ได้ดัก Tag ไว้ ให้แสดง error พื้นฐาน
}
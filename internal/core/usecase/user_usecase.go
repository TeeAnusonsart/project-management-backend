package usecase

import (
	"errors"

	"project-home-iot/internal/core/domain"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	ListUsers() ([]*domain.User, error)
	GetUser(email string) (*domain.User, error)
	DeleteUser(email string) error
}

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) UserUsecase {
	return &userUsecase{userRepo: ur}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	exist, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if exist != nil {
		return errors.New("Email already exist")
	}
	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	if !user.Role.IsValid() {
		return errors.New("invalid role")
	}

	return u.userRepo.Create(user)
}


func (u *userUsecase) ListUsers() ([]*domain.User, error) {
	return u.userRepo.FindAll()
}

// func (u *userUsecase) GetUser(id uint) (*domain.User, error) {
// 	return u.userRepo.FindByID(id)
// }

func (u *userUsecase) GetUser(email string) (*domain.User, error) {
    user, err := u.userRepo.FindByEmail(email)
    if err != nil {
        return nil, errors.New("unauthorized: email not found in system")
    }
    return user, nil
}

// func (u *userUsecase) VerifyAndLinkUser(email string, firebaseUID string) (*domain.User, error) {
//     // 1. ค้นหาว่า Admin ได้สร้าง Email นี้ไว้ในระบบหรือยัง?
//     user, err := u.userRepo.FindByEmail(email)
//     if err != nil {
//         // ถ้าไม่เจอ แสดงว่า Admin ยังไม่ได้สร้างให้
//         return nil, errors.New("user not found. please contact admin")
//     }

//     // 2. ถ้าเจอ แล้ว FirebaseUID ยังว่างอยู่ (แปลว่าเพิ่งเข้าระบบครั้งแรก) ให้ทำการผูกบัญชี
//     if user.UID == "" {
//         err = u.userRepo.UpdateUID(user.Email, firebaseUID)
//         if err != nil {
//             return nil, err
//         }
//         user.UID = firebaseUID // อัปเดต Object ให้เป็นค่าล่าสุด
//     }else if user.UID != firebaseUID {
//         return nil, errors.New("invalid account mapping")
//     }

//     return user, nil
// }

func (u *userUsecase) DeleteUser(email string) error {
	return u.userRepo.Delete(email)
}


// func (u *userUsecase) UploadProfile(id uint, filePath string) error {
// 	return u.userRepo.UpdateProfilePath(id, filePath)
// }

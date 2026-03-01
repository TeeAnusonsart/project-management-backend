package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Provider สำหรับ Wire
func NewFirebaseAuthClient() (*auth.Client, error) {
	ctx := context.Background()
    // ใน production ควรใช้ Environment Variable
	opt := option.WithCredentialsFile("firebase-account.json") 
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	return app.Auth(ctx)
}
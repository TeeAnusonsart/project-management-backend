//go:build wireinject
// +build wireinject

package main

import (
	"project-home-iot/internal/auth"

	"github.com/google/wire"
	"gorm.io/gorm"
)

type AppHandlers struct {
	Auth *auth.AuthHandler
}

func InitializeApp(db *gorm.DB) *AppHandlers {
	wire.Build(
		auth.ProviderSet,
		wire.Struct(new(AppHandlers), "*"),
	)
	return &AppHandlers{}
}

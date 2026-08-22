// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD
// Created: 2026-08-20
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package auth

import (
	"context"
	"octacart/internal/shared"
)

type FormAuthRepo interface {
	CreateProfile(ctx context.Context, profile *FormAuthProfile) error
	FindProfile(ctx context.Context, userId shared.UserId) (*FormAuthProfile, error)
}

type PasswordSecEstimator interface {
	EstimateSecurity(ctx context.Context, password string) error
}

type OidAuthRepo interface {
	CreateProfile(ctx context.Context, profile *OidAuthProfile) error
	FindByProvider(ctx context.Context, providerName string, providerKey string) (*OidAuthProfile, error)
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id shared.UserId) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password string, hashed string) (bool, error)
}

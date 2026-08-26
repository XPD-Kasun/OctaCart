// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD
// Created: 2026-08-20
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package auth

import (
	"context"
	"octacart/internal/shared"
	"time"

	"github.com/rs/zerolog"
)

type FormAuthProfile struct {
	userId       shared.UserId
	passwordHash string
}

type FormAuthConfig struct {
	MinPwLength   int
	LockoutDur    time.Duration
	MaxLockoutDur time.Duration
}

type FormAuthSvc struct {
	logger         *zerolog.Logger // will change to abstraction?
	config         *FormAuthConfig
	hasher         PasswordHasher
	formAuthRepo   FormAuthRepo
	userRepo       UserRepo
	pwSecEstimator PasswordSecEstimator
}

func NewFormAuthSvc(
	config *FormAuthConfig,
	hasher PasswordHasher,
	formAuthRepo FormAuthRepo,
	userRepo UserRepo,
	passSecEstimator PasswordSecEstimator) *FormAuthSvc {

	return &FormAuthSvc{
		config:         config,
		hasher:         hasher,
		formAuthRepo:   formAuthRepo,
		userRepo:       userRepo,
		pwSecEstimator: passSecEstimator,
	}
}

func (fa *FormAuthSvc) Authenticate(
	ctx context.Context,
	email string,
	password string) (*User, error) {

	user, err := fa.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	lockoutDur := int(fa.config.LockoutDur.Seconds() * 2 * float64(user.lockout.lockoutCount))
	userLocked := user.LockoutAttempt(time.Duration(lockoutDur) * time.Second)
	err = fa.userRepo.Update(ctx, user) // need to save the user after locking attempted.

	if err != nil {
		fa.logger.Error().
			Str("err", err.Error()).
			Msg("Could not save locked user state while authenticating at UserRepo.Update")
	}

	if userLocked {
		return nil, UserLocked
	}

	userPWProfile, err := fa.formAuthRepo.FindProfile(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	isOk, err := fa.hasher.CheckPassword(password, userPWProfile.passwordHash)
	if err != nil {
		return nil, err
	}

	if !isOk {
		return nil, InvalidCred
	}

	return user, nil

}

func (fa *FormAuthSvc) CreateProfile(
	ctx context.Context,
	userId shared.UserId,
	password string) (*FormAuthProfile, error) {

	if err := fa.pwSecEstimator.EstimateSecurity(ctx, password); err != nil {
		return nil, err
	}

	hash, err := fa.hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	profile := &FormAuthProfile{
		userId:       userId,
		passwordHash: hash,
	}

	err = fa.formAuthRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (fa *FormAuthSvc) GetProfile(ctx context.Context, userId shared.UserId) (*FormAuthProfile, error) {

	user, err := fa.userRepo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}

	userPWProfile, err := fa.formAuthRepo.FindProfile(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return userPWProfile, nil

}

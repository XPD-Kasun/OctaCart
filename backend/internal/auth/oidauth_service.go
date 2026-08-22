// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD
// Created: 2026-08-22
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package auth

import (
	"context"
	"octacart/internal/shared"
)

type OidAuthProfile struct {
	userId       shared.UserId
	providerName string
	providerKey  string
}

type OidAuthConfig struct {
	DefaultIssuer string
}

type OidAuthSvc struct {
	config    *OidAuthConfig
	authStore OidAuthRepo
	userRepo  UserRepo
}

func NewOidAuthSvc(
	config *OidAuthConfig,
	authStore OidAuthRepo,
	userRepo UserRepo) *OidAuthSvc {

	return &OidAuthSvc{
		config:    config,
		authStore: authStore,
		userRepo:  userRepo,
	}
}

func (oa *OidAuthSvc) Authenticate(
	ctx context.Context,
	providerName string,
	providerKey string) (*User, error) {

	if providerKey == "" {
		providerKey = oa.config.DefaultIssuer
	}

	profile, err := oa.authStore.FindByProvider(ctx, providerName, providerKey)
	if err != nil {
		return nil, err
	}

	user, err := oa.userRepo.FindById(ctx, profile.userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (oa *OidAuthSvc) GetProfileByProvider(
	ctx context.Context,
	providerName string,
	providerKey string) (*OidAuthProfile, error) {

	if providerKey == "" {
		providerKey = oa.config.DefaultIssuer
	}

	return oa.authStore.FindByProvider(ctx, providerName, providerKey)
}

func (oa *OidAuthSvc) CreateProfile(
	ctx context.Context,
	userId shared.UserId,
	providerName string,
	providerKey string) (*OidAuthProfile, error) {

	profile := &OidAuthProfile{
		userId:       userId,
		providerName: providerName,
		providerKey:  providerKey,
	}

	if profile.providerKey == "" {
		profile.providerKey = oa.config.DefaultIssuer
	}

	err := oa.authStore.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

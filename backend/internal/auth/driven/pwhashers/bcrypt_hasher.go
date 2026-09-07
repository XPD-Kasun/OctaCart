// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD, Claude Sonnet
// Created: 2026-08-22
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.
// AI Generated for [auth.PasswordHasher] targeting bcrypt, reviewed and owned by XPD.

package pwhashers

import (
	"octacart/internal/auth"

	goBcrypt "golang.org/x/crypto/bcrypt"
)

// BCryptHasher is a driven adapter implementing [auth.PasswordHasher]
// using the bcrypt algorithm. Place the cost (rounds) between 10-14
// for a good balance between security and latency.
type BCryptHasher struct {
	rounds int
}

func NewBCryptHasher(rounds int) *BCryptHasher {
	return &BCryptHasher{rounds: rounds}
}

// HashPassword implements [auth.PasswordHasher].
func (b *BCryptHasher) HashPassword(password string) (string, error) {
	hashed, err := goBcrypt.GenerateFromPassword([]byte(password), b.rounds)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword implements [auth.PasswordHasher].
func (b *BCryptHasher) CheckPassword(password string, hashed string) (bool, error) {
	err := goBcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err == goBcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// compile-time check that BCryptHasher satisfies the port.
var _ auth.PasswordHasher = (*BCryptHasher)(nil)

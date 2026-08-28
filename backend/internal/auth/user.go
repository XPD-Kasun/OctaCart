// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD
// Created: 2026-08-20
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

// The auth package implements the authentication bounded context for OctaCart.
// It defines domain entities [User], [Lockout].[FormAuthProfile], [OidAuthProfile] as the
// value objects and application services such as [FormAuthSvc] for driving adapters.
// For ports and errors refer ports.go and errors.go files. Adapters should refer these errors when
// applicable.
package auth

import (
	"octacart/internal/shared"
	"time"
)

type Claim string

type Lockout struct {
	lockoutEnabled bool
	lockoutCount   int
	lockoutTill    time.Time
}

func (l *Lockout) LockoutAttempt(lockoutDuration time.Duration) bool {
	if !l.lockoutEnabled {
		return false
	}
	l.lockoutCount--
	if l.lockoutCount < 0 {
		l.lockoutCount = 0
		l.lockoutTill = time.Now().Add(lockoutDuration)
		return true
	}
	return false
}

func (l *Lockout) IsLockedOut() bool {
	return l.lockoutEnabled && l.lockoutCount <= 0
}

func (l *Lockout) LockoutTill() time.Time {
	return l.lockoutTill
}

//func (l *Lockout) Try

type Principal interface {
	Claims() []Claim
	IsAuthenticated() bool
}

type User struct {
	Id              shared.UserId
	Name            string
	Email           string
	claims          []Claim
	isAuthenticated bool
	lockout         Lockout
}

func (u *User) LockoutAttempt(duration time.Duration) bool {
	return u.lockout.LockoutAttempt(duration)
}

func (u *User) ResetLockout(maxAttempts int) {
	u.lockout.lockoutCount = maxAttempts
	u.lockout.lockoutTill = time.Now().Add(-time.Hour)
}

func (u *User) IsAuthenticated() bool {
	return u.isAuthenticated
}

func (u *User) Claims() []Claim { return u.claims }

var AnonymousUser *User = &User{Id: 0, isAuthenticated: false, Name: "Anon"}

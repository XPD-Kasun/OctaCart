// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD
// Created: 2026-08-28
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package auth

import (
	"octacart/internal/shared"
	"testing"
	"time"
)

func TestLockout_LockoutAttempt(t *testing.T) {
	t.Run("lockout disabled/should not lock out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: false, lockoutCount: 3, lockoutTill: time.Now().Add(time.Hour)}
		result := lockout.LockoutAttempt(5 * time.Minute)
		if result {
			t.Error("LockoutAttempt should return false when lockout is disabled")
		}
		if lockout.lockoutCount != 3 {
			t.Errorf("lockoutCount should remain unchanged when disabled, got %d, want 3", lockout.lockoutCount)
		}
	})

	t.Run("lockout enabled/count above zero/should not lock out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: 3, lockoutTill: time.Now()}
		result := lockout.LockoutAttempt(5 * time.Minute)
		if result {
			t.Error("LockoutAttempt should return false when count is above zero")
		}
		if lockout.lockoutCount != 2 {
			t.Errorf("lockoutCount should decrement, got %d, want 2", lockout.lockoutCount)
		}
	})

	t.Run("lockout enabled/count at one/should lock out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: 1, lockoutTill: time.Now()}
		result := lockout.LockoutAttempt(5 * time.Minute)
		if result {
			t.Error("LockoutAttempt should return false when count is positive including 1")
		}
		if lockout.lockoutCount != 0 {
			t.Errorf("lockoutCount should be zero, got %d, want 0", lockout.lockoutCount)
		}
	})

	t.Run("lockout enabled/count at zero/should lock out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: 0, lockoutTill: time.Now()}
		result := lockout.LockoutAttempt(5 * time.Minute)
		if !result {
			t.Error("LockoutAttempt should return true when count is zero")
		}
		if lockout.lockoutCount != 0 {
			t.Errorf("lockoutCount should be zero, got %d, want 0", lockout.lockoutCount)
		}
		if lockout.lockoutTill.Before(time.Now()) {
			t.Error("lockoutTill should be in the future after lockout")
		}
	})

	t.Run("lockout enabled/count below zero/should lock out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: 0, lockoutTill: time.Now()}
		result := lockout.LockoutAttempt(5 * time.Minute)
		if !result {
			t.Error("LockoutAttempt should return true when count is already zero")
		}
		if lockout.lockoutCount != 0 {
			t.Errorf("lockoutCount should remain zero, got %d, want 0", lockout.lockoutCount)
		}
	})
}

func TestLockout_IsLockedOut(t *testing.T) {
	t.Run("lockout disabled/not locked out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: false, lockoutCount: 0, lockoutTill: time.Now().Add(1 * time.Hour)}
		result := lockout.IsLockedOut()
		if result {
			t.Error("IsLockedOut should return false when lockout is disabled")
		}
	})

	t.Run("lockout enabled/count above zero/not locked out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: 1, lockoutTill: time.Now().Add(1 * time.Hour)}
		result := lockout.IsLockedOut()
		if result {
			t.Error("IsLockedOut should return false when count is above zero")
		}
	})

	t.Run("lockout enabled/count at zero/locked out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: 0, lockoutTill: time.Now().Add(1 * time.Hour)}
		result := lockout.IsLockedOut()
		if !result {
			t.Error("IsLockedOut should return true when count is zero and lockout is enabled")
		}
	})

	t.Run("lockout enabled/count below zero/locked out", func(t *testing.T) {
		lockout := Lockout{lockoutEnabled: true, lockoutCount: -1, lockoutTill: time.Now().Add(1 * time.Hour)}
		result := lockout.IsLockedOut()
		if !result {
			t.Error("IsLockedOut should return true when count is negative and lockout is enabled")
		}
	})
}

func TestLockout_LockoutTill(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(10 * time.Minute)
	lockout := Lockout{lockoutEnabled: true, lockoutCount: 0, lockoutTill: futureTime}

	result := lockout.LockoutTill()
	if !result.Equal(futureTime) {
		t.Errorf("LockoutTill should return the stored time, got %v, want %v", result, futureTime)
	}
}

func TestUser_LockoutAttempt(t *testing.T) {
	t.Run("user lockout", func(t *testing.T) {
		user := User{
			Id:              shared.UserId(1),
			Name:            "Test User",
			Email:           "test@example.com",
			isAuthenticated: true,
			lockout:         Lockout{lockoutEnabled: true, lockoutCount: 0, lockoutTill: time.Now().Add(time.Hour)},
		}

		result := user.LockoutAttempt(5 * time.Minute)
		if !result {
			t.Error("User.LockoutAttempt should delegate to Lockout.LockoutAttempt and return true when lockout occurs")
		}
		if user.lockout.lockoutCount != 0 {
			t.Errorf("User lockout count should be updated, got %d, want 0", user.lockout.lockoutCount)
		}
	})
}

func TestUser_ResetLockout(t *testing.T) {
	t.Run("reset lockout", func(t *testing.T) {
		beforeReset := time.Now().Add(time.Hour)
		user := User{
			Id:              shared.UserId(1),
			Name:            "Test User",
			Email:           "test@example.com",
			isAuthenticated: true,
			lockout:         Lockout{lockoutEnabled: true, lockoutCount: 0, lockoutTill: beforeReset},
		}

		user.ResetLockout(5)

		if user.lockout.lockoutCount != 5 {
			t.Errorf("lockoutCount should be reset, got %d, want 5", user.lockout.lockoutCount)
		}
		if !user.lockout.lockoutTill.Before(time.Now()) {
			t.Error("lockoutTill should be updated to some previous time after reset(so not locked)")
		}
	})
}

// Additional stuff. Better to have than not.
func TestUser_IsAuthenticated(t *testing.T) {
	t.Run("authenticated uaer", func(t *testing.T) {
		user := User{
			Id:              shared.UserId(1),
			Name:            "Test User",
			Email:           "test@example.com",
			isAuthenticated: true,
		}
		if !user.IsAuthenticated() {
			t.Error("IsAuthenticated should return true for auth user")
		}
	})

	t.Run("unuthenticated user", func(t *testing.T) {
		user := User{
			Id:              shared.UserId(1),
			Name:            "Test User",
			Email:           "test@example.com",
			isAuthenticated: false,
		}
		if user.IsAuthenticated() {
			t.Error("IsAuthenticated should return false for unauth user")
		}
	})
}

func TestUser_Claims(t *testing.T) {
	t.Run("return user claims", func(t *testing.T) {
		expectedClaims := []Claim{"admin", "user"}
		user := User{
			Id:              shared.UserId(1),
			Name:            "Test User",
			Email:           "test@example.com",
			claims:          expectedClaims,
			isAuthenticated: true,
		}

		result := user.Claims()
		if len(result) != len(expectedClaims) {
			t.Errorf("Claims should return the correct number of claims, got %d, want %d", len(result), len(expectedClaims))
		}
		for i, claim := range expectedClaims {
			if result[i] != claim {
				t.Errorf("Claims[%d] should match, got %v, want %v", i, result[i], claim)
			}
		}
	})

	t.Run("return empty claims for user with no claims", func(t *testing.T) {
		user := User{
			Id:              shared.UserId(1),
			Name:            "Test User",
			Email:           "test@example.com",
			claims:          []Claim{},
			isAuthenticated: true,
		}

		result := user.Claims()
		if len(result) != 0 {
			t.Errorf("Claims should return empty slice, got %d claims", len(result))
		}
	})
}

func TestAnonymousUser(t *testing.T) {
	t.Run("AnonymousUser is not authenticated", func(t *testing.T) {
		if AnonymousUser.IsAuthenticated() {
			t.Error("AnonymousUser should not be authenticated")
		}
	})

	t.Run("AnonymousUser has correct ID", func(t *testing.T) {
		if AnonymousUser.Id != shared.UserId(0) {
			t.Errorf("AnonymousUser ID should be 0, got %d", AnonymousUser.Id)
		}
	})

	t.Run("AnonymousUser has correct name", func(t *testing.T) {
		if AnonymousUser.Name != "Anon" {
			t.Errorf("AnonymousUser Name should be 'Anon', got %s", AnonymousUser.Name)
		}
	})

	t.Run("AnonymousUser has empty claims", func(t *testing.T) {
		claims := AnonymousUser.Claims()
		if len(claims) != 0 {
			t.Errorf("AnonymousUser should have no claims, got %d", len(claims))
		}
	})
}

func TestUser_ImplementsPrincipal(t *testing.T) {
	// Verify that User impl [Principal] interface
	var _ Principal = &User{}
	var _ Principal = AnonymousUser
}

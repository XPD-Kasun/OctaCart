// Copyright 2026 OctaCart. All rights reserved.
// Author: XPD
// Created: 2026-08-20
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package auth

import "errors"

var (
	UserLocked  = errors.New("user locked")
	InvalidCred = errors.New("invalid credentials")

	// [PasswordSecEstimator] port implementors should use this error in case if need to
	// give an error for short passwords than length etc. We suggest to wrap the error with specifics.
	// for example in adapter,
	// if len(password) < b.minLen {
	// 		return fmt.Errorf("password must be at least %d characters: %w", b.minLen, auth.PassLen)
	// }
	// This allows, errors.Is(err, auth.PassLen) {} for callers.
	PassLen = errors.New("password too short")

	// [PasswordSecEstimator] port implementors should use this error
	// when they found a breached password in db or remote rainbow table or
	PassBreach = errors.New("password insecure")

	// [PasswordSecEstimator] port implementors should use this error when they found no profile
	// for the given user. Allows the callers to attach a profile etc..
	ProfileNotFound = errors.New("auth profile not found for provider")
)

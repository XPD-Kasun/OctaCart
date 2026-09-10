// Copyright 2026 OctaCart. All rights reserved.
// Author: OctaCart Team
// Created: 2026-09-07
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package product

import "errors"

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrVariantNotFound        = errors.New("variant not found")
	ErrDuplicateSKU           = errors.New("duplicate SKU within product")
	ErrInvalidTransition      = errors.New("invalid status transition")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrCategoryNotFound       = errors.New("category not found")
	ErrSlugConflict           = errors.New("slug already exists")
	ErrVariantHasActiveOrders = errors.New("variant has active orders")
	ErrCategoryInUse          = errors.New("category is assigned to products")
	ErrMediaNotFound          = errors.New("media not found")
)
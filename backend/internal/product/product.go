// Copyright 2026 OctaCart. All rights reserved.
// Author: OctaCart Team
// Created: 2026-09-07
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

// Package product provides the core domain model and business logic for the
// product catalog bounded context, including products, variants, categories,
// media assets, and dynamic attributes.
package product

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ID types representing opaque entity identifiers.
type ProId int
type ProVariantId int
type CatId int
type MediaId int

// ProductStatus defines lifecycle states of a product.
type ProductStatus string

const (
	StatusDraft    ProductStatus = "Draft"
	StatusActive   ProductStatus = "Active"
	StatusArchived ProductStatus = "Archived"
)

type val struct {
	Type  string
	Value any
}

// Attributes represents an opaque dynamic attribute container.
type Attributes struct {
	attrs map[string]*val
}

// NewAttributes initializes an empty Attributes container.
func NewAttributes() *Attributes {
	return &Attributes{
		attrs: make(map[string]*val),
	}
}

// AddStr sets a string dynamic attribute.
func (a *Attributes) AddStr(name, value string) {
	a.attrs[name] = &val{Type: "str", Value: value}
}

// AddNum sets a numeric dynamic attribute.
func (a *Attributes) AddNum(name string, value float64) {
	a.attrs[name] = &val{Type: "num", Value: value}
}

// AddRange sets an integer range dynamic attribute.
func (a *Attributes) AddRange(name string, min, max int) {
	a.attrs[name] = &val{Type: "range", Value: fmt.Sprintf("%d,%d", min, max)}
}

// AddEnum sets an enumeration dynamic attribute with pipe-joined values.
func (a *Attributes) AddEnum(name string, values ...string) {
	a.attrs[name] = &val{Type: "enum", Value: strings.Join(values, "|")}
}

// GetStr retrieves a string attribute. Returns error if missing or wrong type.
func (a *Attributes) GetStr(name string) (string, error) {
	v, exists := a.attrs[name]
	if !exists {
		return "", fmt.Errorf("attribute %q not found", name)
	}
	if v.Type != "str" {
		return "", fmt.Errorf("attribute %q has type %s, expected str", name, v.Type)
	}
	s, ok := v.Value.(string)
	if !ok {
		return "", fmt.Errorf("attribute %q value is not a string", name)
	}
	return s, nil
}

// GetNum retrieves a numeric attribute. Returns error if missing or wrong type.
func (a *Attributes) GetNum(name string) (float64, error) {
	v, exists := a.attrs[name]
	if !exists {
		return 0, fmt.Errorf("attribute %q not found", name)
	}
	if v.Type != "num" {
		return 0, fmt.Errorf("attribute %q has type %s, expected num", name, v.Type)
	}
	n, ok := v.Value.(float64)
	if !ok {
		return 0, fmt.Errorf("attribute %q value is not float64", name)
	}
	return n, nil
}

// GetRange retrieves an integer range attribute. Returns error if missing or invalid.
func (a *Attributes) GetRange(name string) (min, max int, err error) {
	v, exists := a.attrs[name]
	if !exists {
		return 0, 0, fmt.Errorf("attribute %q not found", name)
	}
	if v.Type != "range" {
		return 0, 0, fmt.Errorf("attribute %q has type %s, expected range", name, v.Type)
	}
	strVal, ok := v.Value.(string)
	if !ok {
		return 0, 0, fmt.Errorf("attribute %q range value is not formatted string", name)
	}
	parts := strings.Split(strVal, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("attribute %q invalid range format: %s", name, strVal)
	}
	minVal, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("attribute %q min parse error: %w", name, err)
	}
	maxVal, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("attribute %q max parse error: %w", name, err)
	}
	return minVal, maxVal, nil
}

// GetEnum retrieves an enum slice attribute. Returns error if missing or wrong type.
func (a *Attributes) GetEnum(name string) ([]string, error) {
	v, exists := a.attrs[name]
	if !exists {
		return nil, fmt.Errorf("attribute %q not found", name)
	}
	if v.Type != "enum" {
		return nil, fmt.Errorf("attribute %q has type %s, expected enum", name, v.Type)
	}
	strVal, ok := v.Value.(string)
	if !ok {
		return nil, fmt.Errorf("attribute %q enum value is not formatted string", name)
	}
	if strVal == "" {
		return []string{}, nil
	}
	return strings.Split(strVal, "|"), nil
}

// Remove deletes an attribute by name.
func (a *Attributes) Remove(name string) {
	delete(a.attrs, name)
}

// Names returns a sorted slice of all attribute keys.
func (a *Attributes) Names() []string {
	keys := make([]string, 0, len(a.attrs))
	for k := range a.attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Has checks if an attribute exists.
func (a *Attributes) Has(name string) bool {
	_, exists := a.attrs[name]
	return exists
}

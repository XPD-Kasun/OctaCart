// Copyright 2026 OctaCart. All rights reserved.
// Author: OctaCart Team
// Created: 2026-09-07
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package product

import (
	"reflect"
	"testing"
)

func TestAttributes(t *testing.T) {
	t.Run("new/should be empty", func(t *testing.T) {
		attrs := NewAttributes()
		names := attrs.Names()
		if len(names) != 0 {
			t.Errorf("expected empty names, got %v", names)
		}
	})

	t.Run("AddStr and GetStr", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddStr("Color", "Red")
		val, err := attrs.GetStr("Color")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "Red" {
			t.Errorf("expected Red, got %s", val)
		}
	})

	t.Run("AddNum and GetNum", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddNum("Weight", 1.5)
		val, err := attrs.GetNum("Weight")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 1.5 {
			t.Errorf("expected 1.5, got %f", val)
		}
	})

	t.Run("AddRange and GetRange", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddRange("Size", 1, 10)
		minVal, maxVal, err := attrs.GetRange("Size")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if minVal != 1 || maxVal != 10 {
			t.Errorf("expected (1, 10), got (%d, %d)", minVal, maxVal)
		}
	})

	t.Run("AddEnum and GetEnum", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddEnum("Size", "S", "M", "L")
		vals, err := attrs.GetEnum("Size")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"S", "M", "L"}
		if !reflect.DeepEqual(vals, expected) {
			t.Errorf("expected %v, got %v", expected, vals)
		}
	})

	t.Run("GetStr missing/should error", func(t *testing.T) {
		attrs := NewAttributes()
		_, err := attrs.GetStr("nope")
		if err == nil {
			t.Error("expected error for missing key, got nil")
		}
	})

	t.Run("Remove/should delete", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddStr("Tag", "Val")
		if !attrs.Has("Tag") {
			t.Fatal("expected tag to exist")
		}
		attrs.Remove("Tag")
		if attrs.Has("Tag") {
			t.Error("expected tag to be removed")
		}
	})

	t.Run("Names/should return sorted", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddStr("Z", "z")
		attrs.AddStr("A", "a")
		attrs.AddStr("M", "m")
		names := attrs.Names()
		expected := []string{"A", "M", "Z"}
		if !reflect.DeepEqual(names, expected) {
			t.Errorf("expected %v, got %v", expected, names)
		}
	})

	t.Run("overwrite/should replace", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddStr("X", "a")
		attrs.AddNum("X", 1)

		numVal, err := attrs.GetNum("X")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if numVal != 1 {
			t.Errorf("expected 1, got %f", numVal)
		}

		_, err = attrs.GetStr("X")
		if err == nil {
			t.Error("expected error when getting as string after overwriting with num, got nil")
		}
	})
}

func TestProductStatus(t *testing.T) {
	t.Run("constants/should match string values", func(t *testing.T) {
		if StatusDraft != "Draft" {
			t.Errorf("expected StatusDraft to be 'Draft', got %s", StatusDraft)
		}
		if StatusActive != "Active" {
			t.Errorf("expected StatusActive to be 'Active', got %s", StatusActive)
		}
		if StatusArchived != "Archived" {
			t.Errorf("expected StatusArchived to be 'Archived', got %s", StatusArchived)
		}
	})
}

// Copyright 2026 OctaCart. All rights reserved.
// Author: OctaCart Team
// Created: 2026-09-07
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package shared

import (
	"testing"
)

func TestPagination_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		perPage  int
		expected int
	}{
		{
			name:     "Page=1, PerPage=10 -> Offset = 0",
			page:     1,
			perPage:  10,
			expected: 0,
		},
		{
			name:     "Page=3, PerPage=20 -> Offset = 40",
			page:     3,
			perPage:  20,
			expected: 40,
		},
		{
			name:     "Page=1, PerPage=1 -> Offset = 0",
			page:     1,
			perPage:  1,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				Page:    tt.page,
				PerPage: tt.perPage,
			}
			actual := p.Offset()
			if actual != tt.expected {
				t.Errorf("expected offset %d, got %d", tt.expected, actual)
			}
		})
	}
}

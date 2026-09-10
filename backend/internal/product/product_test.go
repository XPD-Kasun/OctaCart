// Copyright 2026 OctaCart. All rights reserved.
// Author: OctaCart Team
// Created: 2026-09-07
//
// This file is part of the OctaCart e-Commerce platform. Refer to licensing for usage.

package product

import (
	"errors"
	"reflect"
	"testing"
	"time"
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

func TestNewProduct(t *testing.T) {
	t.Run("valid input/should create draft product", func(t *testing.T) {
		tags := []string{"sale", "new"}
		p, err := NewProduct("My Phone", "Great phone", CatId(10), false, tags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Status() != StatusDraft {
			t.Errorf("expected status %s, got %s", StatusDraft, p.Status())
		}
		if p.Title() != "My Phone" {
			t.Errorf("expected title 'My Phone', got %s", p.Title())
		}
		if p.Description() != "Great phone" {
			t.Errorf("expected description 'Great phone', got %s", p.Description())
		}
		if p.CategoryId() != CatId(10) {
			t.Errorf("expected categoryId 10, got %d", p.CategoryId())
		}
		if p.IsDigital() != false {
			t.Errorf("expected isDigital false, got true")
		}
		if !reflect.DeepEqual(p.Tags(), tags) {
			t.Errorf("expected tags %v, got %v", tags, p.Tags())
		}
		if p.CreatedAt().IsZero() {
			t.Error("expected non-zero createdAt")
		}
		if p.UpdatedAt().IsZero() {
			t.Error("expected non-zero updatedAt")
		}
	})

	t.Run("empty title/should return error", func(t *testing.T) {
		p, err := NewProduct("", "desc", CatId(1), false, nil)
		if err == nil {
			t.Error("expected error for empty title, got nil")
		}
		if p != nil {
			t.Errorf("expected nil product, got %v", p)
		}

		p2, err2 := NewProduct("   ", "desc", CatId(1), false, nil)
		if err2 == nil {
			t.Error("expected error for whitespace title, got nil")
		}
		if p2 != nil {
			t.Errorf("expected nil product, got %v", p2)
		}
	})
}

func TestProduct_Publish(t *testing.T) {
	t.Run("draft product/should transition to active", func(t *testing.T) {
		p, err := NewProduct("Title", "Desc", CatId(1), false, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err = p.Publish()
		if err != nil {
			t.Fatalf("unexpected error on publish: %v", err)
		}
		if p.Status() != StatusActive {
			t.Errorf("expected status %s, got %s", StatusActive, p.Status())
		}
	})

	t.Run("active product/should return ErrInvalidTransition", func(t *testing.T) {
		p, _ := NewProduct("Title", "Desc", CatId(1), false, nil)
		_ = p.Publish()
		err := p.Publish()
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got %v", err)
		}
	})

	t.Run("archived product/should return ErrInvalidTransition", func(t *testing.T) {
		p, _ := NewProduct("Title", "Desc", CatId(1), false, nil)
		_ = p.Publish()
		_ = p.Archive()
		err := p.Publish()
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got %v", err)
		}
	})
}

func TestProduct_Archive(t *testing.T) {
	t.Run("active product/should transition to archived", func(t *testing.T) {
		p, _ := NewProduct("Title", "Desc", CatId(1), false, nil)
		_ = p.Publish()
		err := p.Archive()
		if err != nil {
			t.Fatalf("unexpected error on archive: %v", err)
		}
		if p.Status() != StatusArchived {
			t.Errorf("expected status %s, got %s", StatusArchived, p.Status())
		}
	})

	t.Run("draft product/should return ErrInvalidTransition", func(t *testing.T) {
		p, _ := NewProduct("Title", "Desc", CatId(1), false, nil)
		err := p.Archive()
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got %v", err)
		}
	})
}

func TestProduct_Update(t *testing.T) {
	t.Run("update title/should set title and refresh updatedAt", func(t *testing.T) {
		p, _ := NewProduct("Old Title", "Desc", CatId(1), false, nil)
		oldUpdatedAt := p.UpdatedAt()
		time.Sleep(10 * time.Millisecond)

		newTitle := "New Title"
		p.Update(UpdateProductCmd{Title: &newTitle})

		if p.Title() != "New Title" {
			t.Errorf("expected title 'New Title', got %s", p.Title())
		}
		if !p.UpdatedAt().After(oldUpdatedAt) {
			t.Errorf("expected updatedAt %v to be after %v", p.UpdatedAt(), oldUpdatedAt)
		}
	})

	t.Run("nil fields/should not change", func(t *testing.T) {
		p, _ := NewProduct("Title", "Desc", CatId(1), false, []string{"tag1"})
		oldUpdatedAt := p.UpdatedAt()

		p.Update(UpdateProductCmd{})

		if p.Title() != "Title" {
			t.Errorf("expected title to remain 'Title', got %s", p.Title())
		}
		if p.Description() != "Desc" {
			t.Errorf("expected description to remain 'Desc', got %s", p.Description())
		}
		if p.CategoryId() != CatId(1) {
			t.Errorf("expected categoryId to remain 1, got %d", p.CategoryId())
		}
		if !p.UpdatedAt().Equal(oldUpdatedAt) {
			t.Errorf("expected updatedAt not to change on no-op update, got %v vs %v", p.UpdatedAt(), oldUpdatedAt)
		}
	})
}

func TestProduct_SetSlug(t *testing.T) {
	t.Run("should set slug", func(t *testing.T) {
		p, _ := NewProduct("Title", "Desc", CatId(1), false, nil)
		p.SetSlug("my-slug")
		if p.Slug() != "my-slug" {
			t.Errorf("expected slug 'my-slug', got %s", p.Slug())
		}
	})
}

func TestNewProductVariant(t *testing.T) {
	t.Run("valid input/should create variant", func(t *testing.T) {
		attrs := NewAttributes()
		attrs.AddStr("Color", "Blue")
		v, err := NewProductVariant(ProId(10), "SKU-001", 1999, attrs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.ProductId() != ProId(10) {
			t.Errorf("expected productId 10, got %d", v.ProductId())
		}
		if v.SKU() != "SKU-001" {
			t.Errorf("expected SKU 'SKU-001', got %s", v.SKU())
		}
		if v.Price() != 1999 {
			t.Errorf("expected price 1999, got %d", v.Price())
		}
		if v.StockQty() != 0 {
			t.Errorf("expected stockQty 0, got %d", v.StockQty())
		}
		if v.Attrs() != attrs {
			t.Errorf("expected attrs pointer to match")
		}
		if v.CompareAtPrice() != nil {
			t.Errorf("expected nil compareAtPrice, got %v", v.CompareAtPrice())
		}
		if v.MediaId() != nil {
			t.Errorf("expected nil mediaId, got %v", v.MediaId())
		}
	})

	t.Run("empty sku/should return error", func(t *testing.T) {
		_, err := NewProductVariant(ProId(10), "", 1999, nil)
		if err == nil {
			t.Error("expected error for empty sku, got nil")
		}
		_, err2 := NewProductVariant(ProId(10), "   ", 1999, nil)
		if err2 == nil {
			t.Error("expected error for whitespace sku, got nil")
		}
	})

	t.Run("negative price/should return error", func(t *testing.T) {
		_, err := NewProductVariant(ProId(10), "SKU-001", -1, nil)
		if err == nil {
			t.Error("expected error for negative price, got nil")
		}
	})
}

func TestProductVariant_AdjustStock(t *testing.T) {
	t.Run("positive delta/should increase stock", func(t *testing.T) {
		v, _ := NewProductVariant(ProId(1), "SKU-1", 100, nil)
		err := v.AdjustStock(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.StockQty() != 10 {
			t.Errorf("expected stockQty 10, got %d", v.StockQty())
		}
	})

	t.Run("negative delta within range/should decrease", func(t *testing.T) {
		v, _ := NewProductVariant(ProId(1), "SKU-1", 100, nil)
		_ = v.AdjustStock(10)
		err := v.AdjustStock(-5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.StockQty() != 5 {
			t.Errorf("expected stockQty 5, got %d", v.StockQty())
		}
	})

	t.Run("negative delta below zero/should return ErrInsufficientStock", func(t *testing.T) {
		v, _ := NewProductVariant(ProId(1), "SKU-1", 100, nil)
		_ = v.AdjustStock(3)
		err := v.AdjustStock(-5)
		if !errors.Is(err, ErrInsufficientStock) {
			t.Errorf("expected ErrInsufficientStock, got %v", err)
		}
		if v.StockQty() != 3 {
			t.Errorf("expected stockQty to remain 3, got %d", v.StockQty())
		}
	})
}

func TestProductVariant_SetPrice(t *testing.T) {
	t.Run("changed price/should return old and true", func(t *testing.T) {
		v, _ := NewProductVariant(ProId(1), "SKU-1", 100, nil)
		old, changed := v.SetPrice(200)
		if old != 100 {
			t.Errorf("expected old price 100, got %d", old)
		}
		if !changed {
			t.Error("expected changed to be true")
		}
		if v.Price() != 200 {
			t.Errorf("expected price 200, got %d", v.Price())
		}
	})

	t.Run("same price/should return old and false", func(t *testing.T) {
		v, _ := NewProductVariant(ProId(1), "SKU-1", 100, nil)
		old, changed := v.SetPrice(100)
		if old != 100 {
			t.Errorf("expected old price 100, got %d", old)
		}
		if changed {
			t.Error("expected changed to be false")
		}
		if v.Price() != 100 {
			t.Errorf("expected price 100, got %d", v.Price())
		}
	})
}

func TestNewCategory(t *testing.T) {
	t.Run("valid input/should create category", func(t *testing.T) {
		c, err := NewCategory("Electronics", "electronics", "", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Name() != "Electronics" {
			t.Errorf("expected name 'Electronics', got %s", c.Name())
		}
		if c.Slug() != "electronics" {
			t.Errorf("expected slug 'electronics', got %s", c.Slug())
		}
		if c.ParentId() != nil {
			t.Errorf("expected nil parentId for root, got %v", c.ParentId())
		}
		if c.Depth() != 0 {
			t.Errorf("expected depth 0 for root, got %d", c.Depth())
		}
		if c.Attrs() == nil || len(c.Attrs().Names()) != 0 {
			t.Errorf("expected empty attrs")
		}
	})

	t.Run("child category/should compute depth from parent path", func(t *testing.T) {
		parentId := CatId(5)
		c, err := NewCategory("Phones", "phones", "1.2", &parentId)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Depth() != 2 {
			t.Errorf("expected depth 2 for parentPath '1.2', got %d", c.Depth())
		}
		if c.ParentId() == nil || *c.ParentId() != CatId(5) {
			t.Errorf("expected parentId 5, got %v", c.ParentId())
		}
	})

	t.Run("empty name/should return error", func(t *testing.T) {
		_, err := NewCategory("", "slug", "", nil)
		if err == nil {
			t.Error("expected error for empty name, got nil")
		}
		_, err2 := NewCategory("   ", "slug", "", nil)
		if err2 == nil {
			t.Error("expected error for whitespace name, got nil")
		}
	})
}

func TestCategory_SetPath(t *testing.T) {
	t.Run("should set path", func(t *testing.T) {
		c, _ := NewCategory("Name", "slug", "", nil)
		c.SetPath("1.2.3")
		if c.Path() != "1.2.3" {
			t.Errorf("expected path '1.2.3', got %s", c.Path())
		}
	})
}

func TestCategory_SetSortOrder(t *testing.T) {
	t.Run("should set sort order", func(t *testing.T) {
		c, _ := NewCategory("Name", "slug", "", nil)
		c.SetSortOrder(5)
		if c.SortOrder() != 5 {
			t.Errorf("expected sortOrder 5, got %d", c.SortOrder())
		}
	})
}

func TestBuildCatTree(t *testing.T) {
	t.Run("flat list/should build correct hierarchy with parent links", func(t *testing.T) {
		rootCat := Category{id: CatId(1), name: "Root", path: "1"}
		rootId := CatId(1)
		child1 := Category{id: CatId(2), name: "Child 1", parentId: &rootId, path: "1.2"}
		child2 := Category{id: CatId(3), name: "Child 2", parentId: &rootId, path: "1.3"}

		cats := []Category{rootCat, child1, child2}
		roots := BuildCatTree(cats)

		if len(roots) != 1 {
			t.Fatalf("expected 1 root node, got %d", len(roots))
		}
		rootNode := roots[0]
		if rootNode.Cat.Id() != CatId(1) {
			t.Errorf("expected root id 1, got %d", rootNode.Cat.Id())
		}
		if rootNode.Parent != nil {
			t.Errorf("expected root parent to be nil, got %v", rootNode.Parent)
		}
		if len(rootNode.Children) != 2 {
			t.Fatalf("expected 2 children, got %d", len(rootNode.Children))
		}
		for _, childNode := range rootNode.Children {
			if childNode.Parent != rootNode {
				t.Errorf("expected child's Parent to point to rootNode")
			}
		}
	})

	t.Run("single root/should have nil parent", func(t *testing.T) {
		rootCat := Category{id: CatId(1), name: "Root"}
		roots := BuildCatTree([]Category{rootCat})
		if len(roots) != 1 {
			t.Fatalf("expected 1 root node, got %d", len(roots))
		}
		if roots[0].Parent != nil {
			t.Errorf("expected root parent to be nil, got %v", roots[0].Parent)
		}
	})
}

func TestNewProductMedia(t *testing.T) {
	t.Run("should create media with all fields", func(t *testing.T) {
		m := NewProductMedia(ProId(1), "images/product1.jpg", "Photo of product", 2)
		if m.ProductId() != ProId(1) {
			t.Errorf("expected productId 1, got %d", m.ProductId())
		}
		if m.URI() != "images/product1.jpg" {
			t.Errorf("expected uri 'images/product1.jpg', got %s", m.URI())
		}
		if m.AltText() != "Photo of product" {
			t.Errorf("expected altText 'Photo of product', got %s", m.AltText())
		}
		if m.Order() != 2 {
			t.Errorf("expected order 2, got %d", m.Order())
		}
	})
}


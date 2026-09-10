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
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"octacart/internal/shared"
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

// Product represents the central catalog item aggregate root.
type Product struct {
	id          ProId
	title       string
	slug        string
	description string
	status      ProductStatus
	categoryId  CatId
	tags        []string
	isDigital   bool
	createdAt   time.Time
	updatedAt   time.Time
}

// Id returns the product ID.
func (p *Product) Id() ProId { return p.id }

// Title returns the product title.
func (p *Product) Title() string { return p.title }

// Slug returns the product URL slug.
func (p *Product) Slug() string { return p.slug }

// Description returns the product description.
func (p *Product) Description() string { return p.description }

// Status returns the product lifecycle status.
func (p *Product) Status() ProductStatus { return p.status }

// CategoryId returns the category ID this product belongs to.
func (p *Product) CategoryId() CatId { return p.categoryId }

// Tags returns a copy of product tags.
func (p *Product) Tags() []string {
	if p.tags == nil {
		return nil
	}
	tagsCopy := make([]string, len(p.tags))
	copy(tagsCopy, p.tags)
	return tagsCopy
}

// IsDigital returns true if the product is digital.
func (p *Product) IsDigital() bool { return p.isDigital }

// CreatedAt returns when the product was created.
func (p *Product) CreatedAt() time.Time { return p.createdAt }

// UpdatedAt returns when the product was last updated.
func (p *Product) UpdatedAt() time.Time { return p.updatedAt }

// NewProduct creates a new Product aggregate in Draft status.
func NewProduct(title, description string, categoryId CatId, isDigital bool, tags []string) (*Product, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	now := time.Now()
	var tagsCopy []string
	if tags != nil {
		tagsCopy = make([]string, len(tags))
		copy(tagsCopy, tags)
	}
	return &Product{
		title:       title,
		description: description,
		categoryId:  categoryId,
		status:      StatusDraft,
		tags:        tagsCopy,
		isDigital:   isDigital,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Publish transitions product from Draft to Active.
func (p *Product) Publish() error {
	if p.status != StatusDraft {
		return ErrInvalidTransition
	}
	p.status = StatusActive
	p.updatedAt = time.Now()
	return nil
}

// Archive transitions product from Active to Archived.
func (p *Product) Archive() error {
	if p.status != StatusActive {
		return ErrInvalidTransition
	}
	p.status = StatusArchived
	p.updatedAt = time.Now()
	return nil
}

// SetSlug sets the product URL slug.
func (p *Product) SetSlug(slug string) {
	p.slug = slug
	p.updatedAt = time.Now()
}

// UpdateProductCmd contains mutable fields for Product updates.
type UpdateProductCmd struct {
	Title       *string
	Description *string
	CategoryId  *CatId
	Tags        *[]string
	IsDigital   *bool
}

// Update updates mutable fields on the product.
func (p *Product) Update(cmd UpdateProductCmd) {
	modified := false
	if cmd.Title != nil {
		p.title = *cmd.Title
		modified = true
	}
	if cmd.Description != nil {
		p.description = *cmd.Description
		modified = true
	}
	if cmd.CategoryId != nil {
		p.categoryId = *cmd.CategoryId
		modified = true
	}
	if cmd.Tags != nil {
		if *cmd.Tags == nil {
			p.tags = nil
		} else {
			tagsCopy := make([]string, len(*cmd.Tags))
			copy(tagsCopy, *cmd.Tags)
			p.tags = tagsCopy
		}
		modified = true
	}
	if cmd.IsDigital != nil {
		p.isDigital = *cmd.IsDigital
		modified = true
	}
	if modified {
		p.updatedAt = time.Now()
	}
}

// ProductVariant represents a purchasable SKU variant of a product.
type ProductVariant struct {
	id             ProVariantId
	productId      ProId
	sku            string
	price          shared.Money
	compareAtPrice *shared.Money
	stockQty       int
	attrs          *Attributes
	mediaId        *MediaId
}

// Id returns the variant ID.
func (v *ProductVariant) Id() ProVariantId { return v.id }

// ProductId returns the owning product ID.
func (v *ProductVariant) ProductId() ProId { return v.productId }

// SKU returns the SKU code.
func (v *ProductVariant) SKU() string { return v.sku }

// Price returns the current price in least minor units.
func (v *ProductVariant) Price() shared.Money { return v.price }

// CompareAtPrice returns the optional crossed-out price.
func (v *ProductVariant) CompareAtPrice() *shared.Money {
	if v.compareAtPrice == nil {
		return nil
	}
	val := *v.compareAtPrice
	return &val
}

// StockQty returns current on-hand stock quantity.
func (v *ProductVariant) StockQty() int { return v.stockQty }

// Attrs returns the dynamic attributes of the variant.
func (v *ProductVariant) Attrs() *Attributes { return v.attrs }

// MediaId returns the optional linked media asset ID.
func (v *ProductVariant) MediaId() *MediaId {
	if v.mediaId == nil {
		return nil
	}
	val := *v.mediaId
	return &val
}

// NewProductVariant creates a new variant for a product.
func NewProductVariant(productId ProId, sku string, price shared.Money, attrs *Attributes) (*ProductVariant, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, errors.New("sku cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	return &ProductVariant{
		productId: productId,
		sku:       sku,
		price:     price,
		stockQty:  0,
		attrs:     attrs,
	}, nil
}

// AdjustStock adjusts the variant stock quantity by delta.
func (v *ProductVariant) AdjustStock(delta int) error {
	newStock := v.stockQty + delta
	if newStock < 0 {
		return ErrInsufficientStock
	}
	v.stockQty = newStock
	return nil
}

// UpdateVariantCmd contains mutable fields for ProductVariant updates.
type UpdateVariantCmd struct {
	SKU            *string
	Price          *shared.Money
	CompareAtPrice **shared.Money
	Attrs          *Attributes
	MediaId        **MediaId
}

// Update updates mutable fields on the variant.
func (v *ProductVariant) Update(cmd UpdateVariantCmd) {
	if cmd.SKU != nil {
		v.sku = *cmd.SKU
	}
	if cmd.Price != nil {
		v.price = *cmd.Price
	}
	if cmd.CompareAtPrice != nil {
		if *cmd.CompareAtPrice == nil {
			v.compareAtPrice = nil
		} else {
			val := **cmd.CompareAtPrice
			v.compareAtPrice = &val
		}
	}
	if cmd.Attrs != nil {
		v.attrs = cmd.Attrs
	}
	if cmd.MediaId != nil {
		if *cmd.MediaId == nil {
			v.mediaId = nil
		} else {
			val := **cmd.MediaId
			v.mediaId = &val
		}
	}
}

// SetPrice updates the price and returns the old price and whether it changed.
func (v *ProductVariant) SetPrice(newPrice shared.Money) (oldPrice shared.Money, changed bool) {
	oldPrice = v.price
	if oldPrice != newPrice {
		v.price = newPrice
		return oldPrice, true
	}
	return oldPrice, false
}

// Category represents a hierarchical category node.
type Category struct {
	id        CatId
	name      string
	slug      string
	parentId  *CatId
	path      string
	depth     int
	sortOrder int
	attrs     *Attributes
}

// Id returns the category ID.
func (c *Category) Id() CatId { return c.id }

// Name returns the category display name.
func (c *Category) Name() string { return c.name }

// Slug returns the category URL slug.
func (c *Category) Slug() string { return c.slug }

// ParentId returns the parent category ID or nil if root.
func (c *Category) ParentId() *CatId {
	if c.parentId == nil {
		return nil
	}
	val := *c.parentId
	return &val
}

// Path returns the materialized dot-separated path.
func (c *Category) Path() string { return c.path }

// Depth returns the tree depth level (0 for root).
func (c *Category) Depth() int { return c.depth }

// SortOrder returns the sibling sorting order.
func (c *Category) SortOrder() int { return c.sortOrder }

// Attrs returns the dynamic attribute specifications of the category.
func (c *Category) Attrs() *Attributes { return c.attrs }

// NewCategory creates a new category.
func NewCategory(name, slug string, parentPath string, parentId *CatId) (*Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}
	depth := 0
	if parentPath != "" {
		depth = strings.Count(parentPath, ".") + 1
	}
	var pId *CatId
	if parentId != nil {
		val := *parentId
		pId = &val
	}
	return &Category{
		name:      name,
		slug:      slug,
		parentId:  pId,
		path:      "",
		depth:     depth,
		sortOrder: 0,
		attrs:     NewAttributes(),
	}, nil
}

// UpdateCatCmd contains mutable fields for Category updates.
type UpdateCatCmd struct {
	Name  *string
	Slug  *string
	Attrs *Attributes
}

// Update updates mutable fields on the category.
func (c *Category) Update(cmd UpdateCatCmd) {
	if cmd.Name != nil {
		c.name = *cmd.Name
	}
	if cmd.Slug != nil {
		c.slug = *cmd.Slug
	}
	if cmd.Attrs != nil {
		c.attrs = cmd.Attrs
	}
}

// SetPath sets the materialized path after DB assigns the ID.
func (c *Category) SetPath(path string) {
	c.path = path
}

// SetSortOrder sets the sort order.
func (c *Category) SetSortOrder(order int) {
	c.sortOrder = order
}

// CatNode represents a category tree node.
type CatNode struct {
	Cat      Category
	Parent   *CatNode
	Children []*CatNode
}

// BuildCatTree converts a flat list of categories into a tree of CatNodes.
func BuildCatTree(cats []Category) []*CatNode {
	cache := make(map[CatId]*CatNode, len(cats))
	nodes := make([]*CatNode, len(cats))

	for i, c := range cats {
		node := &CatNode{
			Cat:      c,
			Children: make([]*CatNode, 0),
		}
		nodes[i] = node
		cache[c.id] = node
	}

	var roots []*CatNode
	for _, node := range nodes {
		if node.Cat.parentId == nil {
			node.Parent = nil
			roots = append(roots, node)
		} else if parentNode, exists := cache[*node.Cat.parentId]; exists {
			node.Parent = parentNode
			parentNode.Children = append(parentNode.Children, node)
		} else {
			node.Parent = nil
			roots = append(roots, node)
		}
	}

	return roots
}

// ProductMedia represents an ordered media asset gallery item.
type ProductMedia struct {
	id        MediaId
	productId ProId
	uri       string
	altText   string
	order     int
}

// Id returns the media ID.
func (m *ProductMedia) Id() MediaId { return m.id }

// ProductId returns the product ID this media belongs to.
func (m *ProductMedia) ProductId() ProId { return m.productId }

// URI returns the relative path or direct URL of the media.
func (m *ProductMedia) URI() string { return m.uri }

// AltText returns accessibility/SEO alt text.
func (m *ProductMedia) AltText() string { return m.altText }

// Order returns the display order position.
func (m *ProductMedia) Order() int { return m.order }

// NewProductMedia creates a new ProductMedia item.
func NewProductMedia(productId ProId, uri, altText string, order int) *ProductMedia {
	return &ProductMedia{
		productId: productId,
		uri:       uri,
		altText:   altText,
		order:     order,
	}
}


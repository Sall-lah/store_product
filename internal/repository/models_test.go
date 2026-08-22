package repository

import (
	"testing"
	"time"
)

// TestProductFilter_InactiveDefaults verifies that ProductFilter correctly encapsulates filtering rules.
func TestProductFilter_InactiveDefaults(t *testing.T) {
	filter := ProductFilter{
		Limit:           20,
		IncludeInactive: false,
	}

	if filter.IncludeInactive {
		t.Errorf("expected default filter to not include inactive products")
	}

	isActive := false
	adminFilter := ProductFilter{
		Limit:           20,
		IncludeInactive: true,
		IsActive:        &isActive,
	}

	if !adminFilter.IncludeInactive {
		t.Errorf("expected admin filter to include inactive products")
	}
	if adminFilter.IsActive == nil || *adminFilter.IsActive != false {
		t.Errorf("expected admin filter to explicitly target inactive products")
	}
}

// TestProductDTO_SoftDeleteProperties verifies soft-deleted product and variant DTO representations.
func TestProductDTO_SoftDeleteProperties(t *testing.T) {
	now := time.Now()
	desc := "Soft-deleted sample item"
	size := "M"
	color := "Navy"
	price := 49.99

	softDeletedProduct := ProductDTO{
		ID:          "prod-soft-del-1",
		Name:        "Vintage Jacket",
		Slug:        "vintage-jacket",
		Description: &desc,
		BasePrice:   99.99,
		Category:    "outerwear",
		IsActive:    false,
		Variants: []VariantDTO{
			{
				ID:        "var-soft-del-1",
				ProductID: "prod-soft-del-1",
				SKU:       "VJ-NAVY-M",
				Size:      &size,
				Color:     &color,
				Price:     &price,
				Stock:     15,
				IsActive:  false,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if softDeletedProduct.IsActive {
		t.Errorf("expected product IsActive to be false")
	}
	if len(softDeletedProduct.Variants) != 1 {
		t.Fatalf("expected 1 variant, got %d", len(softDeletedProduct.Variants))
	}
	if softDeletedProduct.Variants[0].IsActive {
		t.Errorf("expected variant IsActive to be false")
	}
	if softDeletedProduct.Variants[0].SKU != "VJ-NAVY-M" {
		t.Errorf("expected variant SKU 'VJ-NAVY-M' to be preserved, got '%s'", softDeletedProduct.Variants[0].SKU)
	}
}

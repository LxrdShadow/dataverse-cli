package commands

import (
	"testing"

	"dvc/models"
)

func TestParseTableFiltersAll(t *testing.T) {
	filters, err := parseTableFilters([]string{})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.EntityScopeAll {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.EntityScopeAll)
	}
	if filters.Management != models.EntityManagementAll {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.EntityManagementAll)
	}
}

func TestParseTableFilters(t *testing.T) {
	filters, err := parseTableFilters([]string{"--custom"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.EntityScopeCustom {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.EntityScopeCustom)
	}
	if filters.Management != models.EntityManagementAll {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.EntityManagementAll)
	}
}

func TestParseTableFiltersUnmanaged(t *testing.T) {
	filters, err := parseTableFilters([]string{"--unmanaged"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Management != models.EntityManagementUnmanaged {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.EntityManagementUnmanaged)
	}
}

func TestParseTableFiltersSystemManaged(t *testing.T) {
	filters, err := parseTableFilters([]string{"--system", "--managed"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.EntityScopeSystem {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.EntityScopeSystem)
	}
	if filters.Management != models.EntityManagementManaged {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.EntityManagementManaged)
	}
}

func TestParseTableFiltersUnmanagedCustom(t *testing.T) {
	filters, err := parseTableFilters([]string{"--unmanaged", "--custom"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.EntityScopeCustom {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.EntityScopeCustom)
	}
	if filters.Management != models.EntityManagementUnmanaged {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.EntityManagementUnmanaged)
	}
}

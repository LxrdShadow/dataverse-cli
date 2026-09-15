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
	if filters.Scope != models.TableScopeAll {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.TableScopeAll)
	}
	if filters.Management != models.TableManagementAll {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.TableManagementAll)
	}
}

func TestParseTableFilters(t *testing.T) {
	filters, err := parseTableFilters([]string{"--custom"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.TableScopeCustom {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.TableScopeCustom)
	}
	if filters.Management != models.TableManagementAll {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.TableManagementAll)
	}
}

func TestParseTableFiltersUnmanaged(t *testing.T) {
	filters, err := parseTableFilters([]string{"--unmanaged"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Management != models.TableManagementUnmanaged {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.TableManagementUnmanaged)
	}
}

func TestParseTableFiltersSystemManaged(t *testing.T) {
	filters, err := parseTableFilters([]string{"--system", "--managed"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.TableScopeSystem {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.TableScopeSystem)
	}
	if filters.Management != models.TableManagementManaged {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.TableManagementManaged)
	}
}

func TestParseTableFiltersUnmanagedCustom(t *testing.T) {
	filters, err := parseTableFilters([]string{"--unmanaged", "--custom"})
	if err != nil {
		t.Errorf("parseTableFilters() error = %v", err)
	}
	if filters.Scope != models.TableScopeCustom {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Scope, models.TableScopeCustom)
	}
	if filters.Management != models.TableManagementUnmanaged {
		t.Errorf("parseTableFilters() = %v, want %v", filters.Management, models.TableManagementUnmanaged)
	}
}

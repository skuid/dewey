package registry

import (
	"encoding/json"
	"fmt"
	"path"
)

// RegistryRepo is a common struct for the format of returns from the registry
type RegistryRepo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// RegistryCatalog is a container for the catalog
type RegistryCatalog struct {
	Repositories []string `json:"repositories"`
}

// FileContent will return the output of the file in either pretty or non-pretty format
func (r *RegistryCatalog) FileContent(pretty bool) ([]byte, error) {
	var content []byte
	var err error
	if pretty {
		content, err = json.MarshalIndent(r, "", "    ")
	} else {
		content, err = json.Marshal(r)
	}

	return content, err
}

// RepoConfig is a container for all repo information
type RepoConfig struct {
	Name           string
	Address        string
	Kind           string
	Username       string
	Password       string
	Orgs           []string
	Repositories   []string
	OutputFilename string `json:"outputFile"`
}

// AddressOrDefault returns the predefined url for each registry implementation
// or the overridden value
func (r RepoConfig) AddressOrDefault(predefined string) string {
	if r.Address != "" {
		return r.Address
	}
	return predefined
}

// Filename returns the filename default of the overridden filenames
func (r RepoConfig) Filename(base string) string {
	if r.OutputFilename != "" {
		return r.OutputFilename
	}

	return path.Join(base, fmt.Sprintf("%s.json", r.Name))
}

// CatalogableRegistry defines an interface for a registry which can be shimmed
type CatalogableRegistry interface {
	GetCatalog() (*RegistryCatalog, error)
}

// ConvertToCatalogableRegistry converts RepoConfig into a CatalogableRegistry implementation
func ConvertToCatalogableRegistry(r RepoConfig) (CatalogableRegistry, error) {
	var catalog CatalogableRegistry
	var err error
	switch r.Kind {
	case "quay":
		catalog = QuayRegistry{r}
	case "dockerhub":
		catalog = DockerhubRegistry{r}
	default:
		err = fmt.Errorf("%s is an invalid registry kind", r.Kind)
	}
	return catalog, err
}

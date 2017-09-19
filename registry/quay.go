package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const quayApiBase = "https://quay.io/api/v1"

type QuayRegistry struct {
	RepoConfig
}

// GetCatalog returns a list of repositories based on the desired list of orgs and static repositories
func (r QuayRegistry) GetCatalog() (*RegistryCatalog, error) {
	repos := r.RepoConfig.Repositories
	baseURL := r.RepoConfig.AddressOrDefault(quayApiBase)
	for _, org := range r.RepoConfig.Orgs {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/repository?private=true&public=true&namespace=%s", baseURL, org), nil)
		if err != nil {
			return nil, err
		}

		if r.RepoConfig.Password != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.RepoConfig.Password))
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		type QuayResp struct {
			Repositories []RegistryRepo `json:"repositories"`
		}
		rep := &QuayResp{}

		if err := json.NewDecoder(resp.Body).Decode(rep); err != nil {
			return nil, err
		}

		for _, repo := range rep.Repositories {
			repos = append(repos, fmt.Sprintf("%s/%s", repo.Namespace, repo.Name))
		}
	}

	return &RegistryCatalog{Repositories: repos}, nil

}

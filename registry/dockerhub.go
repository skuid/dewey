package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DockerhubRegistry struct {
	RepoConfig
}

type dockerhubAuthToken struct {
	Token string `json"token"`
}

type dockerhubResponse struct {
	Count   int            `json:"count"`
	Next    string         `json:"next"`
	Results []RegistryRepo `json"results"`
}

var dockerhubApiBase = "https://hub.docker.com/v2"

// curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${UNAME}'", "password": "'${UPASS}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token
func (r DockerhubRegistry) getAuthToken() (string, error) {
	baseURL := r.RepoConfig.AddressOrDefault(dockerhubApiBase)
	payload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, r.RepoConfig.Username, r.RepoConfig.Password)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users/login", baseURL), bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	token := &dockerhubAuthToken{}
	if err = json.NewDecoder(resp.Body).Decode(token); err != nil {
		return "", err
	}

	if token.Token == "" {
		return "", fmt.Errorf("failed to log into the registry")
	}

	return token.Token, nil
}

func getRepositoryPage(url, token string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("JWT %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parsed := &dockerhubResponse{}
	if err := json.NewDecoder(resp.Body).Decode(parsed); err != nil {
		return nil, err
	}

	repos := make([]string, 0)
	if parsed.Next != "" {
		r, err := getRepositoryPage(parsed.Next, token)
		if err != nil {
			return nil, err
		}
		// append the next page of repos to the current list
		// that will be added from parsed.Results. It's recursive.
		repos = append(repos, r...)
	}

	for _, repo := range parsed.Results {
		repos = append(repos, fmt.Sprintf("%s/%s", repo.Namespace, repo.Name))
	}

	return repos, nil

}

// GetCatalog returns a list of repositories based on the desired list of orgs and static repositories
func (r DockerhubRegistry) GetCatalog() (*RegistryCatalog, error) {
	repos := r.RepoConfig.Repositories
	baseURL := r.RepoConfig.AddressOrDefault(dockerhubApiBase)
	token := ""
	if r.RepoConfig.Password != "" && r.RepoConfig.Username != "" {
		tok, err := r.getAuthToken()
		if err != nil {
			return nil, err
		}
		token = tok
	}
	for _, org := range r.RepoConfig.Orgs {
		catalog, err := getRepositoryPage(fmt.Sprintf("%s/repositories/%s?page=1&page_size=2", baseURL, org), token)
		if err != nil {
			return nil, err
		}
		repos = append(repos, catalog...)
	}

	return &RegistryCatalog{Repositories: repos}, nil
}

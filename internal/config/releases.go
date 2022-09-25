package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
)

// Release is used to save the release informations
type Release struct {
	ID      int             `json:"id,omitempty"`
	TagName string          `json:"tag_name,omitempty"`
	Name    adapter.Version `json:"name,omitempty"`
	Draft   bool            `json:"draft,omitempty"`
	Assets  []*Asset        `json:"assets,omitempty"`
}

// Asset describes the github release asset object
type Asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

// getLatestReleaseNames returns the names of the latest releases
// limited by the "limit" parameter. The first version in the list
// is always is the latest "stable" version.
func GetLatestReleaseNames(limit int) ([]adapter.Version, error) {
	releases, err := getLatestReleases(uint(limit))
	if err != nil {
		return []adapter.Version{}, ErrGetLatestReleaseNames(err)
	}

	var releaseNames []adapter.Version
	var latestStable adapter.Version = ""

	for _, r := range releases {
		releaseNames = append(releaseNames, r.Name)
		if latestStable == "" && strings.HasPrefix(string(r.Name), "stable") {
			latestStable = r.Name
		}
	}

	// Ensure that limit is always lesser than equal to the total
	// number of releases
	if limit > len(releaseNames) {
		limit = len(releaseNames)
	}

	result := make([]adapter.Version, limit)

	// Make latest stable as the first name
	result[0] = latestStable
	// Add other elements to the list
	for i := 1; i < limit; i++ {
		if releaseNames[i-1] != latestStable {
			result[i] = releaseNames[i-1]
		}
	}

	return result, nil
}

// GetLatestReleases fetches the latest releases from the linkerd repository
func getLatestReleases(releases uint) ([]*Release, error) {
	releaseAPIURL := "https://api.github.com/repos/linkerd/linkerd2/releases?per_page=" + fmt.Sprint(releases)
	// We need a variable url here hence using nosec
	// #nosec
	resp, err := http.Get(releaseAPIURL)
	if err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	if resp.StatusCode != http.StatusOK {
		return []*Release{}, ErrGetLatestReleases(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	var releaseList []*Release

	if err = json.Unmarshal(body, &releaseList); err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	if err = resp.Body.Close(); err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	return releaseList, nil
}

// Copyright 2020 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
func getLatestReleaseNames(limit int) ([]adapter.Version, error) {
	releases, err := GetLatestReleases(30)
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
func GetLatestReleases(releases uint) ([]*Release, error) {
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

	body, err := ioutil.ReadAll(resp.Body)
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

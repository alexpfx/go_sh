package merge

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

type Config struct {
	BaseUrl           string `json:"base_url"`
	ApiVersion        string `json:"api_version"`
	Project           string `json:"project"`
	MergeRequestsPath string `json:"mergeRequestsPath"`
	MergeRequestQuery string `json:"mergeRequestQuery"`
	CommitsPath       string `json:"commitsPath"`
}

type MrFetch struct {
	Config Config
}

func (m MrFetch) Fetch(mergeId string) {
	createClient()

	config := m.Config
	url := createUrl(config.BaseUrl, config.Project, config.MergeRequestsPath, config.MergeRequestQuery+mergeId)

	fmt.Println(url)


}
func createUrl(base, project, path, query string) string {
	return strings.Join([]string{base, project, path, query}, "/")
}

func createClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}
package merge

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/common/util"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Merge struct {
	Iid            int    `json:"iid"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	State          string `json:"state"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	TargetBranch   string `json:"target_branch"`
	SourceBranch   string `json:"source_branch"`
	WebUrl         string `json:"web_url"`
	Author         User   `json:"author"`
	MergeCommitSha string `json:"merge_commit_sha"`
	Commit         Commit `json:"commit"`
}

type User struct {
	Username string `json:"username"`
}

type Commit struct {
	Id        string `json:"id"`
	Email     string `json:"author_email"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
}

type FetchResult struct {
	Merge Merge `json:"merge"`
}

type Config struct {
	BaseUrl           string `json:"base_url"`
	MergeRequestsPath string `json:"merge_requests_path"`
	MergeRequestQuery string `json:"merge_request_query"`
	CommitsPath       string `json:"commits_path"`
	Project           string `json:"project"`
}

type Fetcher struct {
	Config Config
	Token  string
}

func (m Fetcher) Fetch(mergeId string) []FetchResult {
	client := createClient()
	config := m.Config
	url := createUrl(config.BaseUrl, config.Project, config.MergeRequestsPath, config.MergeRequestQuery+mergeId)
	req := createRequest(url, m.Token)
	resp, err := client.Do(req)
	util.CheckFatal(err, "")

	body, err := ioutil.ReadAll(resp.Body)
	util.CheckFatal(err, "")

	var merges []Merge
	json.Unmarshal(body, &merges)
	util.CheckFatal(err, "")
	if len(merges) < 1 {
		fmt.Printf("Merge Request não encontrado: %s", mergeId)
	}

	var results []FetchResult

	for _, merge := range merges {
		commit, err := m.fetchCommit(merge.MergeCommitSha, m.Token)
		util.CheckFatal(err, "")
		results = appendResult(results, merge, commit)
	}
	return results
}

func (f Fetcher) addOrDiscard(merge Merge, mrList []FetchResult, filter map[string]string, token string) ([]FetchResult, error) {

	if filter == nil || len(filter) < 1 {
		commit, err := f.fetchCommit(merge.MergeCommitSha, token)
		if err != nil {
			return mrList, err
		}
		return appendResult(mrList, merge, commit), nil
	}

	for k, v := range filter {

		if strings.EqualFold(k, "author") {
			if merge.Author.Username != v {
				return mrList, nil
			}
		}
		if strings.EqualFold(k, "target_branch") {
			if !strings.EqualFold(v, merge.TargetBranch) {
				return mrList, nil
			}
		}
	}

	if merge.MergeCommitSha == "" {
		return mrList, nil
	} // filtros

	commit, err := f.fetchCommit(merge.MergeCommitSha, token)
	if err != nil {
		return mrList, err
	}
	return appendResult(mrList, merge, commit), nil
}

func (f Fetcher) fetchCommit(commitSha, token string) (Commit, error) {
	url := createUrl(f.Config.BaseUrl, f.Config.Project, f.Config.CommitsPath, commitSha)
	client := createClient()
	var commit Commit

	req := createRequest(url, token)

	r, e := client.Do(req)
	if e != nil {
		return commit, e
	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return commit, e
	}

	e = json.Unmarshal(body, &commit)
	if e != nil {
		return commit, e
	}

	createdAt, _ := time.Parse(time.RFC3339, commit.CreatedAt)

	return Commit{
		Id:        commitSha,
		CreatedAt: createdAt.Format("2006-01-02T15:04:05Z"),
		Email:     commit.Email,
		Username: strings.FieldsFunc(commit.Email, func(r rune) bool {
			return r == '@'
		})[0],
	}, e
}

func appendResult(results []FetchResult, merge Merge, commit Commit) []FetchResult {
	merge.Commit = commit
	return append(results, FetchResult{
		Merge: merge,
	})
}

func createUrl(base, project, path, query string) string {
	return strings.Join([]string{base, project, path, query}, "/")
}
func createRequest(url string, token string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("PRIVATE-TOKEN", token)
	return req
}

func createClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}

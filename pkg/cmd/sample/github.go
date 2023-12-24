package sample

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lab5e/spancli/pkg/commonopt"
	"github.com/lab5e/spancli/pkg/helpers"
)

// This is a *very* simple GitHub API client. We just need the list of repositories
// from the API and will handle the rest through ordinary Git commands.
const (
	repositoryAPI       = "https://api.github.com/users/lab5e/repos?per_page=100&ype=public"
	topicSampleName     = "sample"
	topicLanguagePrefix = "lang-"
)

type repo struct {
	Name          string   `json:"name"`
	URL           string   `json:"svn_url"`
	Description   string   `json:"description"`
	Topics        []string `json:"topics"`
	DefaultBranch string   `json:"default_branch"`
	Private       bool     `json:"private"`
	Archived      bool     `json:"archived"`
}

// IsSample returns true if the topic "sample" is set on the repository and
// it's not archived
func (r *repo) IsSample() bool {
	if r.Archived {
		return false
	}
	for _, t := range r.Topics {
		if t == topicSampleName {
			return true
		}
	}
	return false
}

// Language returns the flagged language for samples
func (r *repo) Language() string {
	for _, t := range r.Topics {
		if strings.HasPrefix(t, topicLanguagePrefix) {
			val := strings.Split(t, "-")
			if len(val) == 2 {
				return val[1]
			}
			return "*** unformatted ***"
		}
	}
	return "*** unknown ***"

}

// Keywords returns a concatenated list of topics/keywords for the repo
func (r *repo) Keywords() string {
	return strings.Join(r.Topics, ",")
}

// ListSamples lists all samples in the Lab5e GitHub repository. The
// samples must be tagged with the keywords "sample" to show up.
//
// Samples are named with <language>-<description> and each sample is in its own
// repository.
func ListSamples(format commonopt.ListFormat) error {
	fmt.Println("Reading samples...")
	//
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Get(repositoryAPI)
	if err != nil {
		return fmt.Errorf("error reading GitHub repositories: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 OK but got %d %s from GitHub API", resp.StatusCode, resp.Status)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var repoList []repo
	if err := json.Unmarshal(buf, &repoList); err != nil {
		return err
	}

	t := helpers.NewTableOutput(format)
	t.SetTitle("Available samples")
	t.AppendHeader(table.Row{
		"Language",
		"Name",
		"Description",
		"Keywords",
		"Default Branch",
	})
	for _, sample := range repoList {

		if sample.IsSample() {
			t.AppendRow(table.Row{
				sample.Language(),
				sample.Name,
				sample.Description,
				sample.Keywords(),
				sample.DefaultBranch,
			})
		}
	}
	helpers.RenderTable(t, format.Format)

	return nil
}

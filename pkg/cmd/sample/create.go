package sample

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type createSample struct {
	Name string `required:"yes" description:"Name of sample" long:"name"`
}

// Execute runs the version command
func (c *createSample) Execute([]string) error {
	fmt.Printf("Creating sample %s\n", c.Name)
	samples, err := sampleRepos()
	if err != nil {
		return err
	}
	for _, sample := range samples {
		if sample.Name == c.Name {
			return c.createSample(sample)
		}
	}
	return nil
}

func (c *createSample) createSample(sample repo) error {
	_, err := os.Stat(sample.Name)
	if err == nil {
		return fmt.Errorf("%s already exists", sample.Name)
	}

	// Clone the repo, remove .git directory and init new directory
	_, err = git.PlainClone(sample.Name, false, &git.CloneOptions{
		URL:           sample.URL,
		ReferenceName: plumbing.ReferenceName(sample.DefaultBranch),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return err
	}

	// Strip the .git and .github directories to make it a clean directory. Won't
	// check for errors
	os.RemoveAll(path.Join(sample.Name, ".git"))
	os.RemoveAll(path.Join(sample.Name, ".github"))

	fmt.Printf("Created the %s sample\n", sample.Name)
	return nil
}

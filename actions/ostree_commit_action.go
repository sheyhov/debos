package actions

import (
	"log"
	"os"
	"path"

	"github.com/go-debos/debos"
	"github.com/sjoerdsimons/ostree-go/pkg/otbuiltin"
)

type OstreeCommitAction struct {
	debos.BaseAction `yaml:",inline"`
	Repository string
	Branch     string
	Subject    string
	Command    string
}

func emptyDir(dir string) {
	d, _ := os.Open(dir)
	defer d.Close()
	files, _ := d.Readdirnames(-1)
	for _, f := range files {
		os.RemoveAll(f)
	}
}

func (ot *OstreeCommitAction) Run(context *debos.DebosContext) error {
	ot.LogStart()
	repoPath := path.Join(context.Artifactdir, ot.Repository)

	emptyDir(path.Join(context.Rootdir, "dev"))

	repo, err := otbuiltin.OpenRepo(repoPath)
	if err != nil {
		return err
	}

	_, err = repo.PrepareTransaction()
	if err != nil {
		return err
	}

	opts := otbuiltin.NewCommitOptions()
	opts.Subject = ot.Subject
	ret, err := repo.Commit(context.Rootdir, ot.Branch, opts)
	if err != nil {
		return err
	} else {
		log.Printf("Commit: %s\n", ret)
	}
	_, err = repo.CommitTransaction()
	if err != nil {
		return err
	}

	return nil
}

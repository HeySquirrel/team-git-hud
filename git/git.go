package git

import (
	"fmt"
	tlog "github.com/heysquirrel/tribe/log"
	"github.com/heysquirrel/tribe/shell"
	"strings"
	"time"
)

type Repo struct {
	shell  *shell.Shell
	logger *tlog.Log
	logs   Logs
}

type File struct {
	Name         string
	Contributors Contributors
	Related      []*RelatedFile
	WorkItems    []string
}

func (repo *Repo) git(args ...string) (string, error) {
	return repo.shell.Exec("git", args...)
}

func New(dir string, logger *tlog.Log) (*Repo, error) {
	temp := shell.New(dir, logger)
	out, err := temp.Exec("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}

	repo := new(Repo)
	repo.shell = shell.New(strings.TrimSpace(out), logger)
	repo.logger = logger

	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	repo.logs, err = repo.LogsAfter(sixMonthsAgo)
	if err != nil {
		return nil, err
	}

	logger.Add(fmt.Sprintf("Processed %d logs", len(repo.logs)))

	return repo, err
}

func (repo *Repo) Changes() []*File {
	var results = make([]*File, 0)

	cmdOut, err := repo.git("status", "--porcelain")
	if err != nil {
		repo.logger.Add(err.Error())
		return results
	}

	output := strings.Split(cmdOut, "\n")
	for _, change := range output {
		if len(change) > 0 {
			filename := change[3:len(change)]
			results = append(results, repo.GetFile(filename))
		}
	}

	return results
}

func (repo *Repo) GetFile(filename string) *File {
	logs := repo.logs.ContainsFile(filename)

	file := new(File)
	file.Name = filename
	file.Related = logs.relatedFiles(filename)
	file.Contributors = logs.relatedContributors()
	file.WorkItems = logs.relatedWorkItems()

	return file
}

func (repo *Repo) Related(filename string) ([]*RelatedFile, []string, []*Contributor) {
	relatedLogs := repo.logs.ContainsFile(filename)

	return relatedLogs.relatedFiles(filename), relatedLogs.relatedWorkItems(), relatedLogs.relatedContributors()
}

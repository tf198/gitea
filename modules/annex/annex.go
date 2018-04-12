package annex

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/setting"
)

func ContentLocation(key, repoDir string) (string, error) {

	var gitcmd *exec.Cmd
	gitcmd = exec.Command("git", "annex", "contentlocation", key)
	gitcmd.Dir = repoDir

	keyPath, err := gitcmd.Output()

	if err != nil {
		return "", err
	}

	return path.Join(repoDir, strings.Trim(string(keyPath), "\n")), nil
}

func GetAnnexHandler(ctx *context.Context) {

	if !setting.GitAnnex.Enabled {
		ctx.NotFound("Git Annex not enabled", nil)
	}

	filepath, err := ContentLocation(ctx.Params("key"), ctx.Repo.GitRepo.Path)

	if err != nil {
		ctx.NotFound("Invalid key", err)
		return
	}

	f, err := os.Open(filepath)

	if err != nil {
		ctx.NotFound("Not on this server", err)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		ctx.ServerError("Failed to get modification time", err)
		return
	}

	ctx.ServeContent(ctx.Params("key"), f, stat.ModTime())
}

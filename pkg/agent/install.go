package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xoctopus/x/misc/must"
)

func SkillInstallationPlan(root, gomodcache string, skills []SkillDef) (*Plan, error) {
	dir := filepath.Join(root, ".agents", "skills")

	plan := &Plan{
		SkillsDir:     dir,
		GitIgnorePath: filepath.Join(root, ".gitignore"),
		GitIgnores:    make([]string, 0, len(skills)),
		Skills:        make([]SkillInstall, 0, len(skills)),
	}

	for _, skill := range skills {
		s := skill
		if s.Replace != "" && !filepath.IsAbs(s.Replace) {
			s.Replace = filepath.Join(root, s.Replace)
		}

		install, err := s.Installation(dir, gomodcache)
		if err != nil {
			return nil, err
		}

		plan.GitIgnores = append(plan.GitIgnores, fmt.Sprintf("skills/%s", s.Name))
		plan.Skills = append(plan.Skills, *install)
	}

	return plan, nil
}

type Plan struct {
	SkillsDir     string
	GitIgnorePath string
	GitIgnores    []string
	Skills        []SkillInstall
}

func (p *Plan) Apply() error {
	must.NoErrorF(os.MkdirAll(p.SkillsDir, 0o755), "failed to create skill dir: %q", p.SkillsDir)

	for _, skill := range p.Skills {
		mustMakeSymlink(skill.Destination, skill.Source)
	}

	mustAppendSkillsGitIgnore(p.GitIgnorePath, p.GitIgnores)
	return nil
}

type Installer struct {
	Root       string
	GoModFile  string
	GoModCache string
}

func (i *Installer) init(_ context.Context) {
	if len(i.GoModCache) == 0 {
		i.GoModCache = GOMODCACHE
	}

	if i.Root == "" {
		i.Root = must.NoErrorV(os.Getwd())
	}

	for dir := i.Root; ; dir = filepath.Dir(dir) {
		mod := filepath.Join(dir, "go.mod")

		_, valid := fileCheck(mod)
		if valid {
			i.Root = dir
			i.GoModFile = mod
			break
		}

		if parent := filepath.Dir(dir); parent == dir {
			break
		}
	}

	_, validRoot := dirCheck(i.Root)
	must.BeTrueF(validRoot, "invalid project root: %q", i.Root)
	_, validModFile := fileCheck(i.GoModFile)
	must.BeTrueF(validModFile, "invalid project mod file: %q", i.GoModFile)
}

func (i *Installer) plan(_ context.Context) (*Plan, error) {
	skills, err := ExtractSkillRefsFromGoMod(i.GoModFile)
	if err != nil {
		return nil, err
	}

	return SkillInstallationPlan(i.Root, i.GoModCache, skills)
}

func (i *Installer) Install(ctx context.Context) error {
	i.init(ctx)

	plan, err := i.plan(ctx)
	if err != nil {
		return err
	}
	return plan.Apply()
}

type SkillInstaller struct{}

func (s *SkillInstaller) Run(ctx context.Context) error {
	return (&Installer{}).Install(ctx)
}

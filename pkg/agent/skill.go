package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xoctopus/x/misc/must"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

var GOMODCACHE string

func init() {
	GOMODCACHE = strings.TrimSpace(os.Getenv("GOMODCACHE"))
	if len(GOMODCACHE) == 0 {
		out, _ := exec.Command("go", "env", "GOMODCACHE").Output()
		GOMODCACHE = strings.TrimSpace(string(out))
	}
	must.BeTrueF(len(GOMODCACHE) > 0, "failed to solve GOMODCACHE, got empty")
}

func ExtractSkillRefsFromGoMod(path string) ([]SkillDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	f, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	skills := make([]SkillDef, 0)
	for _, r := range f.Require {
		if !r.Indirect {
			for _, c := range r.Syntax.Comments.Before {
				if len(c.Token) > 0 {
					token := strings.TrimSpace(strings.TrimPrefix(c.Token, "//"))
					if name, found := strings.CutPrefix(token, "+skill:"); found {
						skills = append(skills, SkillDef{
							Name:    name,
							Path:    r.Mod.Path,
							Version: r.Mod.Version,
						})
					}
				}
			}
		}
	}

	replaces := make(map[string][]*modfile.Replace)
	for _, r := range f.Replace {
		replaces[r.Old.Path] = append(replaces[r.Old.Path], r)
	}

	for i := range skills {
		entries, ok := replaces[skills[i].Path]
		if !ok {
			continue
		}

		for _, entry := range entries {
			if entry.Old.Version != "" && entry.Old.Version == skills[i].Version ||
				entry.Old.Version == "" {
				if filepath.IsAbs(entry.New.Path) || strings.HasPrefix(entry.New.Path, ".") {
					skills[i].Replace = entry.New.Path
					break
				}
				skills[i].Path = entry.New.Path
				skills[i].Version = entry.New.Version
				break
			}
		}
	}

	return skills, nil
}

type SkillInstall struct {
	Ref         SkillDef
	ModuleRoot  string
	Source      string
	Destination string
}

type SkillDef struct {
	Name    string
	Path    string
	Version string
	Replace string
}

func (s *SkillDef) Installation(dst, gomodcache string) (*SkillInstall, error) {
	var (
		root string
		src  string
	)
	if len(s.Replace) > 0 {
		root = s.Replace
		src = filepath.Join(s.Replace, ".agents", "skills", s.Name)
	} else {
		mod := must.NoErrorV(module.EscapePath(s.Path))
		version := must.NoErrorV(module.EscapeVersion(s.Version))
		root = filepath.Join(gomodcache, mod+"@"+version)
		src = filepath.Join(root, ".agents", "skills", s.Name)
	}

	for _, dir := range []string{root, src} {
		info, err := os.Stat(dir)
		if err == nil {
			if !info.IsDir() {
				err = fmt.Errorf("invalid skill source %s", dir)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to stat skill dir %s: %w", dir, err)
		}
	}
	return &SkillInstall{
		Ref:         *s,
		ModuleRoot:  root,
		Source:      src,
		Destination: filepath.Join(dst, s.Name),
	}, nil
}

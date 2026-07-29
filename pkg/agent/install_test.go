package agent_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/bdd"

	"github.com/xoctopus/genx/pkg/agent"
)

func TestSkillInstallationPlan(t *testing.T) {
	bdd.From(t).Given("skill definitions with local and module sources", func(t bdd.T) {
		root := t.TempDir()
		gomodcache := filepath.Join(root, "modcache")

		localBar := filepath.Join(root, "local-bar")
		_ = os.MkdirAll(filepath.Join(localBar, ".agents", "skills", "bar"), 0o755)

		fooModRoot := filepath.Join(gomodcache, "github.com/a/foo@v1.0.0")
		_ = os.MkdirAll(filepath.Join(fooModRoot, ".agents", "skills", "foo"), 0o755)

		skills := []agent.SkillDef{
			{
				Name:    "foo",
				Path:    "github.com/a/foo",
				Version: "v1.0.0",
			},
			{
				Name:    "bar",
				Path:    "github.com/b/bar",
				Version: "v1.2.0",
				Replace: "./local-bar",
			},
		}

		t.When("building the installation plan", func(t bdd.T) {
			plan, err := agent.SkillInstallationPlan(root, gomodcache, skills)

			skillsDir := filepath.Join(root, ".agents", "skills")
			t.Then(
				"it should succeed with resolved paths and gitignore entries",
				bdd.Succeed(err),
				bdd.NotBeNil[*agent.Plan](plan),
				bdd.Equal(skillsDir, plan.SkillsDir),
				bdd.Equal(filepath.Join(root, ".gitignore"), plan.GitIgnorePath),
				bdd.HaveLen(plan.GitIgnores, 2),
				bdd.Equal("skills/foo", plan.GitIgnores[0]),
				bdd.Equal("skills/bar", plan.GitIgnores[1]),
				bdd.HaveLen(plan.Skills, 2),

				bdd.Equal(fooModRoot, plan.Skills[0].ModuleRoot),
				bdd.Equal(filepath.Join(fooModRoot, ".agents", "skills", "foo"), plan.Skills[0].Source),
				bdd.Equal(filepath.Join(skillsDir, "foo"), plan.Skills[0].Destination),

				bdd.Equal(localBar, plan.Skills[1].ModuleRoot),
				bdd.Equal(filepath.Join(localBar, ".agents", "skills", "bar"), plan.Skills[1].Source),
				bdd.Equal(filepath.Join(skillsDir, "bar"), plan.Skills[1].Destination),
				bdd.Equal(localBar, plan.Skills[1].Ref.Replace),
			)
		})
	})

	bdd.From(t).Given("a skill with absolute replace path", func(t bdd.T) {
		root := t.TempDir()
		gomodcache := filepath.Join(root, "modcache")
		absReplace := filepath.Join(root, "abs-skill")
		_ = os.MkdirAll(filepath.Join(absReplace, ".agents", "skills", "abs"), 0o755)

		skills := []agent.SkillDef{{
			Name:    "abs",
			Path:    "github.com/c/abs",
			Version: "v0.1.0",
			Replace: absReplace,
		}}

		t.When("building the installation plan", func(t bdd.T) {
			plan, err := agent.SkillInstallationPlan(root, gomodcache, skills)
			t.Then(
				"it should keep the absolute replace path",
				bdd.Succeed(err),
				bdd.Equal(absReplace, plan.Skills[0].ModuleRoot),
				bdd.Equal(absReplace, plan.Skills[0].Ref.Replace),
				bdd.Equal("skills/abs", plan.GitIgnores[0]),
			)
		})
	})

	bdd.From(t).Given("a skill whose source does not exist", func(t bdd.T) {
		root := t.TempDir()
		skills := []agent.SkillDef{{
			Name:    "missing",
			Path:    "github.com/x/missing",
			Version: "v1.0.0",
		}}

		t.When("building the installation plan", func(t bdd.T) {
			_, err := agent.SkillInstallationPlan(root, filepath.Join(root, "modcache"), skills)
			t.Then(
				"it should fail",
				bdd.Failed(err),
				bdd.ErrorContains(err, "failed to stat skill dir"),
			)
		})
	})
}

func TestPlanApply(t *testing.T) {
	bdd.From(t).Given("an installation plan with skill sources", func(t bdd.T) {
		root := t.TempDir()
		skillsDir := filepath.Join(root, ".agents", "skills")
		gitignore := filepath.Join(root, ".gitignore")

		srcFoo := filepath.Join(root, "src", "foo")
		srcBar := filepath.Join(root, "src", "bar")
		_ = os.MkdirAll(srcFoo, 0o755)
		_ = os.MkdirAll(srcBar, 0o755)
		_ = os.WriteFile(filepath.Join(srcFoo, "SKILL.md"), []byte("foo"), 0o644)
		_ = os.WriteFile(gitignore, []byte("node_modules/\n"), 0o644)

		plan := &agent.Plan{
			SkillsDir:     skillsDir,
			GitIgnorePath: gitignore,
			GitIgnores:    []string{"skills/foo", "skills/bar"},
			Skills: []agent.SkillInstall{
				{
					Source:      srcFoo,
					Destination: filepath.Join(skillsDir, "foo"),
				},
				{
					Source:      srcBar,
					Destination: filepath.Join(skillsDir, "bar"),
				},
			},
		}

		t.When("applying the plan", func(t bdd.T) {
			err := plan.Apply()

			fooInfo, fooErr := os.Lstat(filepath.Join(skillsDir, "foo"))
			barInfo, barErr := os.Lstat(filepath.Join(skillsDir, "bar"))
			fooTarget, fooReadErr := os.Readlink(filepath.Join(skillsDir, "foo"))
			barTarget, barReadErr := os.Readlink(filepath.Join(skillsDir, "bar"))
			ignoreData, ignoreErr := os.ReadFile(gitignore)

			t.Then(
				"it should create skill symlinks and append gitignore entries",
				bdd.Succeed(err),
				bdd.Succeed(fooErr),
				bdd.Succeed(barErr),
				bdd.BeTrue(fooInfo.Mode()&os.ModeSymlink != 0),
				bdd.BeTrue(barInfo.Mode()&os.ModeSymlink != 0),
				bdd.Succeed(fooReadErr),
				bdd.Succeed(barReadErr),
				bdd.Equal(srcFoo, fooTarget),
				bdd.Equal(srcBar, barTarget),
				bdd.Succeed(ignoreErr),
				bdd.ContainsSubString(string(ignoreData), "node_modules/"),
				bdd.ContainsSubString(string(ignoreData), "skills/foo"),
				bdd.ContainsSubString(string(ignoreData), "skills/bar"),
			)
		})

		t.When("applying the plan again with existing symlinks", func(t bdd.T) {
			err := plan.Apply()
			ignoreData, _ := os.ReadFile(gitignore)
			countFoo := strings.Count(string(ignoreData), "skills/foo")
			countBar := strings.Count(string(ignoreData), "skills/bar")

			fooTarget, _ := os.Readlink(filepath.Join(skillsDir, "foo"))
			t.Then(
				"it should recreate symlinks without duplicating gitignore entries",
				bdd.Succeed(err),
				bdd.Equal(srcFoo, fooTarget),
				bdd.Equal(1, countFoo),
				bdd.Equal(1, countBar),
			)
		})
	})
}

func TestInstallerInstall(t *testing.T) {
	bdd.From(t).Given("a project with go.mod skill refs and local replace", func(t bdd.T) {
		root := t.TempDir()
		gomodcache := filepath.Join(root, "modcache")
		localBar := filepath.Join(root, "local-bar")

		_ = os.MkdirAll(filepath.Join(localBar, ".agents", "skills", "bar"), 0o755)
		_ = os.WriteFile(filepath.Join(localBar, ".agents", "skills", "bar", "SKILL.md"), []byte("bar"), 0o644)

		fooModRoot := filepath.Join(gomodcache, "github.com/a/foo@v1.0.0")
		_ = os.MkdirAll(filepath.Join(fooModRoot, ".agents", "skills", "foo"), 0o755)
		_ = os.WriteFile(filepath.Join(fooModRoot, ".agents", "skills", "foo", "SKILL.md"), []byte("foo"), 0o644)

		modContent := `module github.com/test/mod

go 1.20

require (
	// +skill:foo
	github.com/a/foo v1.0.0
	// +skill:bar
	github.com/b/bar v1.2.0
)

replace github.com/b/bar => ./local-bar
`
		_ = os.WriteFile(filepath.Join(root, "go.mod"), []byte(modContent), 0o644)

		subdir := filepath.Join(root, "pkg", "sub")
		_ = os.MkdirAll(subdir, 0o755)

		installer := &agent.Installer{
			Root:       subdir,
			GoModCache: gomodcache,
		}

		t.When("installing skills", func(t bdd.T) {
			err := installer.Install(context.Background())

			skillsDir := filepath.Join(root, ".agents", "skills")
			fooLink := filepath.Join(skillsDir, "foo")
			barLink := filepath.Join(skillsDir, "bar")
			fooTarget, fooErr := os.Readlink(fooLink)
			barTarget, barErr := os.Readlink(barLink)
			ignoreData, ignoreErr := os.ReadFile(filepath.Join(root, ".gitignore"))

			t.Then(
				"it should walk up to go.mod root and install skills",
				bdd.Succeed(err),
				bdd.Equal(root, installer.Root),
				bdd.Equal(filepath.Join(root, "go.mod"), installer.GoModFile),
				bdd.Succeed(fooErr),
				bdd.Succeed(barErr),
				bdd.Equal(filepath.Join(fooModRoot, ".agents", "skills", "foo"), fooTarget),
				bdd.Equal(filepath.Join(localBar, ".agents", "skills", "bar"), barTarget),
				bdd.Succeed(ignoreErr),
				bdd.ContainsSubString(string(ignoreData), "skills/foo"),
				bdd.ContainsSubString(string(ignoreData), "skills/bar"),
			)
		})
	})

	bdd.From(t).Given("a root without go.mod", func(t bdd.T) {
		root := t.TempDir()
		installer := &agent.Installer{
			Root:       root,
			GoModCache: filepath.Join(root, "modcache"),
		}

		t.When("installing skills", func(t bdd.T) {
			panicked := false
			func() {
				defer func() {
					err := recover()
					panicked = err != nil
					testx.Expect(t, err.(error), testx.ErrorContains("invalid project mod file"))
				}()
				_ = installer.Install(context.Background())
			}()
			t.Then(
				"it should panic on invalid project root",
				bdd.BeTrue(panicked),
			)
		})
	})
}

func TestSkillInstallerRun(t *testing.T) {
	bdd.From(t).Given("cwd is a project with go.mod skill refs and local replace", func(t bdd.T) {
		root := t.TempDir()
		localBar := filepath.Join(root, "local-bar")
		_ = os.MkdirAll(filepath.Join(localBar, ".agents", "skills", "bar"), 0o755)
		_ = os.WriteFile(filepath.Join(localBar, ".agents", "skills", "bar", "SKILL.md"), []byte("bar"), 0o644)

		modContent := `module github.com/test/mod

go 1.20

require (
	// +skill:bar
	github.com/b/bar v1.2.0
)

replace github.com/b/bar => ./local-bar
`
		_ = os.WriteFile(filepath.Join(root, "go.mod"), []byte(modContent), 0o644)
		t.Chdir(root)

		t.When("running SkillInstaller", func(t bdd.T) {
			err := (&agent.SkillInstaller{}).Run(context.Background())

			barLink := filepath.Join(root, ".agents", "skills", "bar")
			barTarget, barErr := os.Readlink(barLink)
			ignoreData, ignoreErr := os.ReadFile(filepath.Join(root, ".gitignore"))

			t.Then(
				"it should install skills from the working directory project",
				bdd.Succeed(err),
				bdd.Succeed(barErr),
				bdd.Equal(filepath.Join(localBar, ".agents", "skills", "bar"), barTarget),
				bdd.Succeed(ignoreErr),
				bdd.ContainsSubString(string(ignoreData), "skills/bar"),
			)
		})
	})
}

package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xoctopus/genx/pkg/agent"
	"github.com/xoctopus/x/testx/bdd"
)

func TestExtractSkillRefsFromGoMod(t *testing.T) {
	bdd.From(t).Given("a temporary go.mod file", func(t bdd.T) {
		dir := t.TempDir()
		modPath := filepath.Join(dir, "go.mod")
		modContent := `module github.com/test/mod

go 1.20

require (
	// +skill:foo
	github.com/a/foo v1.0.0
	// +skill:bar
	github.com/b/bar v1.2.0
	// +skill:baz-1
	// +skill:baz-2
	github.com/c/baz v1.5.0
	github.com/no-skill v1.1.1
)

replace (
	github.com/b/bar => ./local-bar
	github.com/a/foo => github.com/b/foo v1.0.1
)
`
		_ = os.WriteFile(modPath, []byte(modContent), 0644)

		t.When("extracting skill refs", func(t bdd.T) {
			skills, err := agent.ExtractSkillRefsFromGoMod(modPath)

			t.Then(
				"it should succeed",
				bdd.Succeed(err),
				bdd.HaveLen(skills, 4),

				bdd.Equal("foo", skills[0].Name),
				bdd.Equal("github.com/b/foo", skills[0].Path),
				bdd.Equal("v1.0.1", skills[0].Version),
				bdd.Equal("", skills[0].Replace),

				bdd.Equal("bar", skills[1].Name),
				bdd.Equal("github.com/b/bar", skills[1].Path),
				bdd.Equal("v1.2.0", skills[1].Version),
				bdd.Equal("./local-bar", skills[1].Replace),

				bdd.Equal("baz-1", skills[2].Name),
				bdd.Equal("github.com/c/baz", skills[2].Path),
				bdd.Equal("v1.5.0", skills[2].Version),
				bdd.Equal("baz-2", skills[3].Name),
				bdd.Equal("github.com/c/baz", skills[3].Path),
				bdd.Equal("v1.5.0", skills[3].Version),
			)
		})

		t.When("extracting from non-existent go.mod", func(t bdd.T) {
			_, err := agent.ExtractSkillRefsFromGoMod(filepath.Join(dir, "missing.mod"))
			t.Then(
				"it should fail",
				bdd.Failed(err),
				bdd.ErrorContains(err, "failed to read"),
			)
		})
	})

	bdd.From(t).Given("a invalid go mod file", func(t bdd.T) {
		dir := t.TempDir()
		modPath := filepath.Join(dir, "go.mod")
		modContent := `module github.com/test/mod

go 1.20

require (
	// +skill:foo
	github.com/a/foo v1.0.0
`
		_ = os.WriteFile(modPath, []byte(modContent), 0644)

		t.When("extracting skill refs", func(t bdd.T) {
			skills, err := agent.ExtractSkillRefsFromGoMod(modPath)
			t.Then(
				"it should fail",
				bdd.Failed(err),
				bdd.HaveLen(skills, 0),
			)
		})
	})

}

func TestSkillDefInstallation(t *testing.T) {
	bdd.From(t).Given("a skill definition", func(t bdd.T) {
		t.When("using a local replace", func(t bdd.T) {
			skill := agent.SkillDef{
				Name:    "bar",
				Path:    "github.com/b/bar",
				Version: "v1.2.0",
				Replace: "./local-bar",
			}

			dir := t.TempDir()
			dstDir := filepath.Join(dir, "dst")
			modCacheDir := filepath.Join(dir, "modCache")
			localBarDir := filepath.Join(dir, "local-bar")

			skill.Replace = localBarDir

			_ = os.MkdirAll(filepath.Join(localBarDir, ".agents/skills/bar"), 0755)

			install, err := skill.Installation(dstDir, modCacheDir)
			t.Then(
				"it should succeed",
				bdd.Succeed(err),
				bdd.NotBeNil[*agent.SkillInstall](install),
				bdd.Equal(localBarDir, install.ModuleRoot),
				bdd.Equal(filepath.Join(localBarDir, ".agents/skills/bar"), install.Source),
				bdd.Equal(filepath.Join(dstDir, "bar"), install.Destination),
			)
		})

		t.When("using standard module path", func(t bdd.T) {
			skill := agent.SkillDef{
				Name:    "foo",
				Path:    "github.com/a/Foo",
				Version: "v1.0.0",
			}

			dir := t.TempDir()
			dstDir := filepath.Join(dir, "dst")
			modCacheDir := filepath.Join(dir, "modCache")

			expectedRoot := filepath.Join(modCacheDir, "github.com/a/!foo@v1.0.0")
			expectedSrc := filepath.Join(expectedRoot, ".agents/skills/foo")

			_ = os.MkdirAll(expectedSrc, 0755)

			install, err := skill.Installation(dstDir, modCacheDir)
			t.Then(
				"it should succeed",
				bdd.Succeed(err),
				bdd.Equal(expectedRoot, install.ModuleRoot),
				bdd.Equal(expectedSrc, install.Source),
				bdd.Equal(filepath.Join(dstDir, "foo"), install.Destination),
			)
		})

		t.When("source directory is not a directory", func(t bdd.T) {
			skill := agent.SkillDef{
				Name:    "bad",
				Path:    "github.com/b/bad",
				Version: "v1.0.0",
				Replace: "./bad-replace",
			}

			_ = os.MkdirAll("./bad-replace/.agents/skills", 0755)
			_ = os.WriteFile("./bad-replace/.agents/skills/bad", []byte("not a dir"), 0644)
			defer func() { _ = os.RemoveAll("./bad-replace") }()

			_, err := skill.Installation("/dst", "/modcache")
			t.Then(
				"it should fail",
				bdd.Failed(err),
				bdd.ErrorContains(err, "invalid skill source"),
			)
		})
	})
}

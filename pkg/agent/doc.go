// Package agent installs Agent Skills declared in a project's go.mod into the
// local `.agents/skills` directory.
//
// Skills are versioned with Go modules. A consuming project marks required
// skills on require statements; this package resolves those modules (including
// replace directives), creates symlinks under the project root, and appends
// matching entries to `.gitignore`.
//
// # How to Integrate
//
// Annotate each skill-providing require dependency in go.mod with
// `// +skill:<name>` on the lines above it. One require may declare multiple
// skills.
//
//	module example.com/app
//
//	go 1.22
//
//	require (
//		// +skill:genx
//		github.com/xoctopus/genx v0.3.1
//	)
//
//	replace github.com/xoctopus/genx => ../genx
//
// Then run [Installer.Install] from the project (or any subdirectory). The
// zero-value installer walks up to the nearest go.mod and installs all
// declared skills:
//
//	import (
//		"context"
//		"fmt"
//		"os"
//
//		"github.com/xoctopus/genx/pkg/agent"
//	)
//
//	func main() {
//		if err := (&agent.Installer{}).Install(context.Background()); err != nil {
//			_, _ = fmt.Fprintln(os.Stderr, err)
//			os.Exit(1)
//			return
//		}
//	}
//
// Replace directives are honored: a local or absolute replace path is used as
// the module root; a module replace updates the resolved path and version.
//
// [SkillInstaller] is a thin wrapper that calls [Installer.Install] with a
// zero-value installer. For more control, set [Installer.Root] or
// [Installer.GoModCache] before calling Install. Lower-level helpers
// [ExtractSkillRefsFromGoMod], [SkillDef.Installation], and
// [SkillInstallationPlan] are available when building a custom plan.
//
// # Skill Layout
//
// Each skill module must ship its content at:
//
//	.agents/skills/<name>/
//
// For example, genx provides `.agents/skills/genx/SKILL.md` with supporting
// references. After installation, the consuming project exposes the same path
// via a symlink:
//
//	<project>/.agents/skills/<name> -> <module-root>/.agents/skills/<name>
//
// # Installation Plan
//
// [SkillInstallationPlan] produces a [Plan] that:
//
//  1. Resolves each skill source under GOMODCACHE (`<path>@<version>`) or a
//     replace path.
//  2. Maps destinations to `<root>/.agents/skills/<name>`.
//  3. Records gitignore lines as `skills/<name>`.
//
// [Plan.Apply] creates the skills directory, replaces any existing destination
// with a symlink to the source, and appends missing gitignore entries.
package agent

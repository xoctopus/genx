package agent

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/xoctopus/x/misc/must"
)

func mustMakeSymlink(dst string, src string) {
	exists, _ := symlinkCheck(dst)

	if exists {
		must.NoErrorF(os.RemoveAll(dst), "failed to remove existing path: %q", dst)
	}

	must.NoErrorF(os.Symlink(src, dst), "failed to make symlink %q->%q", dst, src)
}

func mustAppendSkillsGitIgnore(path string, names []string) {
	must.NoErrorF(os.MkdirAll(filepath.Dir(path), 0o755), "failed to create ignore dir: %q", path)

	skills := make(map[string]struct{})
	for _, name := range names {
		skills[name] = struct{}{}
	}

	data, err := os.ReadFile(path)
	must.BeTrueF(err == nil || os.IsNotExist(err), "failed to read ignore file: %q", path)

	for line := range strings.SplitSeq(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		delete(skills, trimmed)
	}

	if len(skills) == 0 {
		return
	}

	names = slices.Collect(maps.Keys(skills))
	sort.Strings(names)

	for _, name := range names {
		if len(data) > 0 && data[len(data)-1] != '\n' {
			data = append(data, '\n')
		}
		data = append(data, []byte(name)...)
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	must.NoErrorF(os.WriteFile(path, data, 0o644), "failed to write ignore file: %q", path)
}

func fileCheck(path string) (exist, valid bool) {
	info, err := os.Stat(path)
	must.BeTrueF(err == nil || os.IsNotExist(err), "failed to stat: %q", path)

	exist = err == nil
	if exist {
		valid = !info.IsDir()
	}
	return
}

func dirCheck(path string) (exists, valid bool) {
	info, err := os.Stat(path)
	must.BeTrueF(err == nil || os.IsNotExist(err), "failed to stat: %q", path)

	exists = err == nil
	if exists {
		valid = info.IsDir()
	}
	return
}

func symlinkCheck(path string) (exist, valid bool) {
	info, err := os.Lstat(path)
	must.BeTrueF(err == nil || os.IsNotExist(err), "failed to stat: %q", path)

	exist = err == nil
	if exist {
		valid = info.Mode()&os.ModeSymlink != 0
	}
	return
}

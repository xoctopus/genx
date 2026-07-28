package helper

import (
	"runtime/debug"
)

func VersionFor(target string) string {
	if i, ok := debug.ReadBuildInfo(); ok {
		if target == "" {
			return i.Main.Version
		}
		for _, m := range i.Deps {
			if m.Path == target {
				return m.Version
			}
		}
	}
	return "devel"
}

package helper

import (
	"cmp"
	"runtime/debug"
)

func VersionFor(target string) string {
	v := ""
	if i, ok := debug.ReadBuildInfo(); ok {
		if target == "" {
			v = i.Main.Version
		} else {
			for _, m := range i.Deps {
				if m.Path == target {
					v = m.Version
					break
				}
			}
		}
	}
	return cmp.Or(v, "devel")
}

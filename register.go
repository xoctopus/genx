package genx

import (
	"sort"

	"github.com/xoctopus/x/mapx"
)

var generators = mapx.NewXmap[string, Generator]()

func RegisterGenerator(g Generator) {
	generators.LoadOrStore(g.Identifier(), g)
}

func GetGenerators(identifiers ...string) (gs []Generator) {
	defer func() {
		sort.Slice(gs, func(i, j int) bool {
			return gs[i].Identifier() < gs[j].Identifier()
		})
	}()

	if len(identifiers) == 0 {
		generators.Range(func(_ string, g Generator) bool {
			gs = append(gs, g)
			return true
		})
		return gs
	}

	ids := map[string]bool{}
	for _, id := range identifiers {
		ids[id] = false
	}
	generators.Range(func(id string, g Generator) bool {
		if scanned, ok := ids[id]; ok && !scanned {
			gs = append(gs, g)
			ids[id] = true
		}
		return true
	})
	return gs
}

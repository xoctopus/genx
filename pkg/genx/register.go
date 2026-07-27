package genx

import (
	"sort"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"
)

var gGenerators = syncx.NewXmap[string, Generator]()

// Register registers a new Generator.
// It panics if a generator with the same identifier has already been registered.
func Register(g Generator) {
	_, loaded := gGenerators.LoadOrStore(g.Identifier(), g)
	must.BeTrueF(!loaded, "generator '%s' has been registered", g.Identifier())
}

// Get retrieves registered generators based on the provided identifiers.
// If no identifiers are provided, it returns all registered generators.
func Get(identifiers ...string) (gs []Generator) {
	defer func() {
		sort.Slice(gs, func(i, j int) bool {
			return gs[i].Identifier() < gs[j].Identifier()
		})
	}()

	if len(identifiers) == 0 {
		gGenerators.Range(func(_ string, g Generator) bool {
			gs = append(gs, g)
			return true
		})
		return gs
	}

	ids := map[string]bool{}
	for _, id := range identifiers {
		ids[id] = false
	}
	gGenerators.Range(func(id string, g Generator) bool {
		if scanned, ok := ids[id]; ok && !scanned {
			gs = append(gs, g)
			ids[id] = true
		}
		return true
	})
	return gs
}

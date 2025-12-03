package dumper_test

import (
	"context"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/genx/internal/dumper"
)

func TestImport(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		i := dumper.NewImport(context.Background(), "github.com/xoctopus/genx/internal/dumper_test")
		Expect(t, i.Path(), Equal("github.com/xoctopus/genx/internal/dumper_test"))
		Expect(t, i.Name(), Equal("dumper_test"))
		Expect(t, i.Alias(), Equal("dumper_test"))
		Expect(t, i.Code(), Equal(`"github.com/xoctopus/genx/internal/dumper_test"`))
		i.MakeAlias()
		Expect(t, i.Code(), Equal(`"github.com/xoctopus/genx/internal/dumper_test"`))
		i.MakeAlias()
		Expect(t, i.Code(), Equal(`internal_dumper_test "github.com/xoctopus/genx/internal/dumper_test"`))
	})

	t.Run("MakeAlias", func(t *testing.T) {
		x := dumper.NewImport(context.Background(), "github.com/xoctopus/genx/testdata/bricks/x/p")

		x.MakeAlias()
		Expect(t, x.Alias(), Equal("x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("testdata_bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("genx_testdata_bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("xoctopus_genx_testdata_bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("com_xoctopus_genx_testdata_bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("github_com_xoctopus_genx_testdata_bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("github_com_xoctopus_genx_testdata_bricks_x_p"))
		x.MakeAlias()
		Expect(t, x.Alias(), Equal("github_com_xoctopus_genx_testdata_bricks_x_p"))
	})

	t.Run("HitCache", func(t *testing.T) {
		dumper.NewImport(context.Background(), "github.com/xoctopus/genx/internal/dumper_test")
	})
}

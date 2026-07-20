@def contextx.From
@def contextx.With
@def contextx.Must
@def contextx.Carry
@def T
--Contextx
type tCtx#T# struct{}

var (
	#T#From =  #contextx.From#[tCtx#T#, #T#]
	With#T# =  #contextx.With#[tCtx#T#, #T#]
	Must#T# =  #contextx.Must#[tCtx#T#, #T#]
	Carry#T# = #contextx.Carry#[tCtx#T#, #T#]
)

@def contextx.From
@def contextx.With
@def contextx.Must
@def contextx.Carry
@def contextx.Carrier
@def context.Context
@def T
@def TParams
@def TNames
--ContextxGeneric
type tCtx#T#[#TParams#] struct{}

// #T#From try extract #T# from context
func #T#From[#TParams#](ctx #context.Context#) (#T#[#TNames#], bool) {
	return #contextx.From#[tCtx#T#[#TNames#], #T#[#TNames#]](ctx)
}

// Must#T# asserts #T# from context
func Must#T#[#TParams#](ctx #context.Context#) #T#[#TNames#] {
	return #contextx.Must#[tCtx#T#[#TNames#], #T#[#TNames#]](ctx)
}

// With#T# inject #T# to context
func With#T#[#TParams#](ctx #context.Context#, v #T#[#TNames#]) #context.Context# {
	return #contextx.With#[tCtx#T#[#TNames#], #T#[#TNames#]](ctx, v)
}

// Carry#T# returns #T# context carrier
func Carry#T#[#TParams#](v #T#[#TNames#]) #contextx.Carrier# {
	return #contextx.Carry#[tCtx#T#[#TNames#], #T#[#TNames#]](v)
}

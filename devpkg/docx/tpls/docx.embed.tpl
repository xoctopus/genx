@def DocOf
@def Ref
@def Prefix
--Anonymous
if doc, ok := #DocOf#(#Ref#.#Field#, #Prefix#, names...); ok {
	return doc, true
}

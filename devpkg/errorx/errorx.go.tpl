
@def CodeType
@def MessagesVar
@def MessageKeyValues
-- Message
var #MessagesVar# = map[#CodeType#]string{
    #MessageKeyValues#
}

@def CodeType
@def MessagesVar
@def github.com/pkg/errors.WithStack
-- NewError
func New#CodeType#Error(code #CodeType#) error {
	return #github.com/pkg/errors.WithStack#(&#CodeType#Error{
		code: code,
		msg:  #MessagesVar#[code],
	})
}


@def CodeType
@def MessagesVar
@def fmt.Sprintf
@def github.com/pkg/errors.WithStack
-- NewErrorf
func New#CodeType#Errorf(code #CodeType#, format string, args ...any) error {
	return #github.com/pkg/errors.WithStack#(&#CodeType#Error{
		code: code,
		msg:  #fmt.Sprintf#(#MessagesVar#[code]+" "+format, args...),
	})
}

@def CodeType
@def MessagesVar
@def fmt.Sprintf
@def github.com/pkg/errors.WithStack
-- NewErrorWrap
func New#CodeType#ErrorWrap(code #CodeType#, cause error) error {
	if cause == nil {
		return nil
	}
	return #github.com/pkg/errors.WithStack#(&#CodeType#Error{
		code: code,
		msg:  #fmt.Sprintf#(#MessagesVar#[code]+" [cause: %+v]", cause),
	})
}


@def CodeType
@def MessagesVar
@def fmt.Sprintf
@def github.com/pkg/errors.WithStack
-- NewErrorWrapf
func New#CodeType#ErrorWrapf(code #CodeType#, cause error, format string, args ...any) error {
	if cause == nil {
		return nil
	}
	return #github.com/pkg/errors.WithStack#(&#CodeType#Error{
		code: code,
		msg: #fmt.Sprintf#(
			#MessagesVar#[code]+" [cause: %+v] "+format,
			append([]any{cause}, args...)...,
		),
	})
}

@def CodeType
-- ErrorDefine
type #CodeType#Error struct {
	code #CodeType#
	msg  string
}

@def CodeType
-- CodeType_Error
func (e *#CodeType#Error) Error() string {
	return e.msg
}

@def CodeType
-- CodeType_Code
func (e *#CodeType#Error) Code() #CodeType# {
	return e.code
}

@def CodeType
@def github.com/pkg/errors.As
-- CodeType_Is
func (e *#CodeType#Error) Is(err error) bool {
	var target *#CodeType#Error
	return #github.com/pkg/errors.As#(err, &target) && target.code == e.code
}

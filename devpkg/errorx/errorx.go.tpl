
@def CodeType
@def CodeMessageCases
@def UnknownCodeFormat
-- CodeMessage
func (e #CodeType#) Message() string {
	switch e {
	default:
		return fmt.Sprintf(#UnknownCodeFormat#, e)
	#CodeMessageCases#
	}
}

@def CodeType
-- NewError
func New#CodeType#Error(code #CodeType#) error {
	return &#CodeType#Error{
		code: code,
	}
}


@def CodeType
-- NewErrorf
func New#CodeType#Errorf(code #CodeType#, msg string, args ...any) error {
	return &#CodeType#Error{
		code: code,
		msg:  msg,
		args: args,
	}
}

@def CodeType
-- NewErrorWrap
func New#CodeType#ErrorWrap(code #CodeType#, cause error) error {
	if cause == nil {
		return nil
	}
	return &#CodeType#Error{
		code:  code,
		args:  []any{cause},
		cause: cause,
	}
}


@def CodeType
-- NewErrorWrapf
func New#CodeType#ErrorWrapf(code #CodeType#, cause error, msg string, args ...any) error {
	if cause == nil {
		return nil
	}
	return &#CodeType#Error{
		code:  code,
		msg:   msg,
		args:  append(args, cause),
		cause: cause,
	}
}

@def CodeType
-- ErrorDefine
type #CodeType#Error struct {
	code  #CodeType#
	msg   string
	args  []any
	cause error
}

@def CodeType
@def fmt.Sprintf
-- CodeType_Error
func (e *#CodeType#Error) Error() string {
	msg := e.code.Message()
	if len(e.msg) > 0 {
		msg += ". " + e.msg
	}
	if e.cause != nil {
		msg += ". [cause: %+v]"
	}
	return #fmt.Sprintf#(msg, e.args...)
}

@def CodeType
-- CodeType_Code
func (e *#CodeType#Error) Code() #CodeType# {
	return e.code
}

@def CodeType
@def errors.As
-- CodeType_Is
func (e *#CodeType#Error) Is(err error) bool {
	var target *#CodeType#Error
	return #errors.As#(err, &target) && target.code == e.code
}

@def CodeType
-- CodeType_Unwrap
func (e *#CodeType#Error) Unwrap() error {
	return e.cause
}
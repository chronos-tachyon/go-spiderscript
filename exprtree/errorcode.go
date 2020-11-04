package exprtree

import (
	"fmt"
)

// ErrorScopeID
// {{{

type ErrorScopeID uint32

func (sid ErrorScopeID) Scope() *ErrorScope {
	return ErrorScopeByID(sid)
}

func (sid ErrorScopeID) IsValid() bool {
	return sid.Scope().IsValid()
}

func (sid ErrorScopeID) Name() string {
	scope := sid.Scope()
	if scope.IsValid() {
		return scope.Name()
	}
	return fmt.Sprintf("[%#08x]", uint32(sid))
}

func (sid ErrorScopeID) String() string {
	return fmt.Sprintf("error scope %q", sid.Name())
}

func (sid ErrorScopeID) GoString() string {
	return fmt.Sprintf("ErrorScopeID(%#x)", uint32(sid))
}

func (sid ErrorScopeID) WithCode(num uint32) ErrorCode {
	cid := ErrorCodeID(num)
	return ErrorCode{sid, cid}
}

// }}}

// ErrorCodeID
// {{{

type ErrorCodeID uint32

func (cid ErrorCodeID) String() string {
	return fmt.Sprintf("ErrorCodeID(%#x)", uint32(cid))
}

func (cid ErrorCodeID) GoString() string {
	return cid.GoString()
}

// }}}

// ErrorCode
// {{{

type ErrorCode struct {
	SID ErrorScopeID
	CID ErrorCodeID
}

func (code ErrorCode) IsZero() bool {
	return code.SID == 0 && code.CID == 0
}

func (code ErrorCode) Scope() *ErrorScope {
	return code.SID.Scope()
}

func (code ErrorCode) IsValid() bool {
	return code.SID.Scope().IsValidCode(code.CID)
}

func (code ErrorCode) ShortName() string {
	if str, ok := code.SID.Scope().CodeName(code.CID); ok {
		return str
	}
	return fmt.Sprintf("[%#08x]", uint32(code.CID))
}

func (code ErrorCode) Name() string {
	buf := takeBuffer()
	defer giveBuffer(buf)

	scope := code.SID.Scope()
	if str, ok := scope.CodeName(code.CID); ok {
		buf.WriteString(scope.Name())
		buf.WriteString("::")
		buf.WriteString(str)
	} else if scope.IsValid() {
		buf.WriteString(scope.Name())
		buf.WriteString("::")
		buf.WriteString(fmt.Sprintf("[%#08x]", uint32(code.CID)))
	} else {
		buf.WriteString(fmt.Sprintf("[%#08x]", uint32(code.SID)))
		buf.WriteString("::")
		buf.WriteString(fmt.Sprintf("[%#08x]", uint32(code.CID)))
	}
	return buf.String()
}

func (code ErrorCode) Description(data map[string]interface{}) (string, bool) {
	if str, ok := code.SID.Scope().CodeDescription(code.CID, data); ok {
		return str, true
	}
	return "", false
}

func (code ErrorCode) String() string {
	return code.Name()
}

func (code ErrorCode) GoString() string {
	return fmt.Sprintf("ErrorCode(%#08x, %#08x)", uint32(code.SID), uint32(code.CID))
}

func (code ErrorCode) As(newSID ErrorScopeID) (ErrorCode, bool) {
	oldSID := code.SID
	oldCID := code.CID
	if newSID == oldSID {
		return ErrorCode{oldSID, oldCID}, true
	} else if newCID, ok := oldSID.Scope().ConvertTo(newSID, oldCID); ok {
		return ErrorCode{newSID, newCID}, true
	} else if newCID, ok := newSID.Scope().ConvertFrom(oldSID, oldCID); ok {
		return ErrorCode{newSID, newCID}, true
	} else {
		return ErrorCode{}, false
	}
}

var _ fmt.Stringer = ErrorCode{}
var _ fmt.GoStringer = ErrorCode{}

// }}}

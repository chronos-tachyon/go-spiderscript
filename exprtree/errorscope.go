package exprtree

import (
	"fmt"
	"sync"
)

var (
	gErrorScopeMutex  sync.Mutex
	gErrorScopeByID   map[ErrorScopeID]*ErrorScope
	gErrorScopeByName map[string]*ErrorScope
)

// ErrorScope
// {{{

type ErrorScopeImpl interface {
	IsValidCode(cid ErrorCodeID) bool
	CodeName(cid ErrorCodeID) (string, bool)
	CodeDescription(cid ErrorCodeID, data map[string]interface{}) (string, bool)
	ConvertTo(sid ErrorScopeID, cid ErrorCodeID) (ErrorCodeID, bool)
	ConvertFrom(sid ErrorScopeID, cid ErrorCodeID) (ErrorCodeID, bool)
}

type ErrorScope struct {
	sid  ErrorScopeID
	name string
	impl ErrorScopeImpl
}

func NewErrorScope(id uint32, name string, impl ErrorScopeImpl) *ErrorScope {
	if impl == nil {
		panic(fmt.Errorf("ErrorScopeImpl is nil"))
	}

	if !reModuleName.MatchString(name) {
		panic(fmt.Errorf("invalid ErrorScope name %q", name))
	}

	sid := ErrorScopeID(id)

	var newScope *ErrorScope

	locked(&gErrorScopeMutex, func() {
		if gErrorScopeByID == nil {
			gErrorScopeByID = make(map[ErrorScopeID]*ErrorScope, 4)
		}

		if gErrorScopeByName == nil {
			gErrorScopeByName = make(map[string]*ErrorScope, 4)
		}

		if oldScope, found := gErrorScopeByID[sid]; found {
			panic(fmt.Errorf("%#v is already registered as %#v", sid, oldScope))
		}

		if oldScope, found := gErrorScopeByName[name]; found {
			panic(fmt.Errorf("name %q is already in use by %#v", name, oldScope))
		}

		newScope = &ErrorScope{
			sid:  sid,
			name: name,
			impl: impl,
		}

		gErrorScopeByID[sid] = newScope
		gErrorScopeByName[name] = newScope
	})

	return newScope
}

func ErrorScopeByID(sid ErrorScopeID) *ErrorScope {
	var scope *ErrorScope
	locked(&gErrorScopeMutex, func() {
		scope = gErrorScopeByID[sid]
	})
	return scope
}

func ErrorScopeByName(name string) *ErrorScope {
	var scope *ErrorScope
	locked(&gErrorScopeMutex, func() {
		scope = gErrorScopeByName[name]
	})
	return scope
}

func (scope *ErrorScope) ID() ErrorScopeID {
	return scope.sid
}

func (scope *ErrorScope) Name() string {
	return scope.name
}

func (scope *ErrorScope) Impl() ErrorScopeImpl {
	return scope.impl
}

func (scope *ErrorScope) IsValid() bool {
	return (scope != nil)
}

func (scope *ErrorScope) String() string {
	if scope == nil {
		return "<nil>"
	}
	return scope.Name()
}

func (scope *ErrorScope) GoString() string {
	if scope == nil {
		return "nil"
	}
	return fmt.Sprintf("ErrorScope(%#08x, %q)", uint32(scope.ID()), scope.Name())
}

func (scope *ErrorScope) IsValidCode(cid ErrorCodeID) bool {
	if scope == nil {
		return false
	}
	return scope.Impl().IsValidCode(cid)
}

func (scope *ErrorScope) CodeName(cid ErrorCodeID) (string, bool) {
	if scope == nil {
		return "", false
	}
	return scope.Impl().CodeName(cid)
}

func (scope *ErrorScope) CodeDescription(cid ErrorCodeID, data map[string]interface{}) (string, bool) {
	if scope == nil {
		return "", false
	}
	return scope.Impl().CodeDescription(cid, data)
}

func (scope *ErrorScope) ConvertTo(sid ErrorScopeID, cid ErrorCodeID) (ErrorCodeID, bool) {
	if scope == nil {
		return 0, false
	}
	if sid == scope.sid {
		return cid, true
	}
	return scope.Impl().ConvertTo(sid, cid)
}

func (scope *ErrorScope) ConvertFrom(sid ErrorScopeID, cid ErrorCodeID) (ErrorCodeID, bool) {
	if scope == nil {
		return 0, false
	}
	if sid == scope.sid {
		return cid, true
	}
	return scope.Impl().ConvertFrom(sid, cid)
}

var _ fmt.Stringer = (*ErrorScope)(nil)
var _ fmt.GoStringer = (*ErrorScope)(nil)

// }}}

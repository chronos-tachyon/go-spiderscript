package exprtree

import (
	"fmt"
	"sort"
	"sync"
)

// StackTrace
// {{{

type StackTrace []StackTraceFrame

// StackTraceFrame
// {{{

type StackTraceFrame struct {
	f *Function
}

func (fr *StackTraceFrame) Function() *Function {
	return fr.f
}

// }}}

// }}}

// Error
// {{{

type Error struct {
	mu    sync.RWMutex
	id    ErrorID
	code  ErrorCode
	trace StackTrace
	data  map[string]interface{}
}

func (interp *Interp) NewError() *Error {
	id := interp.allocateError()
	err := &Error{id: id}
	interp.registerError(err)
	return err
}

func (err *Error) Clear() {
	locked(&err.mu, func() {
		err.code = ErrorCode{}
		err.trace = nil
		for key := range err.data {
			delete(err.data, key)
		}
	})
}

func (err *Error) Clone() *Error {
	var dupe *Error
	locked(&err.mu, func() {
		dupe = &Error{
			code:  err.code,
			trace: err.trace,
			data:  nil,
		}
		if err.data != nil {
			dupe.data = make(map[string]interface{}, len(err.data))
			for key, value := range err.data {
				dupe.data[key] = value
			}
		}
	})
	return dupe
}

func (err *Error) ID() ErrorID {
	return err.id
}

func (err *Error) Code() ErrorCode {
	var code ErrorCode
	locked(&err.mu, func() {
		code = err.code
	})
	return code
}

func (err *Error) Trace() StackTrace {
	var trace StackTrace
	locked(&err.mu, func() {
		trace = err.trace
	})
	return trace
}

func (err *Error) Data(out map[string]interface{}) {
	locked(&err.mu, func() {
		for key, value := range err.data {
			out[key] = value
		}
	})
}

func (err *Error) Keys() []string {
	var keys []string
	locked(&err.mu, func() {
		keys = make([]string, 0, len(err.data))
		for key := range err.data {
			keys = append(keys, key)
		}
	})
	sort.Strings(keys)
	return keys
}

func (err *Error) GetKey(key string) (interface{}, bool) {
	var value interface{}
	var found bool
	locked(&err.mu, func() {
		value, found = err.data[key]
	})
	return value, found
}

func (err *Error) GetKeyOrNil(key string) interface{} {
	value, _ := err.GetKey(key)
	return value
}

func (err *Error) MustGetKey(key string) interface{} {
	value, found := err.GetKey(key)
	if !found {
		panic(fmt.Errorf("BUG: (*Error).GetKey(%q) returned found=false", key))
	}
	return value
}

func (err *Error) SetCode(code ErrorCode) {
	locked(&err.mu, func() {
		err.code = code
	})
}

func (err *Error) SetTrace(trace StackTrace) {
	locked(&err.mu, func() {
		err.trace = trace
	})
}

func (err *Error) SetKey(key string, value interface{}) {
	locked(&err.mu, func() {
		if err.data == nil {
			err.data = make(map[string]interface{}, 8)
		}
		err.data[key] = value
	})
}

func (err *Error) DeleteKey(key string) {
	locked(&err.mu, func() {
		if err.data == nil {
			err.data = make(map[string]interface{}, 8)
		}
		delete(err.data, key)
	})
}

func (err *Error) WithCode(code ErrorCode) *Error {
	err.SetCode(code)
	return err
}

func (err *Error) WithTrace(trace StackTrace) *Error {
	err.SetTrace(trace)
	return err
}

func (err *Error) WithKey(key string, value Value) *Error {
	err.SetKey(key, value)
	return err
}

func (err *Error) WithKeys(data map[string]interface{}) *Error {
	locked(&err.mu, func() {
		for key, value := range data {
			err.data[key] = value
		}
	})
	return err
}

func (err *Error) WithoutKeys(keys []string) *Error {
	locked(&err.mu, func() {
		for _, key := range keys {
			delete(err.data, key)
		}
	})
	return err
}

func (err *Error) As(sid ErrorScopeID) (*Error, bool) {
	var out *Error
	var success bool
	locked(&err.mu, func() {
		if code, ok := err.code.As(sid); ok {
			out = &Error{
				code:  code,
				trace: err.trace,
				data:  nil,
			}
			if err.data != nil {
				out.data = make(map[string]interface{}, len(err.data))
				for key, value := range err.data {
					out.data[key] = value
				}
			}
			success = true
		}
	})
	return out, success
}

func (err *Error) Error() string {
	buf := takeBuffer()
	defer giveBuffer(buf)

	locked(&err.mu, func() {
		buf.WriteString(err.code.Name())
		if description, ok := err.code.Description(err.data); ok {
			buf.WriteString(": ")
			buf.WriteString(description)
		}
	})

	return buf.String()
}

var _ error = (*Error)(nil)

// }}}

// DuplicateModuleError
// {{{

type DuplicateModuleError struct {
	Name string
	Old  *Module
	New  *Module
}

func (err *DuplicateModuleError) Error() string {
	return fmt.Sprintf("duplicate module %q: old %q, new %q", err.Name, err.Old.CanonicalName(), err.New.CanonicalName())
}

var _ error = (*DuplicateModuleError)(nil)

// }}}

// DuplicateSymbolError
// {{{

type DuplicateSymbolError struct {
	Name string
	Old  *Symbol
	New  *Symbol
}

func (err *DuplicateSymbolError) Error() string {
	return fmt.Sprintf("duplicate symbol %q: old %q, new %q", err.Name, err.Old.CanonicalName(), err.New.CanonicalName())
}

var _ error = (*DuplicateSymbolError)(nil)

// }}}

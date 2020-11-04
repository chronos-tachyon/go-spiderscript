package value

import (
	"fmt"
)

func ExtractLiteral(in interface{}, defaultText string) (out string, err error) {
	out = defaultText
	switch x := in.(type) {
	case *Literal:
		out = x.Text
	case *Error:
		err = x.Err
	default:
		panic(fmt.Errorf("expected *value.Literal or *value.Error, got %T", in))
	}
	return
}

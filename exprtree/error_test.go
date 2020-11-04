package exprtree

import (
	"testing"
)

func TestDuplicateModuleError(t *testing.T) {
	fakeModule1 := &Module{cname: "fake::module::one"}
	fakeModule2 := &Module{cname: "fake::module::two"}
	err := &DuplicateModuleError{Name: "alias", Old: fakeModule1, New: fakeModule2}
	expect := "duplicate module \"alias\": old \"fake::module::one\", new \"fake::module::two\""
	if actual := err.Error(); actual != expect {
		t.Errorf("DuplicateModuleError.Error(): expect %q, actual %q", expect, actual)
	}
}

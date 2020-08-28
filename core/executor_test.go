package core

import (
	"testing"
)

func TestExecutor_parse(t *testing.T) {
	e := &Executor{}
	e.parse(`(a or b) or c`, nil)

}

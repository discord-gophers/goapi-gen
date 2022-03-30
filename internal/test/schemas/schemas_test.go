package schemas

import (
	"testing"
)

func TestForCorrectCustomTypes(t *testing.T) {

	_ = CustomGoTypeWithAlias("example")
	_ = CustomGoType("example")

}

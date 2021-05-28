package subtitle

import "testing"

func PrepareStringTest(t *testing.T) {
	source := "   Hi !ya   how r u   ;   "
	dest := prepareString(source)
	want := "Hi! ya how r u;"
	if want != dest {
		t.Fatalf("prepareString() Failed: want %q, have %q", want, dest)
	}
}

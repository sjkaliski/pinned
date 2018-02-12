package pinned

import (
	"testing"
)

func TestVersionManagerAdd(t *testing.T) {
	vm := &VersionManager{}

	// Should fail if invalid date supplied.
	err := vm.Add(&Version{
		Date: "baddate",
	})
	if err == nil {
		t.Fatal("Expected error when adding version with invalid date")
	}

	// Versions should be sorted in descending order.
	vm.Add(&Version{
		Date: "2016-01-02",
	})
	vm.Add(&Version{
		Date: "2017-01-02",
	})
	versions := vm.Versions()
	if versions[0] != "2017-01-02" || versions[1] != "2016-01-02" {
		t.Fatal("Versions not sorted properly")
	}
}

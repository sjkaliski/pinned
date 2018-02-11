package pinned_test

import (
	"github.com/sjkaliski/pinned"
)

var (
	vm      = pinned.VersionManager{}
	version = &pinned.Version{}
)

func ExampleVersionManager() {
	vm = pinned.VersionManager{
		Header: "Example API",
	}
}

func ExampleVersion() {
	// userNameFieldChange is an action which "undoes" the change
	// made in this new version.
	userNameFieldChange := func(mapping map[string]interface{}) map[string]interface{} {
		mapping["full_name"] = mapping["name"]
		delete(mapping, "name")
		return mapping
	}

	// Version sets the date and lists a set of changes to be executed to
	// make the current API compatible with requests pinned to this version.
	version = &pinned.Version{
		Date: "2018-02-11",
		Changes: []*pinned.Change{
			&pinned.Change{
				Description: "Renames `user.full_name` to `user.name`",
				// Actions are executed by type. In this example "User"
				// references any object of type User.
				Actions: map[string]pinned.Action{
					"User": userNameFieldChange,
				},
			},
		},
	}
}

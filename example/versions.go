package main

import (
	"github.com/sjkaliski/pinned"
)

var (
	vm = pinned.VersionManager{}
)

func init() {
	vm.Add(&pinned.Version{
		Date: "2018-03-09",
		Changes: []*pinned.Change{
			{
				Description: "Renames `user.full_name` to `user.name`",
				Actions: map[string]pinned.Action{
					"User": userNameFieldChange,
				},
			},
		},
	})

	vm.Add(&pinned.Version{
		Date: "2018-02-09",
		Changes: []*pinned.Change{
			{
				Description: "Added `user.created_at`",
				Actions: map[string]pinned.Action{
					"User": userCreatedAtChange,
				},
			},
		},
	})

	vm.Add(&pinned.Version{
		Date: "2018-01-09",
	})
}

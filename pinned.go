package pinned

import (
	"errors"
	"net/http"
	"reflect"
	"sort"
	"time"
)

var (
	// ErrInvalidVersion means the version supplied is not included
	// in the versions supplied to the VersionManager or it is malformed.
	ErrInvalidVersion = errors.New("invalid version")

	// ErrNoVersionSupplied means no version was supplied.
	ErrNoVersionSupplied = errors.New("no version supplied")

	// ErrVersionDeprecated means the version is deprecated.
	ErrVersionDeprecated = errors.New("version is deprecated")
)

const (
	defaultLayout = "2006-01-02"
	defaultHeader = "Version"
	defaultQuery  = "v"
)

// VersionManager represents a list of versions.
type VersionManager struct {
	Layout string
	Header string
	Query  string

	versions []*Version
}

// Add adds a version.
func (vm *VersionManager) Add(v *Version) error {
	var err error
	v.layout = vm.layout()
	v.date, err = time.Parse(vm.layout(), v.Date)
	if err != nil {
		return err
	}
	vm.versions = append(vm.versions, v)

	// Ensure the versions are in descending order by date.
	sort.Sort(sort.Reverse(versions(vm.versions)))

	return nil
}

// Latest returns the most current active version.
func (vm *VersionManager) Latest() *Version {
	if len(vm.versions) == 0 {
		return nil
	}

	return vm.versions[0]
}

// Parse evaluates an http.Request object to determine an API version.
// It inspects the query parameters and request headers. Whichever
// is most recent is the version to use.
func (vm *VersionManager) Parse(r *http.Request) (*Version, error) {
	h := r.Header.Get(vm.header())
	q := r.URL.Query().Get(vm.query())

	if h == "" && q == "" {
		return nil, ErrNoVersionSupplied
	}

	hDate, qDate := time.Time{}, time.Time{}

	var err error
	if h != "" {
		hDate, err = time.Parse(vm.layout(), h)
		if err != nil {
			return nil, ErrInvalidVersion
		}
	}
	if q != "" {
		qDate, err = time.Parse(vm.layout(), q)
		if err != nil {
			return nil, ErrInvalidVersion
		}
	}

	t := hDate
	if hDate.Before(qDate) {
		t = qDate
	}

	v, err := vm.getVersionByTime(t)
	if err != nil {
		return nil, err
	}
	if v.Deprecated {
		return v, ErrVersionDeprecated
	}
	return v, nil
}

func (vm *VersionManager) getVersionByTime(t time.Time) (*Version, error) {
	for _, v := range vm.versions {
		if t.Equal(v.date) {
			return v, nil
		}
	}

	return nil, ErrInvalidVersion
}

// Apply processes a Versionable object by applying all changes between the
// latest version and the version requested. The altered object is returned.
//
// Concretely, if the supplied version is two versions behind the latest, the changes
// in those two versions are applied sequentially to the object. This essentially
// "undoes" the changes made to the API so that the object is structured according to
// the specified version.
func (vm *VersionManager) Apply(version *Version, obj Versionable) (map[string]interface{}, error) {
	data := obj.Data()

	for _, v := range vm.versions {
		// If the requested version is >= to the version, do not apply.
		if version.date.After(v.date) || version.date.Equal(v.date) {
			break
		}

		// Iterate through each change and execute
		// actions as appropriate.
		for _, c := range v.Changes {
			typ := reflect.TypeOf(obj).Elem().Name()

			// If there is an action for this obj type
			// execute the action.
			a, ok := c.Actions[typ]
			if ok {
				data = a(data)
			}
		}
	}

	return data, nil
}

// Versions returns a list of all versions as strings.
func (vm *VersionManager) Versions() []string {
	versions := make([]string, len(vm.versions))
	for i := range versions {
		versions[i] = vm.versions[i].String()
	}

	return versions
}

func (vm *VersionManager) layout() string {
	if vm.Layout != "" {
		return vm.Layout
	}

	return defaultLayout
}

func (vm *VersionManager) header() string {
	if vm.Header != "" {
		return vm.Header
	}

	return defaultHeader
}

func (vm *VersionManager) query() string {
	if vm.Query != "" {
		return vm.Query
	}

	return defaultQuery
}

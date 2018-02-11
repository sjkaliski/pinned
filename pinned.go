package pinned

import (
	"errors"
	"net/http"
	"reflect"
	"sort"
	"time"
)

// Errors.
var (
	ErrInvalidVersion    = errors.New("invalid version")
	ErrNoVersionSupplied = errors.New("no version supplied")
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
		return &Version{}
	}

	return vm.versions[0]
}

// Parse evaluates an http.Request object to
// determine an API version. It inspects the query parameters
// and request headers. Whichever is most recent wins.
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

	return vm.getVersionByTime(t)
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
func (vm *VersionManager) Apply(version *Version, obj Versionable) (map[string]interface{}, error) {
	data := obj.Data()

	for _, v := range vm.versions {
		// If the requested version is >= to the version, do not apply.
		if version.date.After(v.date) || version.date.Equal(v.date) {
			break
		}

		// Otherwise, apply changes as necessary.
		for _, c := range v.Changes {
			typ := reflect.TypeOf(obj).Elem().Name()
			f, ok := c.Actions[typ]
			if ok {
				data = f(data)
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

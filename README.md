# pinned

This is a proof-of-concept, date based versioning system for JSON APIs inspired by [Stripe's API versioning](https://stripe.com/blog/api-versioning).

## Overview

This package enables reverse compatibility for a Go API by defining versions and their associated changes. Consequently, API versions can be maintained for long periods of time without minimal effort.

## Example

Consider the following simple example. An API has a User struct which looks like this:

```go
type User struct {
  ID       uint64
  FullName string
}
```

Now we decide we want to rename `FullName` to `Name`. However, this is a breaking change. To ensure stability, prior to change we set a version `2018-02-10`.

After the change, we set a version `2018-02-11`. This version has a change associated with it. This `Change` has an `Action` to be taken on the `User` resource.

This `Action` is a `func` that _reverses the change_ made in the new version.

```go
func userNameFieldChange(mapping map[string]interface{}) map[string]interface{} {
  mapping["full_name"] = mapping["name"]
  delete(mapping, "name")
  return mapping
}
```

There are now two versions, `2018-02-11` and `2018-02-10`. To support the client that requested version `2018-02-10`, the "changes" made in version `2018-02-11` are undone, and the User resource now reflects the requested version.

As versions are added, these changes are sequentially undone. This enables a version to be supported for a long period of time, and allows the developer to focus on new feature development without much concern towards legacy versions.

## Usage

See the included [example](/example) project for sample usage.

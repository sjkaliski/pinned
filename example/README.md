# example

This is an example API that uses `pinned` to manage versions.

In this example, the API maintains three versions.

```
2018-03-09 - User.FullName -> User.Name
2018-02-09 - Remove User.CreatedAt
2018-01-09 - Initial.
```

If a request does not specify a version, or requests the latest, no changes are made. If a request is made for `2018-02-09`, the changes defined in `2018-03-09` are executed. If a request is made for `2018-01-09`, the changes defined in `2018-03-09` and `2018-02-09` are executed in order.

To see what this actually looks like, look at [__snapshots__/example.snapshot](__snapshots__/example.snapshot).

## Testing

Tests are written using [abide](https://github.com/beme/abide). This makes it very simple to snapshot each version of a given API route and ensure future changes keep previous versions stable.

```
$ go test -v -race
```

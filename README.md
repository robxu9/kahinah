# Kahinah, the karma manager

Kahinah is an advisory and update manager for packages. It creates advisories
based on updates (and their dependencies), and maintains connectors with
build systems such as [ABF](https://abf.io).

Kahinah has two components: the library component, and the server component.

The library manages all access to the database, and controls the route of
updates, while the server provides endpoints for clients to access the library.

## Building

Just using `go get` should be enough to get started:

* `go get github.com/robxu9/kahinah/server1` should build the server binary.

## License

Kahinah is licensed under the [MIT license](http://robxu9.mit-license.org).

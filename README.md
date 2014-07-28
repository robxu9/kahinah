# Kahinah v4: less is more

Kahinah _v4_ is a flexible QA system. It provides a system of karma for testing newly built packages, maintains the advisories system, and provides virtual testing of ISOs provided to it.

_v4_ is a rewrite from previous versions: less is more.

## rewrite? what are you thinking

The build time of Kahinah is grossly horrendous now because we pull in too much. Our framework, beego, is nice; but we don't need to do as much as we're making it do.

This rewrite focuses on less spaghetti code and more simplicity; namely, we focus on the following packages only so that we don't bloat the system:
* mux, for routing
* context, for storing variables when chaining calls
* sessions, for providing session management via secure cookies
* negroni, for a decent not-martini-hell middleware, which also nets us data verification via [binding](https://github.com/mholt/binding)
* [render](https://github.com/unrolled/render), which lets us easily render JSON and HTTP templates
* [modl](https://github.com/jmoiron/modl), because I'm lazy and ORMs are still cool
* toml, because I can't deal with JSON or YAML for configuration formats right now and this seems slightly decent

## what I am trying to accomplish

A cleaner codebase. A more manageable. Less "where does everything go" and more "enjoy coding".

Isolating each handler - to each, their own.

Also, making a stable base for an API; right now, the API can't be done because I sure as hell don't isolate the logic. That's something that will be done - each folder contains a part of the system and its logic; and the handlers simply call said functions and then return it in a proper form.

## requirements

* go 1.3+
* git & mktemp -d [for diff generation]
* a decent database

## running

When it's up, you need to copy config.toml.example to config.toml and edit it. Figure out database parameters by looking at the respective driver pages.

Then start it.

## hacking

Use the `Makefile` - just `make` to run it, `make minify` to minify all the static files, `make dist` to generate a tarball.

## struggling

Isolating the logic in _v4_ is the main task here: I'm trying to make this a maintainable system. It's also extremely difficult hooking into external systems; it's best to abstract that and/or have multiple providers in order to be the most flexible system.

## explaining

Kahinah handles:

* newly released builds: the workflow is that Kahinah, once an hour, will check external systems for new builds, or will be pinged to check new builds by external systems (if there is support for that) (this checker will live in the builds package). If it finds any new packages that have just finished building, it will push them to testing (if necessary), create a record for that and allow voting to that package, generating a git diff between the last published version and the one specified. Votes are either `+1`, `-1`, `0` (so just comment only), `+3`, `-3` (maintainer votes), `+9999`, or `-9999` (overriding votes). If a package is in testing, Kahinah will make sure to check if it has an already existing record for it. If a package has just been rejected or published on the external system side, Kahinah will update its record of that package, creating a record if necessary, and will use overriding votes to indicate the status. This way Kahinah stays in sync with the external system(s).

* advisories: more for grouping updates together, this assigns a permanent codename and number to the updates specified (like DIST-YEAR-NUMBER). All updates must either be in the testing state or the published state (so > -3). If any update falls afoul, the advisory is redacted automatically and will be noted, and the creator of the advisory will be notified that their advisory has been redacted due to a failed update. Once an advisory is issued (all packages move to a published state), it can no longer be edited. An advisory has three states: pending (packages still being tested), issued (packages passed qa and are published), and redacted (some packages failed QA, advisory rendered bad). It contains a short summary and a longer description, a type (security, bugfixes, enhancements, recommended, other), and a list of bugs that are fixed by the advisory.

* virtual testing: not really implemented it, but it's anticipated that it'll hook into uitest & libvirt-go to provide a consistent platform for which to test ISOs on.

## derping

I derp a lot. Bugs to the [Github](https://github.com/robxu9/kahinah), pls.

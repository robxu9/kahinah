## Linking Updates to Advisories

Linking updates to advisories can be tricky. Usually, there is no good way of
doing it and we let maintainers link it manually. But in projects where time
is essentially wasted doing this, there needs to be a better way to do this.

We can achieve this by having a file of rules which Kahinah interprets. It is a
JSON file with sections of rules and regex values. It's better to show it than
to explain it, so here's how it would work:

```json
[
    {
        "_comment": "Comment for the ruleset; ignored",
        "rule_name": "name_of_rule",
        "name": "regex_to_use",
        "version": "regex_to_use",
        "files": "regex_to_use",
        "requires": "regex_to_use"
    }
]

```

Get the point? You can use any regex that Go supports, so it would probably be a
good idea to [read up on it](https://golang.org/pkg/regexp/). It's not that
different than most languages. A non-existing regex for a field will allow
everything to be satisfied (so we won't check), while a blank regex will
probably reject everything.

Substitutions can be used as well. For example, if we have a package called
`foo`, and we also want it to catch `foo-plugins` and `foo-more-awesome`, then
we can add a regex like the following:

```json
[
    {
        "_comment": "catch -* from packages",
        "rule_name": "dashed_inc",
        "name": "@NAME@-.*"
    }
]
```

`@NAME@` will be replaced by the package name, respectively. Same goes for
`@VERSION@`. If `@NAME@` needs to be modified (or `@VERSION@`), you can
use `name_cut` and `version_cut` respectively to modify those macros (regex
matched with `*_cut` will be replaced with empty strings).

If there is more than one rule that satisfies the package, the __last__ one
wins. Remember that groups are interpreted from top to bottom, and files are
interpreted in natural order (so it would probably be a good idea to prefix your
files with numbers, e.g. `00-default.json` to `99-custom.json`).

## What exactly are you matching to?

We're matching to _current_ advisories. So, assume I have no advisories,
and I have `pkga`, `pkgb-cool`, `pkgb-cooler`, `pkgb`, and `pkgb-other` that
need advisories created (in that order).

If we have the rules defined here (top-to-bottom interpreted, as always):

```json
[
    {
        "_comment": "Rule 1: connect with @NAME@-* packages",
        "rule_name": "dashed_matching",
        "name": "@NAME@-.*"
    },
    {
        "_comment": "Rule 2: connect @NAME@-* with @NAME@",
        "rule_name": "reversed_dashed_matching",
        "name_cut": "-.*",
        "name": "(?m)@NAME@((?=-).*)*$"
    }
]
```

* `pkga` will not match Rule 1, because there are no existing advisories that
contain `pkga-.*`, and will not match Rule 2, because no existing advisories
that contain `pkga`. So it matches nothing, and a new advisory will be created
for it.
* `pkgb-cool` will not match Rule 1, because there are no existing advisories
that contain `pkgb-cool-.*`, and will not match Rule 2, because there are no
existing advisories that contain `pkgb`. So it matches nothing, and a new
advisory will be created for it.
* `pkgb-cooler` will not match Rule 1, because there are no existing advisories
that contain `pkgb-cooler-.*`, but WILL match Rule 2 because it matches that
there is a package `pkgb-cool` that has `pkgb-` in it. So it will attach onto
`pkgb-cool`'s advisory.
* `pkgb` WILL match Rule 1 with `pkgb-cool`, because it matches `pkgb-.*`, but
it will not match Rule 2 because it does not match `pkgb`. So it will attach
onto `pkgb-cool`'s advisory.
* `pkgb-other` will not match Rule 1 because there are no existing Advisories
that contain `pkgb-other-.*`, but WILL match Rule 2 because it does match
an advisory with an updated `pkgb`. So it will attach onto `pkgb-cool`'s
advisory.

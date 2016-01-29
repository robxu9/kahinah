Kahinah
=======

Kahinah is a QA system. Upon receiving notification of a built package, it runs
a set of tasks (usually determined from where it came from) and then displays it
for users to upvote or downvote until it can be published or rejected.

Requirements
------------

* Git
* mktemp -d
* Go 1.5+
* sqlite3 or mysql

Setup
-----

1. Configure your application in app.toml (see app.toml.example for details).
1. Configure news.md.example.
1. Make sure your databases are created, if needed.
1. Run it.

Tips
----

* Make sure to run `go clean -i -r -x` before every build.
    * Clear out old libraries.
* Update the source every time with `go get -u`.
    * If you get errors, the above cleaning command helps.

Administration
--------------

To be able to control users and permissions, you need to have an admin account.

1. Login with your account to Persona. This will create an account in the database.
1. Add your email address to adminwhitelist (you may add several emails, seperated by semicolons)
1. Restart Kahinah and navigate to /admin. Add the kahinah.admin permission to your users.

It's recommended that you remove your email addresses from adminwhitelist afterwards, and replace them with impossible ones (e.g. 1y8oehowhfoinwdaf).


Whitelist
---------

Kahinah, by default, allows everyone to vote up/down builds. This may have unintended effects.

Currently, there is a whitelist system embedded - set "whitelist" in app.conf to true to use it.

This allows only users in the whitelist to vote. It also adds a permanent notice to the front page
that indicates that it is only available for whitelisted users.

Administrators should use the above administration tool to add the "kahinah.whitelist" permission
to users that should be allowed to vote.

Database
--------

`db_type` is either `sqlite3` or `mysql`. For `sqlite3`, the `db_name` is the filename of the database.

The rest should be self-explainatory. Hopefully.

No major database changes are expected, I think. Maybe we'll introduce migrations with hood or goose
when the time comes to change the database schema, but I don't anticipate major changes anytime soon.

News
----

Markdown.

License
-------

As of v4, MIT licensed. See [the LICENSE file](http://robxu9.mit-license.org)
for more information.

Previous versions were licensed under the AGPLv3.

Kahinah
=======

Kahinah is the QA system for OpenMandriva. It hooks into ABF to retrieve packages that have been recently built and automates the QA portion of the package process (pushing to testing, karma, pushing to updates).

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
## Advisories

Advisories can be considered a group of updates that should be evaluated
together. They can be voted on (+2, +1, 0, -1, -2) and are automatically
adjusted based on updates that arrive via the system.

Advisories usually contain the following properties:

* Advisory ID, usually [GROUP]-[YEAR]:[ADVISORY_NUM].[REV] [^advisory_id] [^db]
* Date
* Description
* List of References [^references]

And also include the following properties, which are mutable:

* List of Updates
* List of Verdicts (similar to Gerrit, {Type, + or - how much}) [^verdict]
* List of Comments [^verdict]

It can generate the following properties from the list of updates:

* Platforms Affected
* List of New Packages
* Pushed (whether an advisory has pushed its updates from testing)
* Closed (whether an advisory is available for further modification)

[^advisory_id]: The ID components are structured typically in this way based on
observations from other distribution advisories. In turn, the _GROUP_ refers to
the issuing product identifier (usually the distribution name, like `FEDORA`
or `SUSE`), followed by the year that the advisory is issued in (e.g. 2015). The
advisory number is the sequential number that corresponds to that specific
advisory _in that specific year_ (so the first advisory of 2015 will have number
1, and will sequentially increment for new advisories). Advisory numbers do not
have to move into production if the advisory fails. Revisions are optional, but
are self-explanatory.

[^db]: the database might represent this a bit differently, by compounding the keys with each component of the advisory ID (e.g. 'group', 'year', 'advisory_num', and 'rev' are all indexes)

[^references]: usually CVE links or other distribution links

[^verdict]: for all intents and purposes, a comment is a neutral verdict.

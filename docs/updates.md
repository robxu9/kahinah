## Updates

Updates are packages that have been submitted to Kahinah for inclusion into
advisories. They are immutable to all except connectors, which typically add the
updates into the system (and modify them as well, in case new updates come in).

Updates usually have the following properties:

* Targeted Platforms
* List of New Packages
* Date
* Name and EVR
* Submitter
* Update Diff (changes between updates)
* Changelog (written changelog)
* Type of Update (security, bug fix, recommended, new package)
* Connector ID [^connector]

Additionally, the following properties are usually tampered with by the
connector:

* Advisory ID (linking an update to an advisory)
* Deprecated (if another update comes though and this one has not been pushed)
* Pushed (this update was pushed through successfully via an advisory)
* Connector Info, in which connectors can put whatever they want here
[^connector]

[^connector]: Sometimes, connectors need to store extra information (such as
the id on the build system side, urls, and more). This is the place to do so.

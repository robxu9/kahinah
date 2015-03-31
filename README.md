# Kahinah: the karma manager

This rewrite carries the same principle as the last one, except
that updates aren't automatically votable. Instead, updates are
immutable, and as they are automatically imported by connectors,
they must be submitted for testing and comment via advisories.

Advisories brings Kahinah in line with standards, at least
where managing package updates are concerned.

A broken mess that was, before.

### Regarding ABF

I still poll for ABF every hour. ABF provides no known form of
push notifications (not that I'm expecting them to) and I'm still
aways off from connecting to an email client and reading email
notifications, so you'll have to bear with me. I may decrease the
polling limit to every 5-10 minutes, depending on if I can get
alternate approval from the administrators (not a high priority
right now though, truth be told).

When updates hit "Build Completed", they will make their way to
the updates section of Kahinah, but WILL NOT get sent to testing
automatically (breaking change!). You need to create an advisory
and fill out the details (that's not optional) and then Kahinah
will send all those updates to testing.

### Developing
go1.4 is currently the targeted version.

Do a `vagrant up` if you want to develop. It's much easier.

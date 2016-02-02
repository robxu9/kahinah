Hello,

Recent activity has occurred on an update you are watching/participating in:

| Update ID | {{.Update.Id}} |
|----------:|----------------|
|      Name | {{.Update.Name}} |
|    Target | {{.Update.Platform}}/{{.Update.Repo}}/{{.Update.Architecture}} |
|      Type | {{.Update.Type}} |
|     Built | {{.Update.BuildDate}} |
|    Status | {{.Update.Status}} |

Here are the list of activities that have occurred:

| Person | did Action | with Comment |
|-------:|------------|--------------|
{{range .Update.Karma}}
| {{if .User.Username}}{{.User.Username}}{{else}}{{.User.Email}}{{end}} | {{.Vote}} | {{.Comment}} |
{{end}}

For more information, check out [the update on Kahinah]({{.URL}}/builds/{{.Update.Id}}).

Cheers,
The Kahinah Bot

-------------------------------
This email was sent by Kahinah, the OpenMandriva QA bot.
Inbound email to this account is not monitored.

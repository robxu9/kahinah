<div id="app">
    <section class="section">
        <div class="container">
            <h1 class="title">{{.Title}}</h1>
            <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Target</th>
                        <th>Type</th>
                        <th>Status</th>
                        <th>Advisory</th>
                        <th>Updated</th>
                    </tr>
                </thead>
                <tbody>
                    {{with .Packages}}
                    {{range .}}
                    <tr>
                        <td class="table-link"><a href="{{urldata "/b/{{.Id}}" .}}">{{.Id}}</a></td>
                        <td class="table-link"><a href="{{urldata "/b/{{.Id}}" .}}">{{.Name}}</a></td>
                        <td>{{.Platform}}: {{.Repo}} ({{.Architecture}})</td>
                        <td class="table-icon"><i class="fa {{if eq .Type "bugfix"}}fa-bug{{end}}{{if eq .Type "security"}}fa-shield{{end}}{{if eq .Type "enhancement"}}fa-gift{{end}}{{if eq .Type "recommended"}}fa-star{{end}}{{if eq .Type "newpackage"}}fa-plus-square{{end}}" title="{{.Type}}"></i></td>
                        <td><img src="{{if eq .Status "testing"}}//b.repl.ca/v1/status-TESTING-yellow.png{{else}}
                            {{if eq .Status "rejected"}}//b.repl.ca/v1/status-REJECTED-red.png{{else}}
                            {{if eq .Status "published"}}//b.repl.ca/v1/status-PUBLISHED-brightgreen.png{{else}}
                            //b.repl.ca/v1/status-UNKNOWN-lightgrey.png{{end}}{{end}}{{end}}" alt="{{.Status}}"></td>
                        <td class="table-link">{{if .Advisory}}<a href="{{urldata "/a/{{.Advisory.Id}}" .}}">{{.Advisory.Id}}</a>{{end}}</td>
                        <td v-text="&quot;{{.Updated | rfc3339}}&quot; | moment &quot;from&quot;"></td>
                    </tr>
                    {{end}}
                    {{end}}
                </tbody>
            </table>
            <br />
            <nav class="navbar">
                <div class="navbar-item is-centered">
                    <form name="input" method="get">
                        <p class="control is-grouped">
                            <a class="button" href="?page={{.PrevPage}}">&larr;</a>
                            <input class="input" name="page" type="text" placeholder="{{.Page}}/{{.Pages}}">
                            <a class="button" href="?page={{.NextPage}}">&rarr;</a>
                        </p>
                    </form>
                </div>
            </nav>
        </div>
    </section>
</div>

<script>
    new Vue({
        el: '#app'
    });
</script>

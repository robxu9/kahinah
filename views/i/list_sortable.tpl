<!-- tablesorter (http://mottie.github.io/tablesorter/) -->
<link rel="stylesheet" href="{{url "/static/css/jquery.tablesorter.pager.min.css"}}">
<!--<link rel="stylesheet" href="{{url "/static/css/theme.default.min.css"}}">-->
<script src="{{url "/static/js/jquery.metadata.min.js"}}"></script>
<script src="{{url "/static/js/jquery.tablesorter.combined.min.js"}}"></script>
<script src="{{url "/static/js/jquery.tablesorter.pager.min.js"}}"></script>

<div id="app">
    <section class="section">
        <div class="container">
            <h1 class="title">{{.Title}}</h1>
            <h2 class="subtitle">{{.Entries}} {{if eq .Entries 1}}entry{{else}}entries{{end}} returned. Click on the headers to sort.</h2>
            <table class="table tablesorter" id="pkgtable">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Target</th>
                        <th>Type</th>
                        <th>Karma</th>
                        <th>Build Date</th>
                    </tr>
                </thead>
                <tbody>
                    {{$out := .}}
                    {{with .Packages}}
                    {{range .}}
                    <tr>
                        <td class="table-link"><a href="{{urldata "/b/{{.Id}}" .}}">{{.Id}}</a></td>
                        <td class="table-link"><a href="{{urldata "/b/{{.Id}}" .}}">{{.Name}}</a></td>
                        <td>{{.Platform}}: {{.Repo}}</td>
                        <td class="table-icon"><i class="fa {{if eq .Type "bugfix"}}fa-bug{{end}}{{if eq .Type "security"}}fa-shield{{end}}{{if eq .Type "enhancement"}}fa-gift{{end}}{{if eq .Type "recommended"}}fa-star{{end}}{{if eq .Type "newpackage"}}fa-plus-square{{end}}" title="{{.Type}}"></i><div style="display: none;">{{.Type}}</div></td>
                        <td>{{$karma := mapaccess .Id $out.PkgKarma}}<img src="{{if eq $karma "0"}}//b.repl.ca/v1/karma-   {{$karma}}-yellow.png{{else}}{{if lt $karma "0"}}//b.repl.ca/v1/karma-  -{{$karma}}-orange.png{{else}}{{if gt $karma "0"}}//b.repl.ca/v1/karma- +{{$karma}}-yellowgreen.png{{end}}{{end}}{{end}}" alt="{{$karma}}"></td>
                        <td unix="{{.BuildDate.Unix}}" v-text="&quot;{{.BuildDate | rfc3339}}&quot; | moment &quot;from&quot;"></td>
                    </tr>
                    {{end}}
                    {{end}}
                </tbody>
            </table>
        </div>
    </section>
</div>

<script>
    $(document).ready(function() {
        $("#pkgtable").tablesorter({
            //theme: "default",
            headerTemplate: "{content} {icon}",
            widgets: ["uitheme", "filter", "zebra"],
            textExtraction: {
                3: function(node, table, cellIndex) {
                    return $(node).find("div").text();
                },
                4: function(node, table, cellIndex) {
                    return $(node).find("img").attr("alt");
                },
                5: function(node, table, cellIndex) {
                    return $(node).attr("unix");
                }
            },
        });
    });

    new Vue({
        el: '#app'
    });
</script>

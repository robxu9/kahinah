{{template "header.tpl" .}}

      <link rel="stylesheet" href="{{url "/static/css/theme.bootstrap.css"}}">

      <script type="text/javascript" src="{{url "/static/js/jquery.tablesorter.min.js"}}"></script>
      <script type="text/javascript" src="{{url "/static/js/jquery.tablesorter.pager.min.js"}}"></script>
      <script type="text/javascript" src="{{url "/static/js/jquery.tablesorter.widgets.min.js"}}"></script>

      <script>
        // With customizations
        $(document).ready(function() { 
          $.extend($.tablesorter.themes.bootstrap, {
            // these classes are added to the table. To see other table classes available,
            // look here: http://twitter.github.com/bootstrap/base-css.html#tables
            table      : 'table table-bordered',
            caption    : 'caption',
            header     : 'bootstrap-header', // give the header a gradient background
            footerRow  : '',
            footerCells: '',
            icons      : '', // add "icon-white" to make them white; this icon class is added to the <i> in the header
            sortNone   : 'bootstrap-icon-unsorted',
            sortAsc    : 'icon-chevron-up glyphicon glyphicon-chevron-up',     // includes classes for Bootstrap v2 & v3
            sortDesc   : 'icon-chevron-down glyphicon glyphicon-chevron-down', // includes classes for Bootstrap v2 & v3
            active     : '', // applied when column is sorted
            hover      : '', // use custom css here - bootstrap class may not override it
            filterRow  : '', // filter row class
            even       : '', // odd row zebra striping
            odd        : ''  // even row zebra striping
          });

          $("#pkgtable").tablesorter({
            theme: "bootstrap",
            headerTemplate: "{content} {icon}",
            widgets: ["uitheme", "filter", "zebra"],
            textExtraction: {
              3: function(node, table, cellIndex) {
                return $(node).find("div").text();
              },
              4: function(node, table, cellIndex) {
                return $(node).find("img").attr("alt");
              },
            },
          });
        });
      </script>

      <div class="row table-responsive">
        <table class="table tablesorter" id="pkgtable">
          <thead>
            <tr>
              <th>Name</th>
              <th>Submitter</th>
              <th>For</th>
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
                <td><a href="{{urldata "/builds/{{.Id}}" .}}">{{.Name}}/{{.Architecture}}</a></td>
                <td>{{.Submitter.Email | emailat}}</td>
                <td>{{.Platform}}/{{.Repo}}</td>
                <td><i class="fa {{if eq .Type "bugfix"}}fa-bug{{end}}{{if eq .Type "security"}}fa-shield{{end}}{{if eq .Type "enhancement"}}fa-gift{{end}}{{if eq .Type "recommended"}}fa-star{{end}}{{if eq .Type "newpackage"}}fa-plus-square{{end}}" title="{{.Type}}"></i><div style="display: none;">{{.Type}}</div></td>
                <td>{{$karma := mapaccess .Id $out.PkgKarma}}<img src="{{if eq $karma "0"}}//b.repl.ca/v1/karma-   {{$karma}}-yellow.png{{else}}{{if lt $karma "0"}}//b.repl.ca/v1/karma-  -{{$karma}}-orange.png{{else}}{{if gt $karma "0"}}//b.repl.ca/v1/karma- +{{$karma}}-yellowgreen.png{{end}}{{end}}{{end}}" alt="{{$karma}}"></td>
                <td>{{.BuildDate | since}}</td>
              </tr>
              {{end}}
            {{end}}
        </table>
        <center><span class="label label-default">{{.Entries}} {{if eq .Entries 1}}entry{{else}}entries returned.{{end}}</span></center>
      </div>

{{template "footer.tpl" .}}
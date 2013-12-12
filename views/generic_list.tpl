{{template "header.tpl" .}}

      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/addons/pager/jquery.tablesorter.pager.css" />
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/addons/pager/jquery.tablesorter.pager.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/addons/pager/jquery.tablesorter.pager.min.js"></script>
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/filter.formatter.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.black-ice.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.blue.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.bootstrap.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.dark.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.default.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.dropbox.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.green.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.grey.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.ice.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/tablesorter/2.13.3/css/theme.jui.css" />
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.metadata.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.tablesorter.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.tablesorter.min.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.tablesorter.widgets-filter-formatter.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.tablesorter.widgets-filter-formatter.min.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.tablesorter.widgets.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/jquery.tablesorter.widgets.min.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-date-iso8601.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-date-month.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-date-two-digit-year.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-date-weekday.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-date.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-feet-inch-fraction.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-file-type.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-ignore-articles.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-input-select.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-ipv6.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/parsers/parser-metric.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/widgets/widget-build-table.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/widgets/widget-editable.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/widgets/widget-grouping.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/widgets/widget-pager.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/widgets/widget-repeatheaders.js"></script>
      <script type="text/javascript" src="//cdn.jsdelivr.net/tablesorter/2.13.3/js/widgets/widget-scroller.js"></script>

      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/jquery.tablecloth/1.0.0/css/bootstrap-tables.css" />
      <link type="text/css" rel="stylesheet" href="//cdn.jsdelivr.net/jquery.tablecloth/1.0.0/css/tablecloth.css" />
      <script type="text/javascript" src="//cdn.jsdelivr.net/jquery.tablecloth/1.0.0/js/jquery.tablecloth.js"></script>

      <script>
        // With customizations
        $(document).ready(function() { 
          $("#pkgtable").tablecloth({
            sortable: true,
            clean: true,
            cleanElements: "th td",
          });
        });
      </script>

      <div class="row">
        <br/>
        <center><span class="label label-default">{{.Entries}} {{if eq .Entries 1}}entry{{else}}entries returned.{{end}}</span></center>
        <table class="table" id="pkgtable">
          <thead>
            <tr>
              <th>Update ID</th>
              <th>Name</th>
              <th>Submitter</th>
              <th>Platform</th>
              <th>Repository</th>
              <th>Architecture</th>
              <th>Date</th>
            </tr>
          </thead>
          <tbody>
            {{with .Packages}}
              {{range .}}
              <tr>
                <td><a href="/builds/{{.Id}}">{{.Id}}</a></td>
                <td>{{.Name}}</td>
                <td>{{.Submitter.Email | emailat}}</td>
                <td>{{.Platform}}</td>
                <td>{{.Repo}}</td>
                <td>{{.Architecture}}</td>
                <td>{{.BuildDate | since}}</td>
              </tr>
              {{end}}
            {{end}}
        </table>
        <center><span class="label label-default">{{.Entries}} {{if eq .Entries 1}}entry{{else}}entries returned.{{end}}</span></center>
      </div>

{{template "footer.tpl" .}}
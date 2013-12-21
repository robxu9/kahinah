{{template "header.tpl" .}}

      {{if .KarmaControls}}<script src="/static/js/post.js"></script>
      <script>
        function postUp() {
          var params = new Array();
          params['type'] = "Up";
          post_to_url(window.location.href, params, "POST");
        }

        function postDown() {
          var params = new Array();
          params['type'] = "Down";
          post_to_url(window.location.href, params, "POST");
        }

        {{if .MaintainerControls}}
        function postMaintainer() {
          var params = new Array();
          params['type'] = "Maintainer";
          post_to_url(window.location.href, params, "POST");
        }
        {{end}}

        {{if .QAControls}}
        function postQA() {
          if (confirm("This is irreversible! Are you sure you want to block with -9999 karma?")) {
            var params = new Array();
            params['type'] = "QABlock";
            post_to_url(window.location.href, params, "POST");
          }
        }
        {{end}}
      </script>{{end}}

      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          <br/>
          {{if eq .Package.Status "testing"}}<div class="panel panel-warning">{{else}}
          {{if eq .Package.Status "rejected"}}<div class="panel panel-danger">{{else}}
          {{if eq .Package.Status "published"}}<div class="panel panel-success">{{else}}
          <div class="panel panel-primary">{{end}}{{end}}{{end}}
            <div class="panel-heading">
              <h1>{{.Package.Name}} <small>[{{.Package.Architecture}}] OMV-{{.Package.BuildDate.Year}}-{{.Package.Id}}</small><div class="pull-right">{{if .KarmaControls}}<a href="#" class="btn" onclick="postUp()"><i class="fa fa-3x {{if .KarmaUpYes}}fa-thumbs-up{{else}}fa-thumbs-o-up{{end}}"></i></a>{{end}} {{.Karma}} {{if .KarmaControls}}<a href="#" class="btn" onclick="postDown()"><i class="fa fa-3x {{if .KarmaDownYes}}fa-thumbs-down{{else}}fa-thumbs-o-down{{end}}"></i></a>{{if .MaintainerControls}}<a href="#" class="btn" onclick="postMaintainer()"><i class="fa fa-3x {{if .KarmaMaintainerYes}}fa-check-square{{else}}fa-check-square-o{{end}}"></i></a>{{end}}{{if .QAControls}}<a href="#" id="qabtn" title="Instant Reject - WARNING - UNREVERSABLE and -9999 karma!" class="btn" onclick="postQA()"><i class="fa fa-3x fa-sort-amount-desc"></i></a>{{end}}{{end}}</div></h1>
            </div>
            <table class="table table-condensed">
              <tbody>
                <tr>
                  <td><b>Submitter</b></td>
                  <td>{{.Package.Submitter.Email | emailat}}</td>
                </tr>
                <tr>
                  <td><b>Platform<b></td>
                  <td>{{.Package.Platform}}</td>
                </tr>
                <tr>
                  <td><b>Repository<b></td>
                  <td>{{.Package.Repo}}</td>
                </tr>
                <tr>
                  <td><b>Update Type<b></td>
                  <td>
                    {{if eq .Package.Type "bugfix"}}<i class="fa fa-bug"></i>{{end}}
                    {{if eq .Package.Type "security"}}<i class="fa fa-shield"></i>{{end}}
                    {{if eq .Package.Type "enhancement"}}<i class="fa fa-gift"></i>{{end}}
                    {{if eq .Package.Type "recommended"}}<i class="fa fa-star-o"></i>{{end}}
                    {{if eq .Package.Type "newpackage"}}<i class="fa fa-plus-square-o"></i>{{end}}
                    {{.Package.Type}}</td>
                </tr>
                <tr>
                  <td><b>URL<b></td>
                  <td><a href="{{.Url}}">{{.Url}}</a></td>
                </tr>
                <tr>
                  <td><b>Packages<b></td>
                  <td>
                    <table class="table table-bordered table-condensed table-responsive">
                      <tbody>
                        {{with .Package.Packages}}
                          {{range .}}
                        <tr>
                          <td>{{.Name}}-{{if gt .Epoch 0}}{{.Epoch}}:{{end}}{{.Version}}-{{.Release}}.{{.Type}}</td>
                        </tr>
                          {{end}}
                        {{end}}
                      </tbody>
                    </table>
                  </td>
                </tr>
                <tr>
                  <td><b>Build Date</b></td>
                  <td>{{.Package.BuildDate}}</td>
                </tr>
                <tr>
                  <td><b>Last Updated</b></td>
                  <td>{{.Package.Updated}}</td>
                </tr>
            </table>
          </div>
        </div>
      </div>

      <!-- diff & changelog -->
      <div class="row">
        <div class="col-md-10 col-md-offset-1">

          <div class="panel panel-info">
            <div class="panel-heading"><button class="btn btn-info" data-toggle="collapse" href="#diff">Git Diff</button></div>
            <div id="diff" class="panel-collapse collapse">
              <pre class="brush: diff">{{.Package.Diff}}</pre>
            </div>
          </div>

          <div class="panel panel-primary">
            <div class="panel-heading"><button class="btn btn-primary" data-toggle="collapse" href="#cnlog">Changelog</button></div>
            <div id="cnlog" class="panel-collapse collapse">
              <pre class="pre-scrollable">{{if .Changelog}}{{.Changelog}}{{else}}Not Available{{end}}</pre>
            </div>
          </div>

        </div>
      </div>

      <!-- want to test? -->
      {{if eq .Package.Status "testing"}}
      <div class="row">
        <div class="col-md-8 col-md-offset-2">
          <div class="panel panel-info">
            <div class="panel-heading">Want to test?</div>
            <div class="panel-body">
              On your <code>{{.Package.Platform}}/{{.Package.Architecture}}</code> machine, do:
              <pre>kahup {{.Package.Id}}</pre>
              to install the above packages from the testing repository onto your computer.<br/>
              <br/>
              When you're done, you can use <code>urpmi --downgrade</code> to revert back to previous versions.
            </div>
          </div>
        </div>
      </div>
      {{end}}

      <!-- results -->
      <div class="row">
        <div class="col-md-3 col-md-offset-3">
          <div class="panel panel-success">
            <div class="panel-heading">Yay</div>
            <table class="table">
              {{with .YayVotes}}
                {{range .}}
              <tr><td>{{.User.Email | emailat}}</td><tr>
                {{end}}
              {{end}}
            </table>
          </div>
        </div>
        <div class="col-md-3">
          <div class="panel panel-danger">
            <div class="panel-heading">Nay</div>
            <table class="table">
              {{with .NayVotes}}
                {{range .}}
              <tr><td>{{.User.Email | emailat}}</td><tr>
                {{end}}
              {{end}}
            </table>
          </div>
        </div>
      </div>

      <link href="http://alexgorbatchev.com/pub/sh/current/styles/shCore.css" rel="stylesheet" type="text/css" />
      <link href="http://alexgorbatchev.com/pub/sh/current/styles/shThemeDefault.css" rel="stylesheet" type="text/css" />
      <script src="http://alexgorbatchev.com/pub/sh/current/scripts/shCore.js" type="text/javascript"></script>
      <script src="http://alexgorbatchev.com/pub/sh/current/scripts/shBrushDiff.js" type="text/javascript"></script>
      <script>SyntaxHighlighter.all();</script>

{{template "footer.tpl" .}}
{{template "header.tpl" .}}

      {{if .KarmaControls}}<script src="/static/js/post.js"></script>
      <script>
        function postUp() {
          var params = new Array();
          params['type'] = "Up"
          post_to_url(window.location.href, params, "POST")
        }

        function postDown() {
          var params = new Array();
          params['type'] = "Down"
          post_to_url(window.location.href, params, "POST")
        }

        {{if .MaintainerControls}}
        function postMaintainer() {
          var params = new Array();
          params['type'] = "Maintainer"
          post_to_url(window.location.href, params, "POST")
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
              <h1>{{.Package.Name}} <small>OMV-{{.Package.BuildDate.Year}}-{{.Package.Id}} {{.Header}}</small><div class="pull-right">{{if .KarmaControls}}<a href="#" class="btn" onclick="postUp()"><i class="fa fa-3x {{if .KarmaUpYes}}fa-thumbs-up{{else}}fa-thumbs-o-up{{end}}"></i></a>{{end}} {{.Karma}} {{if .KarmaControls}}<a href="#" class="btn" onclick="postDown()"><i class="fa fa-3x {{if .KarmaDownYes}}fa-thumbs-down{{else}}fa-thumbs-o-down{{end}}"></i></a>{{if .MaintainerControls}}<a href="#" class="btn" onclick="postMaintainer()"><i class="fa fa-3x {{if .KarmaMaintainerYes}}fa-check-square{{else}}fa-check-square-o{{end}}"></i></a>{{end}}{{end}}</div></h1>
            </div>
            <table class="table">
              <tbody>
                <tr>
                  <td><b>Source</b></td>
                  <td>{{.Package.Handler}}</td>
                </tr>
                <tr>
                  <td><b>Handler ID</b></td>
                  <td>{{.Package.HandleId}}</td>
                </tr>
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
                  <td><b>Architecture<b></td>
                  <td>{{.Package.Architecture}}</td>
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
                    <table class="table table-bordered">
                      <thead>
                        <tr>
                          <th>Type</th>
                          <th>Name</th>
                          <th>Epoch</th>
                          <th>Version</th>
                          <th>Release</th>
                        </tr>
                      </thead>
                      <tbody>
                        {{with .Package.Packages}}
                          {{range .}}
                        <tr>
                          <td>{{.Type}}</td>
                          <td>{{.Name}}</td>
                          <td>{{.Epoch}}</td>
                          <td>{{.Version}}</td>
                          <td>{{.Release}}</td>
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

      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          <div class="panel panel-primary">
            <div class="panel-heading">Changelog</div>
              <pre class="pre-scrollable">{{if .Changelog}}{{.Changelog}}{{else}}Not Available{{end}}</pre>
          </div>
        </div>
      </div>

{{template "footer.tpl" .}}
{{template "header.tpl" .}}

      <div class="row table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th>Update ID</th>
              <th>Name</th>
              <th>Submitter</th>
              <th>Platform</th>
              <th>Repository</th>
              <th>Architecture</th>
              <th>Status</th>
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
                <td>{{.Status}}</td>
                <td>{{.BuildDate | since}}</td>
              </tr>
              {{end}}
            {{end}}
        </table>
      </div>
      <div class="row">
        <div class="col-md-6 col-md-offset-3">
          <form name="input" method="get">
            <div class="input-group">
              <span class="input-group-btn">
                <a href="?page={{.PrevPage}}"><button class="btn btn-default" type="button">&lt;&lt;</button></a>
              </span>
              <span class="input-group-addon">Page</span>
              <input type="text" name="page" class="form-control" placeholder="{{.Page}}">
              <span class="input-group-addon">/ {{.Pages}}</span>
              <span class="input-group-btn">
                <a href="?page={{.NextPage}}"><button class="btn btn-default" type="button">&gt;&gt;</button></a>
              </span>
            </div>
          </form>
        </div>
      </div>

{{template "footer.tpl" .}}
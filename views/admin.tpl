{{template "header.tpl" .}}

      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          <div class="page-header">
            <center><h1>Administrator Interface <small>Users &amp; Permissions</small></h1><center>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-lg-6">
          <h3>Users</h3>
          <form class="form-inline" role="form" method="post">
            {{ .xsrf_data }}
            <div class="form-group">
              <label class="sr-only" for="email">Email Address</label>
              <input type="email" class="form-control" name="email" id="email" placeholder="{{if .User}}{{.User.Email}}{{else}}Email Address{{end}}">
            </div>
            <button type="submit" class="btn btn-default">Submit</button>
          </form>
          {{if .User}}
          <form class="form-inline" role="form" method="post">
            {{ .xsrf_data }}
            <div class="form-group">
              <label class="sr-only" for="add">Add Permission</label>
              <input type="text" class="form-control" name="add" id="add" placeholder="Add Permission">
              <input type="hidden" name="email" value="{{.User.Email}}">
            </div>
            <button type="submit" class="btn btn-default">Add</button>
          </form>
          <form class="form-inline" role="form" method="post">
            {{ .xsrf_data }}
            <div class="form-group">
              <input type="hidden" name="email" value="{{.User.Email}}">
              <label class="sr-only" for="add">Remove Permission</label>
              <select name="rm" class="form-control">
                {{range .User.Permissions}}
                <option>{{.Permission}}</option>
                {{end}}
              </select>
              <button type="submit" class="btn btn-default">Remove</button>
            </div>
          </form>
          {{end}}
        </div>
        <div class="col-lg-6">
          <h3>Permissions</h3>
          <table class="table table-hover table-condensed">
            {{with .Permissions}}
              {{range $key, $value := .}}
            <tr>
              <td>{{$key}}</td>
              <td><div class="list-group">
                  {{range $value}}
                <a href="{{urldata "/admin?email={{.}}" .}}" class="list-group-item">{{.}}</a>
                  {{end}}
              </div></td>
            </tr>
              {{end}}
            {{end}}
          </table>
       </div>
      </div>

{{template "footer.tpl" .}}
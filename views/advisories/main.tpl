{{template "header.tpl" .}}

      <div class="page-header">
        <h1>Advisories <small>Y'all need to get advised.</small></h1>
      </div>

      {{range $index, $element := .Platforms}}
      <div class="row">
        <h2>{{$index}}<a class="btn btn-link pull-right" href="{{urldata "/advisories/{{.}}" $index}}">More &rarr;</a></h2>
        <div class="table-responsive">
          <table class="table table-condensed table-hover">
            {{range $element}}
            <tr>
              <td><a href="{{urldata "/advisories/{{.Id}}" .}}">{{.Prefix}}-{{.Issued.Year}}-{{.AdvisoryId}}</a></td>
              <td>{{.Summary}}</td>
              <td>{{if eq .Type "bugfix"}}<i class="fa fa-bug"></i>{{end}}{{if eq .Type "security"}}<i class="fa fa-shield"></i>{{end}}{{if eq .Type "enhancement"}}<i class="fa fa-gift"></i>{{end}}{{if eq .Type "recommended"}}<i class="fa fa-star-o"></i>{{end}}{{if eq .Type "newpackage"}}<i class="fa fa-plus-square-o"></i>{{end}}{{.Type}}</td>
            </tr>
            {{end}}
          </table>
        </div>
      </div>
      {{end}}

{{template "footer.tpl" .}}
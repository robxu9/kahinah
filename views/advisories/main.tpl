{{template "header.tpl" .}}

      <div class="page-header">
        <h1>Advisories <small>Ya'll need to get advised.</small></h1>
      </div>

      {{range $index, $element := .Platforms}}
      <div class="row">
        <h2>{{$index}}</h2>
        <div class="table-responsive">
          <table class="table table-condensed table-hover">
            {{range $element}}
            <tr>
              <td><a href="{{urldata "/advisories/{{.Id}}" .}}">{{.Prefix}}-{{.Issued.Year}}-{{.AdvisoryId}}</a></td>
              <td>{{.Summary}}</td>
            </tr>
            {{end}}
          </table>
        </div>
        <ul class="pager">
          <li class="next"><a href="{{urldata "/advisories/{{.}}" $index}}">More &rarr;</a></li>
        </ul>
      </div>
      {{end}}

{{template "footer.tpl" .}}
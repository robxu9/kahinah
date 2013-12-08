{{template "header.tpl" .}}

      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>OpenMandriva Update System</h1>
        <p class="lead">Kahinah, the OpenMandriva Update System, is a versatile tool that hooks into <a href="//abf.io">ABF</a>, allowing package developers to focus on developing more without worrying about breakage.</p>
        <p><a class="btn btn-lg btn-success" href="/testing" role="button">See Packages in Testing</a></p>
      </div>

      <!-- Infos -->
      <div class="row">
        <div class="col-lg-6">
          <h2>This is ALPHA quality.</h2>
          <p class="text-danger">This tool has not undergone extensive usage nor testing. Caution is advised. If any updates are not pushed or go missing, please alert OpenMandriva QA.</p>
          <p><a class="btn btn-primary" href="/about" role="button">Contact &raquo;</a></p>
        </div>
        <div class="col-lg-6">
          <h2>News</h2>
          {{with .News}}
            {{range .}}
          <p>{{.}}</p>
            {{end}}
          {{end}}
       </div>
      </div>

{{template "footer.tpl" .}}
{{template "header.tpl" .}}

      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>You don't have permission</h1>
        <p class="lead">You're missing the <code>{{.Permission}}</code> permission.</p>
        <p><a class="btn btn-lg btn-default" href="#" onclick="history.back();" role="button">Back</a></p>
      </div>

{{template "footer.tpl" .}}
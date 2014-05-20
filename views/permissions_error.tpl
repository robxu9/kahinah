{{template "header.tpl" .}}

      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>Yikes! You can't go further!</h1>
        <p class="lead">You're missing the <code>{{.Permission}}</code> permission.</p>
        <p class="lead">...are you sure you logged in?</p>
        <br/>
        <p><a class="btn btn-lg btn-default" href="#" onclick="history.back();" role="button">Back</a></p>
      </div>

{{template "footer.tpl" .}}
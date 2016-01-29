
      <div class="jumbotron">
        <h1>500.</h1>
        <p class="lead">{{.error}}</p>
        <p class="lead">Internal Server Error.</p>
        <br/>
        {{if .stacktrace}}
        <p class="lead">Stacktrace:</p>
        <pre>
{{.stacktrace}}
        </pre>
        {{end}}
        <p><a class="btn btn-lg btn-default" href="#" onclick="history.back();" role="button">Back</a></p>
      </div>


      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>404.</h1>
        {{if .xkcd_today}}
        <p class="lead">Instead, have today's xkcd comic:</p>
        <p class="lead"><img src="{{.xkcd_today}}" title="{{.xkcd_today_title}}" alt="today's xkcd comic" /></p>
        {{else}}
        <p class="lead">I... I dunno what else to tell ya.</p>
        <p class="lead">...I feel so lost.</p>
        {{end}}
        <br/>
        <p><a class="btn btn-lg btn-default" href="#" onclick="history.back();" role="button">Back</a></p>
      </div>

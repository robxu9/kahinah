<section class="section">
    <div class="container">
        <h1 class="title">404</h1>
        <h2 class="subtitle">Page Not Found.</h2>
        {{if .xkcd_today}}
        <p>Instead, have today's xkcd comic:<br/>
            <img src="{{.xkcd_today}}" title="{{.xkcd_today_title}}" alt="today's xkcd comic" />
        </p>
        {{end}}
        <br/>
        <p><a class="button" href="#" onclick="history.back();" role="button">Back</a></p>
    </div>
</section>

<section class="section">
    <div class="container">
        <h1 class="title">500</h1>
        <h2 class="subtitle">Internal Server Error.</h2>
        <p>The error was: {{.error}}</p>
        {{if .stacktrace}}
        <p>Stacktrace:</p>
        <pre>
{{.stacktrace}}
        </pre>
        {{end}}
        <a href="//github.com/robxu9/kahinah">You should report this to the Kahinah developers.</a>
        <br/>
        <p><a class="button" href="#" onclick="history.back();" role="button">Back</a></p>
    </div>
</section>

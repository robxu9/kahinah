<div id="app">
    <section class="section">
        <div class="container">
            <div class="columns">
                <div class="column">
                    <h1 class="title">Welcome to Kahinah.</h1>
                    <p>Kahinah is a versatile tool that hooks into build systems, taking
                        the output of builds from build systems, running them through
                        various tasks that Kahinah has, and allowing developers to
                        accept/reject the builds given the output of various tasks.
                        <br/><br/>
                        Think of it as a pipeline of quality assurance after a product
                        is produced, allowing for more rigorous tests and providing
                        users with the knowledge to make informed decisions.
                        <br/><br/>
                        Package developers can now focus on developing without worrying
                        about extreme breakage.
                    </p>
                </div>
                <div class="column">
                    <h2 class="subtitle">Recent Activity <small><a href="{{url "/i/activity"}}">more...</a></small></h2>
                    <template v-for="activity in activities">
                        <article class="media">
                            <figure class="media-image">
                                <img src="http://placehold.it/60x60">
                            </figure>
                            <div class="media-content">
                                <div class="content">
                                    <p>
                                        <strong>@<strong v-text="activity.User"></strong></strong>
                                        <small>on update <span v-text="activity.ListId"></span></small>
                                        <small v-text="activity.Time | moment &quot;from&quot;"></small>
                                        <br>
                                        <span v-html="activity.Comment"></span>
                                    </p>
                                </div>
                                <nav class="navbar">
                                    <div class="navbar-left">
                                        <a class="navbar-item" v-bind:href="activity.URL">
                                            <span class="icon is-small"><i class="fa fa-reply"></i></span>
                                        </a>
                                    </div>
                                    <div class="navbar-right">
                                        <p class="navbar-item">
                                            now <strong v-text="activity.Karma"></strong> karma
                                        </p>
                                    </div>
                                </nav>
                            </div>
                        </article>
                    </template>
                </div>
            </div>
        </div>
    </section>

    <section class="section">
        <div class="container">
            <div class="columns">
                <div class="column">
                    <h2 class="subtitle">How does it work?</h1>
                    <p>Your build system will build binaries as usual; afterwards,
                        however, is where it changes.
                        <br/><br/>
                        Kahinah will receive a notification of the build, either via
                        polling or webhooks; and will stop any publication of the
                        binaries (if your build system manages publication). It will
                        then collect metadata and run configurable tests (written in
                        JS) and output that into a build view.
                        <br/><br/>
                        Users can then approve or reject a build from there. For
                        more community testing, collective karma can be enabled,
                        which allows users to +1 or -1 and then authors to
                        accept/reject when thresholds are reached.
                        <br/><br/>
                        When thresholds are reached, Kahinah will act accordingly
                        and accept/reject the builds (which may involve publishing
                        them on the build system, throwing them away, etc.).
                    </p>
                </div>
                <div class="column">
                    <h2 class="subtitle">News <small>{{if not .Time.IsZero}}last updated {{.Time}}{{end}}</small></h2>
                    <div class="content">
                        {{.News}}
                    </div>
                </div>
            </div>
        </div>
    </section>
</div>

<script>
    $.getJSON("{{url "/i/activity/json?limit=3"}}", function(data) {
        new Vue({
            el: '#app',
            data: data
        });
    })
</script>

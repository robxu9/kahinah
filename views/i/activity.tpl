<div id="app">
    <section class="section">
        <div class="container">
            <h1 class="title">Recent Activity</h1>
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
            <br />
            <nav class="navbar">
                <div class="navbar-left">
                    <div class="navbar-item">
                        <p class="subtitle is-5">
                            <strong v-text="totalpages"></strong> page(s)
                        </p>
                    </div>
                </div>

                <div class="navbar-right">
                    <p class="navbar-item"><a class="button is-outlined" v-on:click="loadMore" v-if="page < totalpages">Load more...</a><a class="button is-outlined is-disabled" v-else>No more activity.</a></p>
                </div>
            </nav>
        </div>
    </section>
</div>

<script>
    $.getJSON("{{url "/i/activity/json"}}", function(data) {
        new Vue({
            el: '#app',
            data: data,
            methods: {
                loadMore: function(event) {
                    $.getJSON("{{url "/i/activity/json"}}" + "?page=" + (data.page + 1), function(newData) {
                        data.activities = data.activities.concat(newData.activities);
                        data.page = newData.page;
                        data.totalpages = newData.totalpages;
                    });
                }
            }
        });
    });
</script>

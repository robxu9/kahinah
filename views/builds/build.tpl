<div id="app">
    <section class="section">
        <div class="container">
            <nav class="navbar">
                <div class="navbar-left">
                    <p class="title is-5">
                        <strong v-text="nevr"></strong>: update <span v-text="ID"></span>
                        <small>last updated <span v-text="Updated | moment &quot;from&quot;"></span></small>
                    </p>
                </div>

                <div class="navbar-right">
                    <p class="navbar-item"><span v-text="Status"></span></p>
                    <p class="navbar-item"><a href="#activity" v-bind:class="[&quot;button&quot;, &quot;is-large&quot;, statusModifier]"><strong v-text="TotalKarma"></strong></a></p>
                </div>
            </nav>

            <div class="columns">
                <div class="column"> <!-- buildlist info -->
                    <table class="table is-bordered">
                        <tbody>
                            <tr>
                                <td><strong>Platform</strong></td>
                                <td v-text="Platform"></td>
                            </tr>
                            <tr>
                                <td><strong>Channel</strong></td>
                                <td v-text="Channel"></td>
                            </tr>
                            <tr>
                                <td><strong>Variants</strong></td>
                                <td><ul><li v-for="a in Arch">
                                    <span v-text="a"></span>
                                </li>
                                </ul></td>
                            </tr>
                            <tr>
                                <td><strong>Submitter</strong></td>
                                <td><a v-bind:href="userURL" v-text="Submitter"></a></td>
                            </tr>
                            <tr>
                                <td><strong>Build Type</strong></td>
                                <td><span class="icon"><i v-bind:class="[&quot;fa&quot;, typeClass]"></i></span> <span v-text="Type"></span></td>
                            </tr>
                            <tr>
                                <td><strong>Build Date</strong></td>
                                <td><span v-text="BuildDate | moment &quot;dddd, D MMMM YYYY, HH:mm:ss Z&quot;"></span><br />which was <span v-text="BuildDate | moment &quot;from&quot;"></span></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="column"> <!-- links & artifacts -->
                    <div v-if="Advisory" class="message is-success">
                        <div class="message-header">
                            Advisory
                        </div>
                        <div class="message-body">
                            <strong>This update has been attached to an advisory:</strong>
                            <a v-bind:href="Advisory.Url">Link to Advisory</a>
                        </div>
                    </div>
                    <div v-if="Acceptable" class="message is-warning">
                        <div class="message-header">
                            Acceptable
                        </div>
                        <div class="message-body">
                            <strong>This update is now eligible for publishing:</strong>
                            <a href="#advisory">Fill out Advisory</a>
                        </div>
                    </div>
                    <div v-if="Rejectable" class="message is-danger">
                        <div class="message-header">
                            Rejectable
                        </div>
                        <div class="message-body">
                            <strong>This update is now eligible for rejecting:</strong>
                            <a href="#activity">Reject</a>
                        </div>
                    </div>
                    <div class="message is-primary">
                        <div class="message-header">
                            Links
                        </div>
                        <div class="message-body">
                            <ul>
                                <li v-for="l in Links">
                                    <a v-bind:href="l.Url"><strong v-text="l.Name | linkName"></strong></a>
                                </li>
                            </ul>
                        </div>
                    </div>
                    <div class="message">
                        <div class="message-header">
                            Artifacts
                        </div>
                        <div class="message-body">
                            <ul>
                                <li v-for="a in Artifacts">
                                    <a v-bind:href="a.Url"><strong v-text="a.Name"></strong>-<strong v-text="a.Epoch"></strong>:<strong v-text="a.Version"></strong>-<strong v-text="a.Release"></strong>.<strong v-text="a.Arch"></strong> (<span v-text="a.Type"></span>)</a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>

            <div id="activity"> <!-- recent activity -->
                <h2 class="subtitle">Activity</h2>
                <template v-for="a in Activity">
                    <article class="media">
                        <figure class="media-image">
                            <img src="http://placehold.it/60x60">
                        </figure>
                        <div class="media-content">
                            <div class="content">
                                <p>
                                    <strong>@<strong v-text="a.User"></strong></strong>
                                    <small v-if="a.User == User">this is you</small>
                                    <small v-text="a.Time | moment &quot;from&quot;"></small>
                                    <strong v-text="a.Karma | karmaParse"></strong>
                                    <br>
                                    <span v-html="a.Comment"></span>
                                </p>
                            </div>
                        </div>
                    </article>
                </template>
                <div class="notification" v-if="!User">
                    You must login to add karma or comments.
                </div>
                <article class="media">
                    <div class="media-content">
                        <form class="form-inline" role="form" method="post" id="actform">
                            {{ .xsrf_data }}
                            <div class="content">
                                <p>
                                    <strong>add some karma, @<strong v-text="User"></strong>?</strong>
                                    <small>markdown is supported</small>
                                    <br>
                                    <textarea class="textarea" name="comment" form="actform" style="width:100%"></textarea>
                                </p>
                            </div>
                            <nav class="navbar">
                                <div class="navbar-left">
                                    <p class="control is-grouped">
                                        <label class="radio">
                                            <input type="radio" name="type" value="Neutral">
                                            <i class="fa fa-comment"></i> Comment Only
                                        </label>
                                        <label class="radio">
                                            <input type="radio" name="type" value="Up">
                                            <i class="fa fa-arrow-up"></i> Upvote
                                        </label>
                                        <label class="radio">
                                            <input type="radio" name="type" value="Down">
                                            <i class="fa fa-arrow-down"></i> Downvote
                                        </label>
                                        <label v-if="Maintainer" class="radio">
                                            <input type="radio" name="type" value="Maintainer">
                                            <i class="fa fa-upload"></i> Maintainer Vote
                                        </label>
                                        <label v-if="UserIsQA" class="radio">
                                            <input type="radio" name="type" value="QAPush">
                                            <i class="fa fa-plus"></i> QA Push
                                        </label>
                                        <label v-if="UserIsQA" class="radio">
                                            <input type="radio" name="type" value="QABlock">
                                            <i class="fa fa-minus"></i> QA Block
                                        </label>
                                        <label v-if="Rejectable" class="radio">
                                            <input type="radio" name="type" value="Reject">
                                            <i class="fa fa-hand-paper"></i> Reject
                                        </label>
                                    </p>
                                </div>
                                <div class="navbar-right">
                                    <p class="navbar-item">
                                        <button class="button is-primary" type="submit">Submit</button>
                                    </p>
                                </div>
                            </nav>
                        </form>
                    </div>
                </article>
            </div>

            <div>
                <h2 class="subtitle">Stages</h2>
                <!-- show each stage with a pass/fail, waiting for manual approve, reject, etc -->
            </div>

            <div class="columns">
                <div class="column" v-if="!Advisory"> <!-- advisory left -->
                    <h2 class="subtitle">Advisory</h2>
                    <form class="form-inline" role="form" method="post" id="advisoryForm">
                        <div class="content">
                            <p>
                                This update has not been attached to an advisory yet. You should fill in
                                details about this update to the advisory (and link to an existing one,
                                if necessary), then save it. You may only accept this stage if an advisory is
                                published.
                            </p>
                        </div>
                    </form>
                </div>
                <div class="column"> <!-- diff right -->
                </div>
            </div>

        </div>
    </section>
</div>

<script>
    Vue.filter('linkName', function (value) {
        return value.replace("_mainURL", "Main URL")
            .replace("_changelogURL", "Changelog")
            .replace("_scmlogURL", "SCM Commits");
    });
    Vue.filter('karmaParse', function (value) {
        switch (value) {
        case "+":
            return "+1";
        case "-":
            return "-1";
        case "*":
            return "+" + this.MaintainerKarma;
        case "v":
            return "-" + this.BlockKarma;
        case "^":
            return "+" + this.PushKarma;
        case "_":
            return "\xB10";
        }
    });
    $.getJSON("{{urldata "/b/{{.ID}}/json" .}}", function(data, status, jqXHR) {
        var vm = new Vue({
            el: '#app',
            data: data,
            computed: {
                statusModifier: function() {
                    if (data.Status == "testing") {
                        return "is-warning";
                    } else if (data.Status == "rejected") {
                        return "is-danger";
                    } else if (data.Status == "published") {
                        return "is-success";
                    } else {
                        return "is-info";
                    }
                },
                evr: function() {
                    var evr = "unknown";
                    data.Artifacts.forEach(function(a) {
                        if (a.Type == "source") {
                            evr = a.Epoch + ":" + a.Version + "-" + a.Release;
                        }
                    });
                    return evr;
                },
                nevr: function() {
                    return data.Name + " " + this.evr;
                },
                userURL: function() {
                    return "{{url "/u/"}}" + data.Submitter;
                },
                typeClass: function() {
                    if (data.Type == "bugfix") {
                        return "fa-bug";
                    } else if (data.Type == "security") {
                        return "fa-shield";
                    } else if (data.Type == "enhancement") {
                        return "fa-gift";
                    } else if (data.Type == "recommended") {
                        return "fa-star";
                    } else if (data.Type == "newpackage") {
                        return "fa-plus-square";
                    }
                    return "";
                }
            }
        });
    });
</script>

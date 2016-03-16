<div id="app">
    <section class="section">
        <div class="container">
            <!-- filtering -->
            <nav class="navbar">
                <!-- filters... -->
                <div class="navbar-left">
                    <div class="navbar-item">
                        <p class="control is-grouped">
                            <button class="button is-primary">Platform</button>
                            <input v-model="filters.platform" debounce="500" class="input" type="text" placeholder="(all)">
                        </p>
                    </div>
                    <div class="navbar-item">
                        <p class="control is-grouped">
                            <button class="button is-info">Channel</button>
                            <input v-model="filters.channel" debounce="500" class="input" type="text" placeholder="(all)">
                        </p>
                    </div>
                    <div class="navbar-item">
                        <p class="control is-grouped">
                            <button class="button is-success">Status</button>
                            <span class="select">
                                <select v-model="filters.status" debounce="500">
                                    <option selected></option>
                                    <option>pending</option>
                                    <option>running</option>
                                    <option>success</option>
                                    <option>failed</option>
                                </select>
                            </span>
                        </p>
                    </div>
                    <div class="navbar-item">
                        <p class="control is-grouped">
                            <button class="button is-warning">Limit</button>
                            <input v-model="filters.limit" debounce="500" class="input" type="number" min="1" max="500" placeholder="(50)">
                        </p>
                    </div>
                </div>

                <!-- pagination -->
                <div class="navbar-right">
                    <div class="navbar-item">
                        <p class="control is-grouped">
                            <button class="button" v-on:click="filters.page -= 1"><span class="icon"><i class="fa fa-arrow-left"></i></span></button>
                            <input v-model="filters.page" debounce="500" class="input" type="number" min="1" v-bind:max="result.pages.total">
                            <button class="button" v-on:click="filters.page += 1"><span class="icon"><i class="fa fa-arrow-right"></i></span></button>
                        </p>
                    </div>
                    <p class="navbar-item"><strong v-text="result.pages.total"></strong> pages</p>
                </div>
            </nav>
            <!-- results -->
            <template v-for="l in result.lists">
                <hr>
                <div class="columns">
                    <div class="column is-2"> <!-- ID -->
                        <a v-bind:href="baseURL + &quot;b/&quot; + l.ID"><strong v-text="l.ID"></strong></a>
                    </div>
                    <div class="column is-2"> <!-- name -->
                        <a v-bind:href="baseURL + &quot;b/&quot; + l.ID"><span v-text="l.Name"></span></a>
                    </div>
                    <div class="column is-4"> <!-- target -->
                        <span v-text="l.Platform"></span>: <span v-text="l.Channel"></span> (<span v-text="l.Variants | semicolontocomma"></span>)
                    </div>
                    <div class="column is-1"> <!-- status -->
                        <!-- running, pending, success, failed -->
                        <span v-if="l.StageResult == &quot;running&quot;" class="tag">Running</span>
                        <span v-if="l.StageResult == &quot;pending&quot;" class="tag is-warning">Pending</span>
                        <span v-if="l.StageResult == &quot;success&quot;" class="tag is-success">Success</span>
                        <span v-if="l.StageResult == &quot;failed&quot;" class="tag is-danger">Failed</span>
                    </div>
                    <div class="column is-1" v-if="l.AdvisoryID"> <!-- advisory -->
                        <span v-bind:href="baseURL + &quot;a/&quot; + l.AdvisoryID" class="tag is-info">Advisory</span>
                    </div>
                    <div class="column"> <!-- updated at -->
                        <span v-text="l.UpdatedAt | moment &quot;from&quot;"></span>
                    </div>
                </div>
            </template>
        </div>
    </section>
</div>

<script>
    Vue.filter('semicolontocomma', function (value) {
        return value.split(";").join(", ");
    });

    var params = $.param({
        platform: "{{.Platform}}",
        channel: "{{.Channel}}",
        status: "{{.Status}}",
        limit: "{{.Limit}}",
        page: "{{.Page}}"
    });

    $.getJSON("{{urldata "/i/list/json?" .}}" + params, function(data, status, jqXHR) {
        // set the current page
        params.page = data.pages.current;

        // setup vue.js
        var vm = new Vue({
            el: '#app',
            data: {
                filters: {
                    platform: params.platform,
                    channel: params.channel,
                    status: params.status,
                    limit: params.limit,
                    page: params.page
                },
                result: data
            },
            watch: {
                'filters': {
                    handler: function(newData, oldData) {
                        // retrieve new data based on the filters
                        var outer = this;
                        $.getJSON("{{urldata "/i/list/json?" .}}" + $.param(this.filters), function(newData, status, jqXHR) {
                            outer.filters.page = newData.pages.current;
                            outer.result = newData;
                        });
                    },
                    deep: true
                }
            },
            computed: {
                baseURL: function() {
                    return "{{url "/"}}";
                }
            }
        });
    });
</script>

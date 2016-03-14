<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="">
    <meta name="author" content="">
    <meta name="_xsrf" content="{{.xsrf_token}}">
    <link rel="shortcut icon" href="{{url "/static/img/favicon.png"}}">
    <title>{{if .Title}}{{.Title}} | {{end}}Kahinah</title>

    <!-- bulma -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.0.16/css/bulma.min.css" rel="stylesheet">

    <!-- Font Awesome -->
    <link href="//maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css" rel="stylesheet" integrity="sha256-3dkvEK0WLHRJ7/Csr0BZjAWxERc5WH7bdeUya2aXxdU= sha512-+L4yy6FRcDGbXJ9mPG8MT/3UCDzwR9gPeyFNMCtInsol++5m3bk2bXWKdZjvybmohrAsn3Ua5x8gfLnbE1YkOg=="
        crossorigin="anonymous">

    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
        <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->

    <!-- Offline support -->
    <link href="{{url "/static/css/offline-theme-chrome.css"}}" rel="stylesheet">
    <link href="{{url "/static/css/offline-language-english.css"}}" rel="stylesheet">

    <!-- Javascript Files -->
    <script>
        window.urlPrefix = "{{url ""}}";
    </script>
    <script src="//code.jquery.com/jquery-2.2.0.min.js"></script>
    <script src="//code.jquery.com/ui/1.11.4/jquery-ui.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/moment.js/2.11.2/moment-with-locales.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/moment-timezone/0.5.0/moment-timezone-with-data.min.js"></script>
    <!-- Vue.js -->
    {{if eq .environment "dev"}}
    <script src="https://cdnjs.cloudflare.com/ajax/libs/vue/1.0.16/vue.js"></script>
    <script src="{{url "/static/js/vue-moment.js"}}"></script>
    <script>Vue.config.debug = true</script>
    {{else}}
    <script src="https://cdnjs.cloudflare.com/ajax/libs/vue/1.0.16/vue.min.js"></script>
    <script src="{{url "/static/js/vue-moment.min.js"}}"></script>
    {{end}}
    <script src="{{url "/static/js/xsrf.js"}}"></script>
    <!-- Offline support -->
    <script src="{{url "/static/js/offline.min.js"}}"></script>
</head>
<body>
    <header class="header">
        <div class="container">
            <div class="header-left">
                <div class="header-item">
                    <span class="icon"><i class="fa fa-tasks"></i></span> Kahinah
                </div>
                <a class="header-tab {{if eq .Nav 0}}is-active{{end}}" href="{{url "/"}}">Dashboard</a>
                <a class="header-tab {{if eq .Nav 1}}is-active{{end}}" href="{{url "/i/activity"}}">Recent Activity</a>
                <a class="header-tab {{if eq .Nav 2}}is-active{{end}}" href="{{url "/i/list"}}">Lists</a>
                <a class="header-tab {{if eq .Nav 3}}is-active{{end}}" href="{{url "/i/advisories"}}">Advisories</a>
            </div>

            <!-- Hamburger Menu -->
            <span class="header-toggle">
                <span></span>
                <span></span>
                <span></span>
            </span>

            <div class="header-right header-menu">
                <span class="header-item">
                    <a href="{{url "/u"}}">{{.authenticated}}</a>
                </span>
                <span class="header-item">
                    {{if .authenticated}}
                    <a class="button" href="{{url "/u/logout"}}">Logout</a>
                    {{else}}
                    <a class="button" href="{{url "/u/login"}}">Login</a>
                    {{end}}
                </span>
            </div>
        </div>
    </header>

    {{if .flash_err}}<div class="notification is-danger">{{.flash_err}}</div>{{end}}
    {{if .flash_warn}}<div class="notification is-warning">{{.flash_err}}</div>{{end}}
    {{if .flash_info}}<div class="notification is-info">{{.flash_err}}</div>{{end}}

    {{yield}}

    <section class="hero is-info is-right">
        <div class="hero-content">
            <div class="container">
                <h1 class="title">
                    Kahinah, v4
                </h1>
                <h2 class="subtitle">
                    Copyright &copy; 2013-{{.copyright}} Robert Xu.
                    Licensed under the MIT license. <a href="https://github.com/robxu9/kahinah" style="color: #fff;"><span class="icon"><i class="fa fa-github"></i></span></a>
                </h2>
            </div>
        </div>
    </section>
</body>

</html>

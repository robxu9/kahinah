<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="">
    <meta name="author" content="">

    <meta name="_xsrf" content="{{.xsrf_token}}" />

    <link rel="shortcut icon" href="{{url "/static/img/favicon.png"}}">

    <title>{{.Title}} | Kahinah</title>

    <!-- Bootstrap core CSS -->
    <link href="https://maxcdn.bootstrapcdn.com/bootswatch/3.3.6/paper/bootstrap.min.css" rel="stylesheet" integrity="sha256-ZSKfhECi0yCEmGVAuQLWTHtJEn2vBNPexEWsJCIC/Nc= sha512-b+mSnD4QXw1uYwkgJ3d0XxoMXo+ZKHJNTNNFIddJ0IazcwKvKYtIlWADZ1JEREJzxUG428sfCw7qDuswAFcrOQ==" crossorigin="anonymous">

    <!-- Font Awesome -->
    <link href="https://maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css" rel="stylesheet" integrity="sha256-3dkvEK0WLHRJ7/Csr0BZjAWxERc5WH7bdeUya2aXxdU= sha512-+L4yy6FRcDGbXJ9mPG8MT/3UCDzwR9gPeyFNMCtInsol++5m3bk2bXWKdZjvybmohrAsn3Ua5x8gfLnbE1YkOg==" crossorigin="anonymous">

    <link href="{{url "/static/css/justified-nav.css"}}" rel="stylesheet">

    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="//oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
      <script src="//oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js"></script>
    <![endif]-->

    <!-- Bootstrap core JavaScript -->
    <script>window.urlPrefix = "{{url ""}}";</script>
    <script src="//code.jquery.com/jquery-2.2.0.min.js"></script>
    <script src="{{url "/static/js/xsrf.js"}}"></script>
  </head>

  <body>

    <div class="container">

      <nav class="navbar navbar-default" role="navigation">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navbar-collapse">
            <span class="sr-only">Toggle Navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a href="{{url "/"}}" class="navbar-brand"><i class="fa fa-tasks"></i> Kahinah</a>
        </div>

        <div class="collapse navbar-collapse" id="navbar-collapse">
          <ul class="nav navbar-nav">
            <li{{if eq .Loc 0}} class="active" {{end}}><a href="{{url "/"}}">Home</a></li>

            <li {{if eq .Loc 1}}class="active" {{end}}class="dropdown"> <!-- builds -->
              <a href="#" class="dropdown-toggle" data-toggle="dropdown">Builds <b class="caret"></b></a>
              <ul class="dropdown-menu">
                <li{{if eq .Loc 1}}{{if eq .Tab 1}} class="active"{{end}}{{end}}><a href="{{url "/builds/testing"}}">Testing</a></li>
                <li{{if eq .Loc 1}}{{if eq .Tab 2}} class="active"{{end}}{{end}}><a href="{{url "/builds/published"}}">Published</a></li>
                <li{{if eq .Loc 1}}{{if eq .Tab 3}} class="active"{{end}}{{end}}><a href="{{url "/builds/rejected"}}">Rejected</a></li>
                <li{{if eq .Loc 1}}{{if eq .Tab 4}} class="active"{{end}}{{end}}><a href="{{url "/builds"}}">All</a></li>
              </ul>
            </li>

            <li {{if eq .Loc 2}}class="active" {{end}}class="dropdown"> <!-- advisories -->
              <a href="#" class="dropdown-toggle" data-toggle="dropdown">Advisories <b class="caret"></b></a>
              <ul class="dropdown-menu">
                <li{{if eq .Loc 2}}{{if eq .Tab 1}} class="active"{{end}}{{end}}><a href="{{url "/advisories"}}">Recent</a></li>
                <li{{if eq .Loc 2}}{{if eq .Tab 2}} class="active"{{end}}{{end}}><a href="{{url "/advisories/all"}}">All</a></li>
                <li class="divider"></li>
                <li{{if eq .Loc 2}}{{if eq .Tab -1}} class="active"{{end}}{{end}}><a href="{{url "/advisories/new"}}">New Advisory</a></li>
              </ul>
            </li>

            <li {{if eq .Loc 3}}class="active" {{end}}class="dropdown"> <!-- virtual testing -->
              <a href="#" class="dropdown-toggle" data-toggle="dropdown">Virtual Testing <b class="caret"></b></a>
              <ul class="dropdown-menu">
                <li{{if eq .Loc 3}}{{if eq .Tab 1}} class="active"{{end}}{{end}}><a href="{{url "/vtests/running"}}">Currently Running</a></li>
                <li{{if eq .Loc 3}}{{if eq .Tab 2}} class="active"{{end}}{{end}}><a href="{{url "/vtests/recent"}}">Recent Tests</a></li>
                <li{{if eq .Loc 3}}{{if eq .Tab 3}} class="active"{{end}}{{end}}><a href="{{url "/vtests/platform"}}">By Platform</a></li>
              </ul>
            </li>

            <li {{if eq .Loc 4}}class="active" {{end}}class="dropdown"> <!-- appstream -->
              <a href="#" class="dropdown-toggle" data-toggle="dropdown">AppStream Check <b class="caret"></b></a>
              <ul class="dropdown-menu">
                <li{{if eq .Loc 4}}{{if eq .Tab 1}} class="active"{{end}}{{end}}><a href="{{url "/appstream/desktop"}}">Desktop Applications</a></li>
                <li{{if eq .Loc 4}}{{if eq .Tab 2}} class="active"{{end}}{{end}}><a href="{{url "/appstream/console"}}">Console Applications</a></li>
                <li{{if eq .Loc 4}}{{if eq .Tab 3}} class="active"{{end}}{{end}}><a href="{{url "/appstream/unclassified"}}">Unclassified</a></li>
                <li class="divider"></li>
                <li{{if eq .Loc 4}}{{if eq .Tab 4}} class="active"{{end}}{{end}}><a href="{{url "/appstream/api"}}">API</a></li>
              </ul>
            </li>


            <li{{if eq .Loc -1}}class="active"{{end}}><a href="{{url "/about"}}">About</a></li>

          </ul>

          <!-- login -->
          <div class="navbar-right">
            {{if .authenticated}}
            <p class="navbar-text">{{.authenticated}}</p>
            <a href="{{url "/user/logout"}}"><button class="btn btn-warning navbar-btn">Logout</button></a>
            {{else}}
            <a href="{{url "/user/login"}}"><button class="btn btn-primary navbar-btn">Login</button></a>
            {{end}}
          </div>
        </div>
      </nav>

      {{if .flash_err}}<div class="alert alert-danger">{{.flash_err}}</div>{{end}}
      {{if .flash_warn}}<div class="alert alert-warning">{{.flash_warn}}</div>{{end}}
      {{if .flash_info}}<div class="alert alert-success">{{.flash_info}}</div>{{end}}

      {{yield}}

      <!-- Site footer -->
      <div class="footer">
        Copyright &copy; 2013-{{.copyright}} Robert Xu.<div class="pull-right">Licensed under the MIT license - <a href="//github.com/robxu9/kahinah"><i class="fa fa-github"></i> Github</a></div>
      </div>

    </div> <!-- /container -->

    <script src="//code.jquery.com/ui/1.11.4/jquery-ui.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js" integrity="sha256-KXn5puMvxCw+dAYznun+drMdG1IFl3agK0p/pqT9KAo= sha512-2e8qq0ETcfWRI4HJBzQiA3UoyFk6tbNyG+qSaIBZLyW9Xf3sWZHN/lxe9fTh1U45DpPf07yj94KsUHHWe4Yk1A==" crossorigin="anonymous"></script>

  </body>
</html>

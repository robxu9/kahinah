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
    <link href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet">

    <!-- Font Awesome -->
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.min.css" rel="stylesheet">

    <link href="{{url "/static/css/justified-nav.css"}}" rel="stylesheet">

    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="//oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
      <script src="//oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js"></script>
    <![endif]-->

    <!-- Bootstrap core JavaScript -->
    <script>window.urlPrefix = "{{url ""}}";</script>
    <script src="//code.jquery.com/jquery-2.0.3.min.js"></script>
    <script src="{{url "/static/js/xsrf.js"}}"></script>

    <script src="//login.persona.org/include.js"></script>
    <script src="{{url "/static/js/persona.js"}}"></script>

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
          <a href="{{url "/"}}" class="navbar-brand">Kahinah</a>
        </div>

        <div class="collapse navbar-collapse" id="navbar-collapse">
          <ul class="nav navbar-nav">
            <li{{if eq .Tab 0}} class="active"{{end}}><a href="{{url "/"}}">Home</a></li>
            <li{{if eq .Tab 1}} class="active"{{end}}><a href="{{url "/testing"}}">Testing</a></li>
            <li{{if eq .Tab 2}} class="active"{{end}}><a href="{{url "/published"}}">Published</a></li>
            <li{{if eq .Tab 3}} class="active"{{end}}><a href="{{url "/rejected"}}">Rejected</a></li>
            <li{{if eq .Tab 4}} class="active"{{end}}><a href="{{url "/builds"}}">All</a></li>
            <li{{if eq .Tab 5}} class="active"{{end}}><a href="{{url "/about"}}">About</a></li>
          </ul>

          <!-- login -->
          <div class="navbar-right">
            <p class="navbar-text" id="persona-user"></p>
            <button type="button" class="btn btn-primary navbar-btn" style="display: none" id="login">Persona Login</button><button type="button" class="btn btn-warning navbar-btn" style="display: none" id="logout">Logout</button>
          </div>
        </div>
      </nav>

      {{if .flash.error}}<div class="alert alert-danger">{{.flash.error}}</div>{{end}}
      {{if .flash.warning}}<div class="alert alert-warning">{{.flash.warning}}</div>{{end}}
      {{if .flash.notice}}<div class="alert alert-success">{{.flash.notice}}</div>{{end}}

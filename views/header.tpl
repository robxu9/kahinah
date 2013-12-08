<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="shortcut icon" href="../../docs-assets/ico/favicon.png">

    <title>{{.Title}} | Kahinah</title>

    <!-- Bootstrap core CSS -->
    <link href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet">

    <!-- Font Awesome -->
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.min.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="/static/css/justified-nav.css" rel="stylesheet">

    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="//oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
      <script src="//oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js"></script>
    <![endif]-->

    <!-- Bootstrap core JavaScript -->
    <script src="//code.jquery.com/jquery-2.0.3.min.js"></script>
    <script src="//login.persona.org/include.js"></script>
    <script src="/static/js/persona.js"></script>

  </head>

  <body>

    <div class="container">

      <div class="masthead">
        <h3 class="text-muted">Kahinah<button type="button" class="btn btn-success pull-right" style="display: none" id="login">Persona Login</button><button type="button" class="btn btn-warning pull-right" style="display: none" id="logout">Logout</button></h3>
        <ul class="nav nav-justified">
          <li{{if eq .Tab 0}} class="active"{{end}}><a href="/">Home</a></li>
          <li{{if eq .Tab 1}} class="active"{{end}}><a href="/testing">Testing Updates</a></li>
          <li{{if eq .Tab 2}} class="active"{{end}}><a href="/published">Published Updates</a></li>
          <li{{if eq .Tab 3}} class="active"{{end}}><a href="/rejected">Rejected Updates</a></li>
          <li{{if eq .Tab 4}} class="active"{{end}}><a href="/builds">All Updates</a></li>
          <li{{if eq .Tab 5}} class="active"{{end}}><a href="/about">About</a></li>
        </ul>
      </div>

      {{if .flash.error}}<div class="alert alert-danger">{{.flash.error}}</div>{{end}}
      {{if .flash.warning}}<div class="alert alert-warning">{{.flash.warning}}</div>{{end}}
      {{if .flash.notice}}<div class="alert alert-success">{{.flash.notice}}</div>{{end}}

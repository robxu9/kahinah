$("document").ready(function(){

  $("#logout").hide();

  function loggedIn(email){
    $("#login").hide();
    $("#logout").show();
    $("#persona-user").text(email);
  }

  function loggedOut(){
    $("#logout").hide();
    $("#login").show();
    $("#persona-user").text("");
  }

  var user = null;

  $.ajax({
    url: window.urlPrefix + '/auth/check',
    success: function(res, status, xhr) {
      if (res === "") {
        loggedOut();
      }
      else {
        loggedIn(res);
        user = res;
      }
    },
    async: false
  });

  navigator.id.watch({
    loggedInUser: user,
    onlogin: function(assertion) {
      $("#login").text("Logging in...");
      $.ajax({
        type: 'POST',
        url: window.urlPrefix + '/auth/login',
        data: {assertion: assertion},
        success: function(res, status, xhr) {
          location.reload(true);
        },
        error: function(xhr, status, err) { 
          $("#login").text("Failed!");
          $("#login").attr("class","btn btn-danger navbar-btn");
          navigator.id.logout();
        },
      });
    },
    onlogout: function() {
      $.get(window.urlPrefix + '/auth/logout', function() {
        location.reload(true);
      });
    }
  });

  $("#login").on("click", function(e) {
    e.preventDefault();
    $("#login").text("[Popup Appeared]");
    $("#login").attr("class", "btn btn-info navbar-btn");
    navigator.id.request({
      siteName: "Kahinah",
      oncancel: function() {
        $("#login").text("Login Canceled");
        $("#login").attr("class","btn btn-danger navbar-btn");
      },
    });
  });

  $("#logout").on("click", function(e) {
    e.preventDefault();
    navigator.id.logout();
  });

});

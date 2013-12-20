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

  $("#login").on("click", function(e) {
    e.preventDefault();
    navigator.id.get(mailVerified, {siteName: "Kahinah"});
  });

  $("#logout").on("click", function(e) {
    e.preventDefault();
    $.get('/auth/logout', onlogout);
  });

  function onlogout() {
    location.reload(true);
  }

  function mailVerified(assertion){
    $("#login").text("Logging in...");
    $.ajax({
      type: 'POST',
      url: '/auth/login',
      data: {assertion: assertion},
      success: function(res, status, xhr) {
        location.reload(true);
      },
      error: function(xhr, status, err) { 
        $("#login").text("Failed to login! Try again?");
      },
    });
  }

  $.get('/auth/check', function (res) {
    if (res === "") loggedOut();
    else loggedIn(res);
  });
});

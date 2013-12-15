$("document").ready(function(){

  $("#logout").hide();

  function loggedIn(email){
    $("#login").hide();
    $("#logout").show();
    $("#logout").text("Logout [" + email + "]");
  }

  function loggedOut(){
    $("#logout").hide();
    $("#login").show();
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
    $.ajax({
      type: 'POST',
      url: '/auth/login',
      data: {assertion: assertion},
      success: function(res, status, xhr) {
        location.reload(true);
      },
      //error: function(xhr, status, err) { alert("Didn't login: " + err); }
    });
  }

  $.get('/auth/check', function (res) {
    if (res === "") loggedOut();
    else loggedIn(res);
  });
});

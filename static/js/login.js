$('#login-btn').on('click', function(e) {
  $.post('/api/login', {
    username: $('#username').val(),
    password: $('#password').val()
  })
  .done(function() {
    window.location.href = '/order';
  })
  .fail(function(xhr) {
    if (xhr.status == 400) {
      var msg = "Invalid username / password";
    } else {
      var msg = "Internal server error";
    }
    $('#status').text(msg).show();
  });
});

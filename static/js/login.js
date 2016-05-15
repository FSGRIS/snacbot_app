window.onload = function() {
  $('#login-btn').on('click', function(e) {
    var body = {
      email: $('#email').val(),
      password: $('#password').val()
    };
    $.post('/api/login', JSON.stringify(body), 'json')
      .done(function() {
        window.location.href = '/order';
      })
      .fail(function(xhr) {
        $('#status').text(xhr.responseText).show();
      });
  });
}

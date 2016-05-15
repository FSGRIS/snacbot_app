window.onload = function() {
  $('#login-btn').on('click', function(e) {
    var data = {
      email: $('#email').val(),
      password: $('#password').val()
    };
    $.post('/api/login', JSON.stringify(data), 'json')
      .done(function() {
        window.location.href = '/order';
      })
      .fail(function(xhr) {
        $('#status').text(xhr.responseText).show();
      });
  });
}

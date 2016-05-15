window.onload = function() {
  $('#create-btn').on('click', function(e) {
    if ($('#password').val() !== $('#confirmPassword').val()) {
      $('#status').text("Passwords do not match.").show();
    }
    var body = {
      email: $('#email').val(),
      password: $('#password').val(),
      orgName: $('#orgName').val(),
      orgCode: $('#orgCode').val()
    };
    $.post('/api/create_account', JSON.stringify(body), 'json')
      .done(function() {
        window.location.href = '/order';
      })
      .fail(function(xhr) {
        $('#status').text(xhr.responseText).show();
      });
  });
}

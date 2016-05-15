window.onload = function() {
  $('#create-btn').on('click', function(e) {
    if ($('#password').val() !== $('#confirmPassword').val()) {
      $('#status').text("Passwords do not match.").show();
    }
    var data = {
      email: $('#email').val(),
      password: $('#password').val(),
      orgName: $('#orgName').val(),
      orgCode: $('#orgCode').val()
    };
    $.post('/api/create_account', JSON.stringify(data), 'json')
      .done(function() {
        window.location.href = '/order';
      })
      .fail(function(xhr) {
        $('#status').text(xhr.responseText).show();
      });
  });
}

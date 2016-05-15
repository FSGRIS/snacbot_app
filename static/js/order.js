window.onload = function() {
  var deliveryLocations = [
    { id: 1, x: 390, y: 250, selected: false },
    { id: 2, x: 250, y: 420, selected: true },
    { id: 3, x: 150, y: 140, selected: false }
  ];
  var selectedLocation;
  for (var i = 0; i < deliveryLocations.length; i++) {
    var l = deliveryLocations[i];
    if (l.selected) {
      selectedLocation = l;
      break;
    }
  }

  var w = 500;
  var h = 500;

  var svg = d3.select('#map')
    .append('svg')
    .attr('width', w)
    .attr('height', h)

  svg.append('defs')
    .append('pattern')
    .attr('id', 'bg')
    .attr('patternUnits', 'userSpaceOnUse')
    .attr('width', w)
    .attr('height', h)
    .append('image')
    .attr('xlink:href', '/static/img/map.jpg')
    .attr('width', w)
    .attr('height', h);

  svg.append('rect')
    .attr('x', 0)
    .attr('y', 0)
    .attr('width', w)
    .attr('height', h)
    .attr('fill', 'url(#bg)');

  var circles = svg.selectAll('circle')
    .data(deliveryLocations)
    .enter()
    .append('circle');

  function update() {
    circles.attr('fill', function(d) { return d.selected ? 'red' : 'white' });
  }

  circles
    .attr('cx', function(d) { return d.x; })
    .attr('cy', function(d) { return d.y; })
    .attr('r', function(d) { return 10; })
    .attr('fill', function(d) { return d.selected ? 'red' : 'white' })
    .attr('stroke', 'black')
    .attr('stroke-width', 1)
    .on('click', function(d, i) {
      d.selected = true;
      selectedLocation = d;
      for (var j = 0; j < deliveryLocations.length; j++) {
        if (j != i) {
          deliveryLocations[j].selected = false;
        }
      }
      update();
    });

  function showStatus(msg, success) {
    $('#status')
      .toggleClass(false)
      .toggleClass('alert ' + (success ? 'alert-success' : 'alert-danger'))
      .text(msg)
      .show();
  }

  $('#order-btn').on('click', function() {
    var snacks = [];
    $('.quantity').each(function() {
      var id = parseInt($(this).attr('data-snack-id'));
      var quantity = parseInt($(this).find('select option:selected').text());
      if (quantity > 0) {
        snacks.push({id: id, quantity: quantity})
      }
    });
    if (snacks.length === 0) {
      showStatus('No snacks selected', false);
      return;
    }
    var body = {
      locationID: selectedLocation.id,
      saveLocation: $('#save-location').is(':checked'),
      snacks: snacks,
    };
    $.post('/api/order', JSON.stringify(body), 'json')
      .done(function() {
        showStatus('Order placed!', true);
      })
      .fail(function(xhr) {
        showStatus('Order failed: ' + xhr.responseText);
      });
  });
}

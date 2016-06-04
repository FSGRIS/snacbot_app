window.onload = function() {
  $.get('/api/locations')
    .done(populateMap)
    .fail(function(xhr) {
      alert("Could not get locations: " + xhr.responseText);
    });
};

function populateMap(locById) {
  var w = 222 * 2;
  var h = 207 * 2;

  var locations = [];
  Object.keys(locById).forEach(function(lid) {
    locations.push({
      id: parseInt(lid),
      x: (locById[lid].x * 40),
      y: h - (locById[lid].y * 40),
      selected: false
    });
  });
  var selectedLocation = null;
  for (var i = 0; i < locations.length; i++) {
    var l = locations[i];
    if (l.selected) {
      selectedLocation = l;
      break;
    }
  }

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
    .attr('xlink:href', '/static/img/newmap.jpg')
    .attr('width', w)
    .attr('height', h);

  svg.append('rect')
    .attr('x', 0)
    .attr('y', 0)
    .attr('width', w)
    .attr('height', h)
    .attr('fill', 'url(#bg)');

  var circles = svg.selectAll('circle')
    .data(locations)
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
      for (var j = 0; j < locations.length; j++) {
        if (j != i) {
          locations[j].selected = false;
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
    if (selectedLocation === null) {
      showStatus('No location selected', false);
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

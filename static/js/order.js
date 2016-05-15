window.onload = function() {
  var deliveryPoints = [
    { id: 1, x: 390, y: 250, selected: false },
    { id: 2, x: 250, y: 420, selected: true },
    { id: 3, x: 150, y: 140, selected: false }
  ];

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
    .data(deliveryPoints)
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
      for (var j = 0; j < deliveryPoints.length; j++) {
        if (j != i) {
          deliveryPoints[j].selected = false;
        }
      }
      update();
    });
}

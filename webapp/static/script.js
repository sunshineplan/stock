checkTime = function () {
  var date = new Date();
  var hour = date.getUTCHours();
  var day = date.getDay();
  if (hour >= 1 && hour <= 8 && day >= 1 && day <= 5)
    return true;
  return false;
}

addColor = function (stock, val) {
  if (stock !== null && stock.name != 'n/a') {
    var last = parseFloat(stock.last);
    switch (val) {
      case 'change':
      case 'percent':
        if (parseFloat(stock.change) > 0) return { color: 'red' };
        else if (parseFloat(stock.change) < 0) return { color: 'green' };
      case 'now':
        return color(last, stock.now);
      case 'high':
        return color(last, stock.high);
      case 'low':
        return color(last, stock.low);
      case 'open':
        return color(last, stock.open);
    };
  };
}

color = function (last, value) {
  if (last < parseFloat(value)) return { color: 'red' };
  else if (last > parseFloat(value)) return { color: 'green' };
}

gotoStock = function (stock) {
  window.location = `/stock/${stock.index}/${stock.code}`;
}

timeLabels = function (start, end) {
  var times = [];
  for (var i = 0; start <= end; i++) {
    times[i] = `${Math.floor(start / 60).toString().padStart(2, '0')}:${(start % 60).toString().padStart(2, '0')}`;
    start++;
  }
  return times;
}

$(document).on('click', '#login', () => {
  if ($('#username').val() != 'admin')
    localStorage.setItem('username', $('#username').val());
});

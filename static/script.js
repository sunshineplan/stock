autocomplete = {
  source: (request, response) => {
    fetch('/suggest?' + new URLSearchParams({ keyword: request.term }))
      .then(response => response.json()).then(data => {
        if (!data)
          response(['No matches found.']);
        else
          response($.map(data, item => {
            return `${item.Index}:${item.Code} ${item.Name} ${item.Type}`;
          }));
      });
  },
  select: (event, ui) => {
    if (ui.item.value == 'No matches found.')
      event.preventDefault();
    else {
      var stock = ui.item.value.split(' ')[0].split(':');
      window.location.replace(`/stock/${stock[0]}/${stock[1]}`);
    };
  },
  minLength: 2,
  autoFocus: true,
  position: {
    of: '.search'
  }
};

$(document).on('click', '#login', () => {
  if ($('#username').val() != 'admin')
    localStorage.setItem('username', $('#username').val());
});

checkTime = function () {
  var date = new Date();
  var hour = date.getUTCHours();
  var day = date.getDay();
  if (hour >= 1 && hour <= 8 && day >= 1 && day <= 5) {
    return true;
  };
  return false;
};

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
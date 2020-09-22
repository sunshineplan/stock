index = $('.code').data('index');
code = $('.code').text();

if (code != 'n/a') {
  update_chart(index, code);
  update_realtime(index, code);
  setInterval(() => update_chart(index, code, ct = true), 60000);
  setInterval(() => update_realtime(index, code, ct = true), 3000);
};

$.get('/star', data => {
  if (data == '1') {
    $('.star').addClass('stared');
    $('.star').text('star');
  };
});

function timeLabels(start, end) {
  var times = [];
  for (var i = 0; start <= end; i++) {
    times[i] = `${Math.floor(start / 60).toString().padStart(2, '0')}:${(start % 60).toString().padStart(2, '0')}`;
    start++;
  }
  return times;
}

Chart.defaults.global.maintainAspectRatio = false;
Chart.defaults.global.legend.display = false;
Chart.defaults.global.hover.mode = 'index';
Chart.defaults.global.hover.intersect = false;
Chart.defaults.global.tooltips.mode = 'index';
Chart.defaults.global.tooltips.intersect = false;
Chart.defaults.global.tooltips.displayColors = false;
Chart.defaults.global.animation.duration = 0;
chart = new Chart($('#chart'), {
  type: 'line',
  data: {
    labels: timeLabels(9 * 60 + 30, 11 * 60 + 30).concat(timeLabels(13 * 60 + 1, 15 * 60)),
    datasets: [
      {
        label: 'Price',
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        borderColor: 'red',
        backgroundColor: 'red',
        pointRadius: 0,
        pointHoverRadius: 3
      }
    ]
  },
  options: {
    scales: {
      xAxes: [{
        gridLines: {
          drawTicks: false
        },
        ticks: {
          padding: 10,
          autoSkipPadding: 100,
          maxRotation: 0
        }
      }],
      yAxes: [{
        gridLines: {
          drawTicks: false
        },
        ticks: {
          padding: 12
        }
      }]
    },
    annotation: {
      annotations: [
        {
          id: 'PreviousClose',
          type: 'line',
          mode: 'horizontal',
          scaleID: 'y-axis-0',
          borderColor: 'black',
          borderWidth: .75
        }
      ]
    }
  }
});

$(document).on('click', '.star', function () {
  if ($(this).hasClass('stared'))
    $.post('/star', { action: 'unstar' }, () => {
      $('.star').removeClass('stared');
      $('.star').removeClass('unstar');
      $('.star').text('star_border');
    });
  else
    $.post('/star', () => {
      $('.star').addClass('stared');
      $('.star').text('star');
    });
});
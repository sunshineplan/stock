Chart.defaults.global.maintainAspectRatio = false
Chart.defaults.global.legend.display = false
Chart.defaults.global.hover.mode = 'index'
Chart.defaults.global.hover.intersect = false
Chart.defaults.global.tooltips.mode = 'index'
Chart.defaults.global.tooltips.intersect = false
Chart.defaults.global.tooltips.displayColors = false
Chart.defaults.global.animation.duration = 0

intraday = {
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
        gridLines: { drawTicks: false },
        ticks: {
          padding: 10,
          autoSkipPadding: 100,
          maxRotation: 0
        }
      }],
      yAxes: [{
        gridLines: { drawTicks: false },
        ticks: { padding: 12 }
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
}

const chart = { template: "<canvas class='chart' id='stockChart'></canvas>" }

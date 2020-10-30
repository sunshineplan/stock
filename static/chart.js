Chart.defaults.global.maintainAspectRatio = false
Chart.defaults.global.legend.display = false
Chart.defaults.global.hover.mode = 'index'
Chart.defaults.global.hover.intersect = false
Chart.defaults.global.tooltips.mode = 'index'
Chart.defaults.global.tooltips.intersect = false
Chart.defaults.global.tooltips.displayColors = false
Chart.defaults.global.animation.duration = 0

chart = new Chart(document.getElementById('chart'), {
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
})

new Vue({
  data: {
    index: realtime.index,
    code: realtime.code,
  },
  created() { this.start() },
  methods: {
    start: function () {
      if (this.code != 'n/a') {
        this.load(this.index, this.code);
        setInterval(() => this.load(this.index, this.code, ct = true), 60000);
      }
    },
    load: function (index, code, ct = false) {
      if (checkTime() || !ct) {
        fetch('/get?' + new URLSearchParams({ index: index, code: code, q: 'chart' }))
          .then(response => response.json()).then(json => {
            if (json !== null) {
              if (json.chart != null)
                chart.data.datasets.forEach(dataset => {
                  dataset.data = json.chart;
                });
              if (json.last != 0) {
                chart.options.scales.yAxes[0].ticks.suggestedMin = json.last / 1.01;
                chart.options.scales.yAxes[0].ticks.suggestedMax = json.last * 1.01;
                chart.annotation.options.annotations[0].value = json.last;
              }
              chart.update();
            };
          });
      };
    }
  }
})

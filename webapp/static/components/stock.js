const stock = {
  data() {
    return {
      refresh: Number(app.dataset.refresh) + 1,
      autoUpdate: [],
      stock: {},
      chart: ''
    }
  },
  computed: {
    index() { return this.$route.params.index },
    code() { return this.$route.params.code }
  },
  provide() { return { Stock: Vue.computed(() => this.stock) } },
  watch: {
    $route(to) {
      if (to.name == 'stock' && this.code != 'n/a')
        this.loadChart(true).then(() => this.loadRealtime(true))
    }
  },
  components: { autocomplete, realtime, chart },
  template: `
<div class='content'>
  <header>
    <autocomplete></autocomplete>
    <div class='home' @click="$router.push('/')">
      <div class='icon'>
        <i class='material-icons'>home</i>
      </div>
      <a>Home</a>
    </div>
  </header>
  <realtime></realtime>
  <chart></chart>
</div>`,
  mounted() {
    document.title = 'My Stocks'
    this.start()
    this.chart = new Chart(stockChart, intraday)
  },
  beforeUnmount() {
    for (; this.autoUpdate.length > 0;) clearInterval(this.autoUpdate.pop())
    this.chart.destroy()
  },
  methods: {
    start: function () {
      if (this.code != 'n/a') {
        this.loadChart(true).then(() => this.loadRealtime(true))
        this.autoUpdate.push(setInterval(() => this.loadRealtime(true), this.refresh * 1000))
        this.autoUpdate.push(setInterval(() => this.loadChart(true), 60000))
      }
    },
    loadRealtime: function (force) {
      if (checkTime() || force && this.code)
        return post('/get', { index: this.index, code: this.code, q: 'realtime' })
          .then(response => response.json())
          .then(stock => {
            this.stock = stock
            if (stock && stock.name != 'n/a') {
              document.title = `${stock.name} ${stock.now} ${stock.percent}`
              var last = this.chart.data.datasets[0].data
              if (last.length) {
                last[last.length - 1].y = stock.now
                this.chart.update()
              }
            }
          })
    },
    loadChart: function (force) {
      if (checkTime() || force && this.code) {
        return post('/get', { index: this.index, code: this.code, q: 'chart' })
          .then(response => response.json()).then(json => {
            if (json) {
              if (json.chart)
                this.chart.data.datasets.forEach(dataset => {
                  dataset.data = json.chart
                })
              if (json.last != 0) {
                this.chart.options.scales.yAxes[0].ticks.suggestedMin = json.last / 1.01
                this.chart.options.scales.yAxes[0].ticks.suggestedMax = json.last * 1.01
                this.chart.annotation.elements.PreviousClose.options.value = json.last
              }
              this.chart.update()
            }
          })
      }
    }
  }
}

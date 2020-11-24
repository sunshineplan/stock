const stock = {
  provide() { return { Stock: Vue.computed(() => this.stock) } },
  data() {
    return {
      refresh: Number(app.dataset.refresh) + 1,
      autoUpdate: [],
      stock: {},
      chart: '',
      data: [],
      update: ''
    }
  },
  computed: {
    index() { return this.$route.params.index },
    code() { return this.$route.params.code }
  },
  watch: {
    update(_, last) { if (last != '') this.updateChart() },
    async $route(to) { if (to.name == 'stock' && this.code != 'n/a') await this.load() }
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
    <realtime></realtime>
  </header>
  <chart></chart>
</div>`,
  mounted() {
    document.title = 'My Stocks'
    this.start()
  },
  beforeUnmount() {
    for (; this.autoUpdate.length > 0;) clearInterval(this.autoUpdate.pop())
    this.chart.destroy()
  },
  methods: {
    async start() {
      this.chart = new Chart(stockChart, intraday)
      if (this.code != 'n/a') {
        await this.load()
        this.autoUpdate.push(setInterval(this.loadRealtime, this.refresh * 1000))
        this.autoUpdate.push(setInterval(this.loadChart, 60000))
      }
    },
    async load() {
      await this.loadRealtime(true)
      await this.loadChart(true)
      this.chart.options.scales.yAxes[0].ticks.suggestedMin = this.stock.last / 1.01
      this.chart.options.scales.yAxes[0].ticks.suggestedMax = this.stock.last * 1.01
      this.chart.annotation.elements.PreviousClose.options.value = this.stock.last
      this.updateChart()
    },
    async loadRealtime(force) {
      if (checkTime() || force && this.code) {
        let response = await post('/get', { index: this.index, code: this.code, q: 'realtime' })
        let stock = await response.json()
        if (stock.name) {
          this.stock = stock
          this.update = stock.update
          document.title = `${stock.name} ${stock.now} ${stock.percent}`
        }
      }
    },
    async loadChart(force) {
      if (checkTime() || force && this.code) {
        let response = await post('/get', { index: this.index, code: this.code, q: 'chart' })
        let json = await response.json()
        if (json.chart) this.data = json.chart
      }
    },
    updateChart() {
      if (this.data.length && this.stock.now) {
        let data = this.data
        data[data.length - 1].y = this.stock.now
        this.chart.data.datasets[0].data = data
        this.chart.update()
      }
    }
  }
}

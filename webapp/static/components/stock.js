const stock = {
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
    update() {
      document.title = `${this.stock.name} ${this.stock.now} ${this.stock.percent}`
      updateChart()
    }
  },
  provide() { return { Stock: Vue.computed(() => this.stock) } },
  watch: {
    async $route(to) {
      if (to.name == 'stock' && this.code != 'n/a') {
        await this.loadChart(true)
        await this.loadRealtime(true)
      }
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
        await this.loadChart(true)
        await this.loadRealtime(true)
        document.title = `${this.stock.name} ${this.stock.now} ${this.stock.percent}`
        this.chart.options.scales.yAxes[0].ticks.suggestedMin = this.stock.last / 1.01
        this.chart.options.scales.yAxes[0].ticks.suggestedMax = this.stock.last * 1.01
        this.chart.annotation.elements.PreviousClose.options.value = this.stock.last
        this.updateChart()
        //this.autoUpdate.push(setInterval(this.loadRealtime, this.refresh * 1000))
        //this.autoUpdate.push(setInterval(this.loadChart, 60000))
        this.autoUpdate.push(setInterval(this.loadRealtime, 500, true))
        this.autoUpdate.push(setInterval(this.loadChart, 500, true))
      }
    },
    async loadRealtime(force) {
      if (checkTime() || force && this.code) {
        let response = await post('/get', { index: this.index, code: this.code, q: 'realtime' })
        this.stock = await response.json()
        this.update = this.stock.update
      }
    },
    async loadChart(force) {
      if (checkTime() || force && this.code) {
        let response = await post('/get', { index: this.index, code: this.code, q: 'chart' })
        let json = await response.json()
        this.data = json.chart
      }
    },
    updateChart() {
      if (this.data.length) {
        let data = this.data
        data[data.length - 1].y = this.stock.now
        this.chart.data.datasets[0].data = data
        this.chart.update()
      }
    }
  }
}

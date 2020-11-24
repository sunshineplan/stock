const stocks = {
  components: { autocomplete },
  data() {
    return {
      columns: {
        '指数': 'index',
        '代码': 'code',
        '名称': 'name',
        '最新': 'now',
        '涨跌': 'change',
        '涨幅': 'percent',
        '最高': 'high',
        '最低': 'low',
        '开盘': 'open',
        '昨收': 'last'
      },
      stocks: [],
      sortable: '',
      refresh: Number(app.dataset.refresh) + 1,
      autoUpdate: '',
      fetching: ''
    }
  },
  template: `
<div class='content'>
  <header>
    <autocomplete></autocomplete>
  </header>
  <div class='table-responsive'>
    <table class='table table-hover table-sm'>
      <thead>
        <tr>
          <th v-for='(val, key) in columns'>{{ key }}</th>
        </tr>
      </thead>
      <tbody id='sortable'>
        <tr v-for='stock in stocks' :key='stock.index+stock.code' @click='gotoStock(stock)'>
          <td v-for='val in columns' :style='addColor(stock, val)'>{{ stock[val] }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>`,
  created() { this.start() },
  mounted() {
    document.title = 'My Stocks'
    this.sortable = new Sortable(sortable, {
      animation: 150,
      delay: 500,
      swapThreshold: 0.5,
      onStart: () => this.stop(),
      onEnd: () => this.start(),
      onUpdate: this.onUpdate
    })
  },
  beforeUnmount() {
    this.stop()
    this.sortable.destroy()
  },
  methods: {
    start() {
      this.load(true)
      this.autoUpdate = setInterval(this.load, this.refresh * 1000)
    },
    stop() {
      this.fetching.abort()
      clearInterval(this.autoUpdate)
    },
    load(force) {
      if (checkTime() || force) {
        this.fetching = new AbortController()
        fetch('/mystocks', { signal: this.fetching.signal })
          .then(response => response.json())
          .then(json => { this.stocks = json })
      }
    },
    onUpdate: function (evt) {
      post('/reorder', {
        old: `${this.stocks[evt.oldIndex].index} ${this.stocks[evt.oldIndex].code}`,
        new: `${this.stocks[evt.newIndex].index} ${this.stocks[evt.newIndex].code}`
      })
    }
  }
}

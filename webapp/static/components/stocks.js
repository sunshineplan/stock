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
        <tr v-for='stock in stocks' @click='gotoStock(stock)'>
          <td v-for='val in columns' :style='addColor(stock, val)'>{{ stock[val] }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>`,
  created() { this.start() },
  mounted() {
    document.title = 'My Stocks'
    this.sortable = $('#sortable').sortable({
      start: () => this.stop(),
      stop: () => window.location = '/',
      update: (event, ui) => this.reorder(ui)
    })
  },
  beforeUnmount() {
    this.stop()
    $('#sortable').sortable('destroy')
  },
  methods: {
    start: function () {
      this.load(true)
      this.autoUpdate = setInterval(() => this.load(), this.refresh * 1000)
    },
    stop: function () {
      this.fetching.abort()
      clearInterval(this.autoUpdate)
    },
    load: function (force) {
      if (checkTime() || force) {
        this.fetching = new AbortController()
        fetch('/mystocks', { signal: this.fetching.signal })
          .then(response => response.json())
          .then(json => { this.stocks = json })
      }
    },
    reorder: ui => {
      var orig, dest
      orig = ui.item.find('td')[0].textContent + ' ' + ui.item.find('td')[1].textContent
      if (ui.item.prev().length != 0)
        dest = ui.item.prev().find('td')[0].textContent + ' ' + ui.item.prev().find('td')[1].textContent
      else dest = '#TOP_POSITION#'
      post('/reorder', { orig, dest })
    }
  }
}

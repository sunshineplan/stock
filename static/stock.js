Vue.component("mystocks", {
  template: `
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
`,
  props: {
    stocks: Array,
    columns: Object,
    sortable: Object
  },
  mounted() { $('#sortable').sortable(this.sortable) },
  methods: {
    addColor: addColor,
    gotoStock: gotoStock
  }
})

mystocks = new Vue({
  el: "#mystocks",
  data: {
    Columns: {
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
    Stocks: [],
    Sortable: {
      start: () => mystocks.stop(),
      stop: () => window.location = '/',
      update: (event, ui) => mystocks.reorder(ui)
    },
    autoUpdate: '',
    fetching: ''
  },
  created() { this.start() },
  methods: {
    start: function () {
      this.load();
      this.autoUpdate = setInterval(() => this.load(ct = true), 3000);
    },
    stop: function () {
      this.fetching.abort();
      clearInterval(this.autoUpdate);
    },
    load: function (ct = false) {
      if (checkTime() || !ct) {
        this.fetching = new AbortController();
        var signal = this.fetching.signal;
        fetch('/mystocks', { signal: signal })
          .then(response => response.json())
          .then(data => { this.Stocks = data });
      };
    },
    reorder: ui => {
      var orig, dest;
      orig = ui.item.find('td')[0].textContent + ' ' + ui.item.find('td')[1].textContent;
      if (ui.item.prev().length != 0)
        dest = ui.item.prev().find('td')[0].textContent + ' ' + ui.item.prev().find('td')[1].textContent;
      else dest = '#TOP_POSITION#';
      fetch('/reorder', { method: 'POST', body: new URLSearchParams({ orig: orig, dest: dest }) });
    }
  },
  beforeDestroy() { this.stop() }
})

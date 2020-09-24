Vue.component("realtime", {
  template: `
<div>
  <div style='display: flex; font-size: 2rem;'>
    <i class='material-icons star' :class='stared ? "stared" : ""' @click='star'>{{ stared ? 'star' : 'star_border' }}</i>
    <span>{{ stock.name }}</span>(<span>{{ stock.code }}</span>)&nbsp;&nbsp;&nbsp;
    <span :style='addColor(stock, "now")'>{{ stock.now }}</span>&nbsp;&nbsp;&nbsp;
    <span :style='addColor(stock, "percent")'>{{ stock.percent }}</span>
  </div>
  <div style='min-height: 52px;'>
    <table style='float: left; table-layout: fixed;' :style='{ width: width }'>
      <tbody>
        <tr>
          <td>昨收: <span>{{ stock.last }}</span></td>
          <td>涨跌: <span :style='addColor(stock, "change")'>{{ stock.change }}</span></td>
          <td>涨幅: <span :style='addColor(stock, "percent")'>{{ stock.percent }}</span></td>
        </tr>
        <tr>
          <td>最高: <span :style='addColor(stock, "high")'>{{ stock.high }}</span></td>
          <td>最低: <span :style='addColor(stock, "low")'>{{ stock.low }}</span></td>
          <td>开盘: <span :style='addColor(stock, "open")'>{{ stock.open }}</span></td>
        </tr>
      </tbody>
    </table>
    <table v-if='stock.sell5 !== null || stock.buy5 !== null'>
      <tbody>
        <tr>
          <td>
            <span style='display: inline-flex'>
              卖盘:&nbsp;
              <div class='buysell' style='color: red' v-for='sell in stock.sell5'>{{ sell[0] }}-{{ sell[1] }}</div>
            </span>
          </td>
        </tr>
        <tr>
          <td>
            <span style='display: inline-flex'>
              买盘:&nbsp;
              <div class='buysell' style='color: green' v-for='buy in stock.buy5'>{{ buy[0] }}-{{ buy[1] }}</div>
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <small>更新时间: <span class='update'>{{ stock.update }}</span></small>
</div>
`,
  props: { stock: Object },
  data() { return { stared: false } },
  created() {
    fetch('/star').then(response => response.text())
      .then(data => { if (data == '1') this.stared = true });
  },
  computed: {
    width: function () {
      if (this.stock.sell5 === null && this.stock.buy5 === null) return '480px';
      else return '360px';
    }
  },
  methods: {
    star: function () {
      if (this.stared)
        fetch('/star', { method: 'POST', body: new URLSearchParams({ action: 'unstar' }) })
          .then(() => this.stared = false)
      else
        fetch('/star', { method: 'POST' }).then(() => this.stared = true);
    },
    addColor: addColor
  }
})

realtime = new Vue({
  el: "#realtime",
  data: {
    index: document.getElementById('realtime').dataset.index,
    code: document.getElementById('realtime').dataset.code,
    Stock: {}
  },
  created() { this.start() },
  methods: {
    start: function () {
      this.load();
      setInterval(() => this.load(ct = true), 3000);
    },
    load: function (ct = false) {
      if (checkTime() || !ct)
        fetch('/get?' + new URLSearchParams({ index: this.index, code: this.code, q: 'realtime' }))
          .then(response => response.json())
          .then(stock => {
            this.Stock = stock;
            if (stock !== null && stock.name != 'n/a') {
              document.title = `${stock.name} ${stock.now} ${stock.percent}`;
              var last = chart.data.datasets[0].data;
              if (last.length != 0) {
                last[last.length - 1].y = stock.now;
                chart.update();
              };
            };
          });
    }
  }
})

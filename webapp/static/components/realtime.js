const realtime = {
  inject: ['Stock'],
  data() { return { stared: false } },
  computed: {
    width() {
      if (this.stock.sell5 === null && this.stock.buy5 === null) return '480px'
      else return '360px'
    },
    stock() { return this.Stock.value }
  },
  template: `
<div>
  <div style='display: flex; font-size: 2rem;'>
    <i class='material-icons star' :class='stared ? "stared" : ""' @click='star'>{{ stared ? 'star' : 'star_border' }}</i>
    <span>{{ stock.name }}</span>(<span>{{ stock.code }}</span>)
    <i class='material-icons open' @click='open'>open_in_new</i>&nbsp;&nbsp;&nbsp;
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
              <div class='sellbuy' style='color: red' v-for='sell in stock.sell5'>{{ sell.Price }}-{{ sell.Volume }}</div>
            </span>
          </td>
        </tr>
        <tr>
          <td>
            <span style='display: inline-flex'>
              买盘:&nbsp;
              <div class='sellbuy' style='color: green' v-for='buy in stock.buy5'>{{ buy.Price }}-{{ buy.Volume }}</div>
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <small>更新时间: <span class='update'>{{ stock.update }}</span></small>
</div>`,
  created() {
    fetch('/star').then(response => response.text())
      .then(text => { if (text == '1') this.stared = true })
  },
  methods: {
    star: function () {
      if (this.stared)
        post('/star', { action: 'unstar' })
          .then(() => this.stared = false)
      else post('/star').then(() => this.stared = true)
    },
    open: function () { window.open('http://stockpage.10jqka.com.cn/' + this.stock.code) }
  }
}

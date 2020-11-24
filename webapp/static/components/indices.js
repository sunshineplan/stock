const indices = {
  data() {
    return {
      indices: {},
      names: { '沪': '上证指数', '深': '深证成指', '创': '创业板指', '中': '中小板指' },
      fields: ['now', 'change', 'percent']
    }
  },
  template: `
<div class='indices' v-if='Object.keys(indices).length !== 0'>
  <a v-for='(val, key) in names' :id='key' @click='gotoStock(indices[key])'>
    <span class='short'>{{ key }}</span>
    <span class='full'>{{ val }}</span>
    <span v-for='field in fields' :style='addColor(indices[key], field)'>&nbsp;&nbsp;{{ indices[key][field] }}</span>
  </a>
</div>`,
  created() { this.start() },
  methods: {
    start() {
      this.load(true)
      setInterval(this.load, 10000)
    },
    load(force) {
      if (checkTime() || force) {
        fetch('/indices')
          .then(response => response.json())
          .then(json => { this.indices = json })
      }
    }
  }
}

Vue.component("indices", {
  template: `
<div class='indices' v-if='Object.keys(indices).length !== 0'>
  <a v-for='(val, key) in names' :id='key' @click='gotoStock(indices[key])'>
    <span class='short'>{{ key }}</span>
    <span class='full'>{{ val }}</span>
    <span v-for='field in fields' :style='addColor(indices[key], field)'>&nbsp;&nbsp;{{ indices[key][field] }}</span>
  </a>
</div>
`,
  props: { indices: Object },
  data() {
    return {
      names: { '沪': '上证指数', '深': '深证成指', '创': '创业板指', '中': '中小板指' },
      fields: ['now', 'change', 'percent']
    }
  },
  methods: {
    addColor: addColor,
    gotoStock: gotoStock
  }
})

new Vue({
  el: "#indices",
  data: { Indices: {} },
  created() { this.start() },
  methods: {
    start: function () {
      this.load();
      setInterval(() => this.load(ct = true), 30000);
    },
    load: function (ct = false) {
      if (checkTime() || !ct) {
        fetch('/indices')
          .then(response => response.json())
          .then(json => { this.Indices = json });
      };
    }
  }
})

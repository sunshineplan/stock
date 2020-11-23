const App = Vue.createApp({ data() { return { user: app.dataset.user } } })

const routes = [
  { path: '/', component: stocks },
  { path: '/login', component: login },
  { path: '/setting', component: setting },
  { name: 'stock', path: '/stock/:index/:code', component: stock }
]

const router = VueRouter.createRouter({
  history: VueRouter.createWebHistory(),
  routes
})
App.use(router)

App.mixin({
  methods: {
    addColor(stock, val) {
      if (stock && stock.name != 'n/a') {
        switch (val) {
          case 'change':
          case 'percent':
            return color(stock.change)
          case 'now':
            return color(stock.last, stock.now)
          case 'high':
            return color(stock.last, stock.high)
          case 'low':
            return color(stock.last, stock.low)
          case 'open':
            return color(stock.last, stock.open)
        }
      }
    },
    gotoStock(stock) { this.$router.push(`/stock/${stock.index}/${stock.code}`) },
  }
})

App.component('indices', indices)

App.mount('#app')

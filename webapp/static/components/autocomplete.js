const autocomplete = {
  data() {
    return {
      suggest: '',
      autoComplete: ''
    }
  },
  template: `
<div class='search'>
  <div class='icon'>
    <i class='material-icons'>search</i>
  </div>
  <input v-model.trim='suggest' id='suggest'>
</div>`,
  mounted() {
    this.autoComplete = new autoComplete({
      selector: '#suggest',
      data: {
        src: async () => {
          if (this.suggest.length >= 2) {
            let source = await post('/suggest', { keyword: this.suggest })
            let data = await source.json()
            return data.map(i => `${i.Index}:${i.Code} ${i.Name} ${i.Type}`)
          }
          return []
        },
        cache: false
      },
      searchEngine: (query, record) => { return record },
      placeHolder: 'Search Stock',
      threshold: 1,
      debounce: 300,
      maxResults: 5,
      resultsList: {
        render: true,
        container: source => {
          source.setAttribute('id', 'suggest-list')
          source.setAttribute('class', 'suggest-list')
        }
      },
      resultItem: { content: (data, src) => { src.innerHTML = data.match } },
      noResults: () => {
        let result = document.createElement('li')
        result.setAttribute('class', 'no_result')
        result.setAttribute('tabindex', '1')
        result.innerHTML = 'No Results'
        document.querySelector('#suggest-list').appendChild(result)
      },
      onSelection: feedback => {
        let stock = feedback.selection.value.split(' ')[0].split(':')
        this.$router.push(`/stock/${stock[0]}/${stock[1]}`)
        document.querySelector('#suggest').value = ''
      }
    })
  }
}

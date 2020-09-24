Vue.component("autocomplete", {
  template: `
<div class='search'>
  <div class='icon'>
    <i class='material-icons'>search</i>
  </div>
  <input placeholder='Search Stock' id='search'>
</div>`,
  mounted() {
    $("#search").autocomplete({
      source: (request, response) => {
        fetch('/suggest?' + new URLSearchParams({ keyword: request.term }))
          .then(response => response.json()).then(data => {
            if (!data) response(['No matches found.']);
            else response($.map(data, item => {
              return `${item.Index}:${item.Code} ${item.Name} ${item.Type}`;
            }));
          });
      },
      select: (event, ui) => {
        if (ui.item.value == 'No matches found.') event.preventDefault();
        else {
          var stock = ui.item.value.split(' ')[0].split(':');
          window.location.replace(`/stock/${stock[0]}/${stock[1]}`);
        };
      },
      minLength: 2,
      autoFocus: true,
      position: { of: '.search' }
    })
  },
})

new Vue({ el: '#autocomplete' })

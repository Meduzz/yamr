<inactives>
  <div class="alert alert-warning" if={rows==null||rows.length==0}>No data available.</div>
  <div class="alert alert-danger" if={error!=null}>{error}</div>

  <h3>Inactive domains</h3>
  <table class="table table-striped">
    <thead>
      <tr>
        <th>#</th>
        <th>Domain</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr each={rows}>
        <td>{Id}</td>
        <td>{Name}</td>
        <td>
          <a href="#" onclick={activate}>Activate</a>
        </td>
      </tr>
    </tbody>
  </table>

  <script>
    this.mixin(RestMixin)

    this.rows = []
    this.error = null
    this.page = 0
    this.limit = 20


    const rest = this.initRest('/admin/', 'domains', {Session:opts.session.Id})
    const self = this

    this.on('mount', listInactive)

    activate(e) {
      let request = {
        contentType:'application/json',
        url:"/admin/activate/"+e.item.Id,
        method:"GET",
        success:listInactive,
        error:function(e) {
          self.error = e
          self.update()
        },
        headers:{Session:opts.session.Id}
      }

      $.ajax(request)

      return false
    }

    nextPage(e) {
      self.page++
      listInactive()
      return false
    }

    prevPage(e) {
      self.page--
      if (self.page > -1) {
        listInactive()
      } else {
        self.page = 0
      }
      return false
    }

    function listInactive() {
      rest.list(self.page, self.limit, (rows) => {
        self.rows = rows
        self.update()
      })
    }
  </script>
</inactives>

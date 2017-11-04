<users>
  <div class="alert alert-warning" if={rows==null||rows.length==0}>No data available.</div>
  <div class="alert alert-danger" if={error!=null}>{error}</div>

  <h3>Users</h3>
  <table class="table table-striped">
    <thead>
      <tr>
        <th>#</th>
        <th>Username</th>
        <th>Admin</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      <tr each={rows}>
        <td>{Id}</td>
        <td>{Username}</td>
        <td>{Admin ? 'Yes' : 'No'}</td>
        <td></td>
      </tr>
    </tbody>
  </table>
  <nav>
    <ul class="pager">
      <li class="previous {parent.rows.length==0 || parent.page == 0 ? 'disabled' : ''}">
        <a href="#" onclick={parent.prevPage}><span aria-hidden="true">&larr;</span> Previous</a>
      </li>
      <li class="next {parent.data.length==0 ? 'disabled' : ''}">
        <a href="#" onclick={parent.nextPage}>Next <span aria-hidden="true">&rarr;</span></a>
      </li>
    </ul>
  </nav>


  <script>
    this.mixin(RestMixin)

    this.rows = []
    this.error = null
    this.page = 0
    this.limit = 20

    const rest = this.initRest('/admin/', 'users', {Session:opts.session.Id})
    const self = this

    this.on('mount', listUsers)

    nextPage(e) {
      self.page++
      listUsers()
      return false
    }

    prevPage(e) {
      self.page--
      if (self.page > -1) {
        listUsers()
      } else {
        self.page = 0
      }
      return false
    }

    function listUsers() {
      rest.list(self.page, self.limit, (rows) => {
        self.rows = rows
        self.update()
      })
    }
  </script>
</users>

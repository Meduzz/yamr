<packages>
  <div class="alert alert-warning" if={rows==null||rows.length==0}>No data available.</div>
  <div class="alert alert-danger" if={error!=null}>{error}</div>
  <jsonform method="POST" action="/api/packages">
    <h3>Your packages</h3>
    <table class="table table-striped">
      <thead>
        <tr>
          <th>Package</th>
          <th>Password</th>
          <th>Public</th>
          <th></th<
        </tr>
      </thead>
      <tbody>
        <tr each={parent.rows}>
          <td>{Name}</td>
          <td>{Password}</td>
          <td>{Public}</td>
          <td><a href="#" onclick={parent.parent.edit}>Edit</a></td><!-- <- LOL -->
        </tr>
      </tbody>
      <tfoot>
        <tr>
          <td><input class="form-control" type="text" placeholder="se.kodiak.tools" name="Name"/></td>
          <td><input class="form-control" type="text" placeholder="Password" name="Password"/></td>
          <td><input type="checkbox" name="Public"/></td>
          <td><button type="submit" class="btn btn-default">Save</button></td>
        </tr>
      </tfoot>
    </table>
    <!--
      TODO add pagination
      TODO react to other status codes 400 among others...
    -->
  </jsonform>

  <div class="panel panel-default">
    <div class="panel-body">
      <p>When uploading a jar, it's package must match one of the ones you specified above. Also a basic auth header with your username and the package password are expected, or the upload will be rejected.</p>
      <p>When downloading a jar from a package, that has public set to off, the same basic auth header are expected.</p>
    </div>
  </div>

  <script>
    this.mixin(RestMixin)

    this.rows = []
    this.error = null

    const rest = this.initRest('/api/', 'packages', {Session:opts.session.Id})
    const self = this

    this.settings = function() {
      return {
        success:(ok) => {
          self.tags.jsonform.Name = ''
          self.tags.jsonform.Password = ''

          rest.list(rows => {
            self.rows = rows
            self.update()
          })
        },
        failure:(xhr, status, err) => {
          self.error = err
          self.update()
        },
        fields:{"Name":"text","Password":"text","Public":"boolean"},
        headers:{
          Session:opts.session.Id
        }
      }
    }

    this.on('mount', () => {
      rest.list((rows) => {
        self.rows = rows
        self.update()
      })
    })

    edit(e) {
      this.tags.jsonform.Name.value = e.item.Name
      this.tags.jsonform.Password.value = e.item.Password
      this.tags.jsonform.Public.checked = e.item.Public
      this.update()
      return false
    }
  </script>
</packages>

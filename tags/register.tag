<register>
  <div class="container-fluid">
    <div class="alert alert-danger" if = {error != null}>{error}</div>
    <jsonform method="POST" action="/api/register">
      <div class="col-xs-6">
        <h3>Sign up</h3>
        <div class="form-group">
          <label>Username:</label>
          <input type="text" name="username" placeholder="Username" class="form-control">
        </div>
        <div class="form-group">
          <label>Password:</label>
          <input type="password" name="password" placeholder="Password" class="form-control">
        </div>
        <button class="btn btn-default" type="submit">Sign up</button>
      </div>
      <div class="col-xs-6">
        <div class="panel panel-default">
          Unless you intend to upload and administrate artifacts, you actually do not need to sign up.
          Credentials to download and or upload artifacts can be given to you by the administrator for your domain.
        </div>
      </div>
    </jsonform>
  </div>

  <script>
    const self = this

    this.settings = function() {
      return {
        success:(ok) => {
          self.username = ''
          self.password = ''
          window.location.assign('/#/')
        },
        failure:(xhr, status, err) => {
          self.error = err
          self.update()
        },
        fields:{"username":"text","password":"text"}
      }
    }
  </script>
</register>

<register>
  <div class="container-fluid">
    <div class="alert alert-danger" if = {error != null}>{error}</div>
    <jsonform method="POST" action="/api/register">
      <h3>Sign up</h3>
      <div class="form-group">
        <label>Username:</label>
        <input type="text" name="username" placeholder="Username" class="form-control">
      </div>
      <div class="form-group">
        <label>Password:</label>
        <input type="password" name="password" placeholder="Password" class="form-control">
      </div>
      <div class="form-group">
        <label>Top package:</label>
        <input type="text" name="package" placeholder="Top package" class="form-control">
      </div>
      <button class="btn btn-default" type="submit">Sign up</button>
    </jsonform>
  </div>

  <script>
    const self = this

    this.settings = function() {
      return {
        success:(ok) => {
          self.username = ''
          self.password = ''
          self.package = ''
          window.location.assign('/#/')
        },
        failure:(xhr, status, err) => {
          self.error = err
          self.update()
        },
        fields:{"username":"text","password":"text","package":"text"}
      }
    }
  </script>
</register>

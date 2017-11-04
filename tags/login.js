riot.tag2('login', '<ul class="nav navbar-nav navbar-right" if="{session != null}"> <li if="{session.Admin}"><a href="#/users">Users</a></li> <li if="{session.Admin}"><a href="#/inactive">Inactive</a></li> <li><a href="#/domains">Domains</a></li> <li><a href="/">Log out</a></li> </ul> <div class="navbar-form navbar-right" if="{session == null}"> <jsonform method="post" action="/api/login"> <div class="form-group"> <input type="text" name="username" placeholder="Username" class="form-control"> </div> <div class="form-group"> <input type="password" name="password" placeholder="Password" class="form-control"> </div> <button class="btn btn-default" type="submit"> Sign in </button> <a class="btn btn-default" href="#/register"> Sign up </a> </jsonform> </div>', 'login ul.nav,[data-is="login"] ul.nav{ margin-right: 0px; }', '', function(opts) {
    this.session = null
    this.bus = opts.bus

    const self = this

    this.settings = function() {
      return {
        success:(session) => {
          self.session = session
          self.bus.trigger('session.started', session)
          self.update()
          window.location.assign("#/")
        },
        failure:(xhr, status, err) => {

          alert("Login failed.")
          console.log(err)
          self.password = ''
        },
        fields:{"username":"text","password":"text"}
      }
    }
});

(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);throw new Error("Cannot find module '"+o+"'")}var f=n[o]={exports:{}};t[o][0].call(f.exports,function(e){var n=t[o][1][e];return s(n?n:e)},f,f.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
riot.tag2('apply', '<jsonform method="PUT" action="/api/domain/apply"> <h3>Apply for domain</h3> <div class="form-group"> <label>Reverse domain</label> <input type="text" name="Name" class="form-control"> </div> <div> <button type="submit" class="btn btn-default">Apply</a> </div> </jsonform>', '', '', function(opts) {
    this.error = null
    const self = this

    this.settings = function() {
      return {
        success:(ok) => {
          self.tags.jsonform.Name.value = ''

          window.location.assign("#/domains")
          return false
        },
        failure:(xhr, status, err) => {
          self.error = err
          self.update()
        },
        fields:{"Name":"text"},
        headers:{
          Session:opts.session.Id
        }
      }
    }
});

},{}],2:[function(require,module,exports){
riot.tag2('content', '<register if="{controller == ⁗register⁗}"></register> <profile if="{controller == ⁗profile⁗}"></profile> <search if="{controller == ⁗home⁗}"></search>', '', '', function(opts) {
    this.controller = ""

    let r = riot.route.create()
    const self = this

    r("/register", register)
    r("/profile", profile)
    r("/", home)

    function register() {
      self.controller = "register"
      self.update()
    }

    function profile() {
      self.controller = "profile"
      self.update()
    }

    function home() {
      self.controller = "home"
      self.update()
    }
});

},{}],3:[function(require,module,exports){
riot.tag2('domains', '<div class="alert alert-warning" if="{rows==null||rows.length==0}">No data available.</div> <div class="alert alert-danger" if="{error!=null}">{error}</div> <h3>Your domains</h3> <table class="table table-striped"> <thead> <tr> <th>Domain</th> <th>Status</th> <th></th> </tr> </thead> <tbody> <tr each="{rows}"> <td>{Name}</td> <td>{Active ? \'Verified\' : \'Not verified\'}</td> <td> <a href="#/packages/{Id}" if="{Active}">Packages</a> </td> </tr> </tbody> </table> <nav> <ul class="pager"> <li class="previous {rows.length==0 || page == 0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{prevPage}"><span aria-hidden="true">&larr;</span> Previous</a> </li> <li class="next {rows.length==0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{nextPage}">Next <span aria-hidden="true">&rarr;</span></a> </li> </ul> </nav> <div> <a class="btn btn-default" href="#/apply"><span class="glyphicon glyphicon-plus"></span> Domain</a> </div>', '', '', function(opts) {
    this.mixin(RestMixin)

    this.rows = []
    this.error = null
    this.page = 0
    this.limit = 20

    const rest = this.initRest('/api/', 'domains', {Session:opts.session.Id})
    const self = this

    this.on('mount', listDomains)

    this.nextPage = function(e) {
      self.page++
      listDomains()
      return false
    }.bind(this)

    this.prevPage = function(e) {
      self.page--
      if (self.page > -1) {
        listDomains()
      } else {
        self.page = 0
      }
      return false
    }.bind(this)

    function listDomains() {
      rest.list(self.page, self.limit, (rows) => {
        self.rows = rows
        self.update()
      })
    }
});

},{}],4:[function(require,module,exports){
require('./apply')
require('./content')
require('./domains')
require('./inactives')
require('./jsonform')
require('./login')
require('./main')
require('./packages')
require('./register')
require('./search')
require('./users')
},{"./apply":1,"./content":2,"./domains":3,"./inactives":5,"./jsonform":6,"./login":7,"./main":8,"./packages":9,"./register":10,"./search":11,"./users":12}],5:[function(require,module,exports){
riot.tag2('inactives', '<div class="alert alert-warning" if="{rows==null||rows.length==0}">No data available.</div> <div class="alert alert-danger" if="{error!=null}">{error}</div> <h3>Inactive domains</h3> <table class="table table-striped"> <thead> <tr> <th>#</th> <th>Domain</th> <th></th> </tr> </thead> <tbody> <tr each="{rows}"> <td>{Id}</td> <td>{Name}</td> <td> <a href="#" onclick="{activate}">Activate</a> </td> </tr> </tbody> </table>', '', '', function(opts) {
    this.mixin(RestMixin)

    this.rows = []
    this.error = null
    this.page = 0
    this.limit = 20

    const rest = this.initRest('/admin/', 'domains', {Session:opts.session.Id})
    const self = this

    this.on('mount', listInactive)

    this.activate = function(e) {
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
    }.bind(this)

    this.nextPage = function(e) {
      self.page++
      listInactive()
      return false
    }.bind(this)

    this.prevPage = function(e) {
      self.page--
      if (self.page > -1) {
        listInactive()
      } else {
        self.page = 0
      }
      return false
    }.bind(this)

    function listInactive() {
      rest.list(self.page, self.limit, (rows) => {
        self.rows = rows
        self.update()
      })
    }
});

},{}],6:[function(require,module,exports){
riot.tag2('jsonform', '<form onsubmit="{submitAction}"> <yield></yield> </form>', '', '', function(opts) {
    this.mixin(EntityMixin)

    let action = opts.action
    let method = opts.method

    let settings = this.parent.settings()
    let success = settings.success
    let failure = settings.failure
    let fields = settings.fields
    let headers = settings.headers || {}

    let formReader = this.initEntity({properties:Object.keys(fields).map(field => {
      return {
        read:(ctx, entity) => {
          if (fields[field] == 'text') {
            entity[field] = ctx[field].value
          } else if (fields[field] == 'number') {
            entity[field] = parseInt(ctx[field].value, 10)
          } else {
            entity[field] = ctx[field].checked
          }
          return entity
        }
      }
    })})

    this.submitAction = function(e) {
      let formData = JSON.stringify(formReader.bind(this))
      let settings = {
        contentType:'application/json',
        data:formData,
        dataType:'json',
        url:action,
        method:method,
        success:success,
        error:failure,
        headers:headers
      }

      $.ajax(settings)

      e.preventDefault()
      return false
    }.bind(this)
});

},{}],7:[function(require,module,exports){
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

},{}],8:[function(require,module,exports){
module.exports=require(4)
},{"./apply":1,"./content":2,"./domains":3,"./inactives":5,"./jsonform":6,"./login":7,"./main":8,"./packages":9,"./register":10,"./search":11,"./users":12}],9:[function(require,module,exports){
riot.tag2('packages', '<div class="alert alert-warning" if="{rows==null||rows.length==0}">No data available.</div> <div class="alert alert-danger" if="{error!=null}">{error}</div> <jsonform method="POST" action="/api/packages?id={opts.pageId}"> <h3>Your packages</h3> <table class="table table-striped"> <thead> <tr> <th>Package</th> <th>Password</th> <th>Public</th> <th></th> </tr> </thead> <tbody> <tr each="{parent.rows}"> <td>{Name}</td> <td>{Password}</td> <td>{Public}</td> <td><a href="#" onclick="{parent.parent.edit}">Edit</a></td> </tr> </tbody> <tfoot> <tr> <td><input class="form-control" type="text" placeholder="se.kodiak.tools" name="Name"></td> <td><input class="form-control" type="text" placeholder="top secret" name="Password"></td> <td> <input type="checkbox" name="Public"> <input type="hidden" name="Id"> </td> <td><button type="submit" class="btn btn-default">Save</button></td> </tr> </tfoot> </table> <nav> <ul class="pager"> <li class="previous {parent.rows.length==0 || parent.page == 0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{parent.prevPage}"><span aria-hidden="true">&larr;</span> Previous</a> </li> <li class="next {parent.data.length==0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{parent.nextPage}">Next <span aria-hidden="true">&rarr;</span></a> </li> </ul> </nav> </jsonform> <div class="panel panel-default"> <div class="panel-body"> <p>When uploading a jar, it\'s package must match one of the ones you specified above. Also a basic auth header with your username and the package password are expected, or the upload will be rejected.</p> <p>When downloading a jar from a package, that has public set to off, the same basic auth header are expected.</p> </div> </div>', '', '', function(opts) {
    this.mixin(RestMixin)

    this.rows = []
    this.error = null
    this.page = 0
    this.limit = 20

    const rest = this.initRest('/api/', 'packages', {Session:opts.session.Id})
    const self = this

    this.settings = function() {
      return {
        success:(ok) => {
          self.tags.jsonform.Name.value = ''
          self.tags.jsonform.Password.value = ''

          listPackages()
        },
        failure:(xhr, status, err) => {
          self.error = err
          self.update()
        },
        fields:{"Id":"number","Name":"text","Password":"text","Public":"boolean"},
        headers:{
          Session:opts.session.Id
        }
      }
    }

    this.on('mount', listPackages)

    this.edit = function(e) {
      this.tags.jsonform.Name.value = e.item.Name
      this.tags.jsonform.Password.value = e.item.Password
      this.tags.jsonform.Public.checked = e.item.Public
      this.tags.jsonform.Id.value = e.item.Id
      this.update()
      return false
    }.bind(this)

    this.nextPage = function(e) {
      self.page++
      listPackages()
      return false
    }.bind(this)

    this.prevPage = function(e) {
      self.page--
      if (self.page > -1) {
        listPackages()
      } else {
        self.page = 0
      }
      return false
    }.bind(this)

    function listPackages() {
      rest.list(self.page, self.limit, {"id":opts.pageId}, (rows) => {
        self.rows = rows
        self.update()
      })
    }
});

},{}],10:[function(require,module,exports){
riot.tag2('register', '<div class="container-fluid"> <div class="alert alert-danger" if="{error != null}">{error}</div> <jsonform method="POST" action="/api/register"> <div class="col-xs-6"> <h3>Sign up</h3> <div class="form-group"> <label>Username:</label> <input type="text" name="username" placeholder="Username" class="form-control"> </div> <div class="form-group"> <label>Password:</label> <input type="password" name="password" placeholder="Password" class="form-control"> </div> <button class="btn btn-default" type="submit">Sign up</button> </div> <div class="col-xs-6"> <div class="panel panel-default"> Unless you intend to upload and administrate artifacts, you actually do not need to sign up. Credentials to download and or upload artifacts can be given to you by the administrator for your domain. </div> </div> </jsonform> </div>', '', '', function(opts) {
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
});

},{}],11:[function(require,module,exports){
riot.tag2('search', '<div class="container-fluid"> <div class="alert alert-warning" if="{(data==null||data.length==0) && q.value.length > 0}">No data available.</div> <div class="alert alert-danger" if="{error!=null}">{error}</div> <form onsubmit="{unsubmit}"> <div class="form-group"> <label>Search</label> <input type="text" name="q" class="form-control" onkeyup="{search}"> </div> </form> <table class="table table-striped"> <thead> <tr> <th>Package</th> <th>Artifact</th> <th>Version</th> </tr> </thead> <tbody> <tr each="{data}"> <td>{Group}</td> <td>{Name}</td> <td>{Version}</td> </tr> </tbody> </table> <nav> <ul class="pager"> <li class="previous {data.length==0 || page == 0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{prevPage}"><span aria-hidden="true">&larr;</span> Previous</a> </li> <li class="next {data.length==0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{nextPage}">Next <span aria-hidden="true">&rarr;</span></a> </li> </ul> </nav> </div>', '', '', function(opts) {
    this.bus = opts.bus
    this.data = []
    this.timer = null
    this.page = 0
    this.limit = 20
    self = this

    this.search = function(e) {
      if (this.timer != null) {
        clearTimeout(this.timer)
      }
      this.timer = setTimeout(() => doSearch(e.target.value), 250)
    }.bind(this)

    this.unsubmit = function(e) {
      return false
    }.bind(this)

    this.nextPage = function(e) {
      self.page++
      doSearch(self.q.value)
      return false
    }.bind(this)

    this.prevPage = function(e) {
      self.page--
      if (self.page > -1) {
        doSearch(self.q.value)
      } else {
        self.page = 0
      }
      return false
    }.bind(this)

    function doSearch(query) {
      headers = {}

      if (opts.session != null) {
        headers['Session'] = opts.session.Id
      }

      let settings = {
        dataType:'json',
        url:'/api/search?q='+query+'&page='+self.page+'&limit='+self.limit,
        method:'GET',
        success:searchSuccess,
        error:searchFailure,
        headers:headers
      }

      $.ajax(settings)
    }

    function searchSuccess(ok) {
      self.data = ok
      self.update()
    }

    function searchFailure(xhr, status, err) {
      self.error = err
      self.update()
    }
});

},{}],12:[function(require,module,exports){
riot.tag2('users', '<div class="alert alert-warning" if="{rows==null||rows.length==0}">No data available.</div> <div class="alert alert-danger" if="{error!=null}">{error}</div> <h3>Users</h3> <table class="table table-striped"> <thead> <tr> <th>#</th> <th>Username</th> <th>Admin</th> <th></th> </tr> </thead> <tbody> <tr each="{rows}"> <td>{Id}</td> <td>{Username}</td> <td>{Admin ? \'Yes\' : \'No\'}</td> <td></td> </tr> </tbody> </table> <nav> <ul class="pager"> <li class="previous {parent.rows.length==0 || parent.page == 0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{parent.prevPage}"><span aria-hidden="true">&larr;</span> Previous</a> </li> <li class="next {parent.data.length==0 ? \'disabled\' : \'\'}"> <a href="#" onclick="{parent.nextPage}">Next <span aria-hidden="true">&rarr;</span></a> </li> </ul> </nav>', '', '', function(opts) {
    this.mixin(RestMixin)

    this.rows = []
    this.error = null
    this.page = 0
    this.limit = 20

    const rest = this.initRest('/admin/', 'users', {Session:opts.session.Id})
    const self = this

    this.on('mount', listUsers)

    this.nextPage = function(e) {
      self.page++
      listUsers()
      return false
    }.bind(this)

    this.prevPage = function(e) {
      self.page--
      if (self.page > -1) {
        listUsers()
      } else {
        self.page = 0
      }
      return false
    }.bind(this)

    function listUsers() {
      rest.list(self.page, self.limit, (rows) => {
        self.rows = rows
        self.update()
      })
    }
});

},{}]},{},[4])
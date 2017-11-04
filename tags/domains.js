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

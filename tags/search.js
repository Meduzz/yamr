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

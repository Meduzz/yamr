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

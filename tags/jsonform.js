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

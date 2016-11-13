var EntityMixin = {
  initEntity:function(options) {
    return new EntityHelper(options)
  }
}

function EntityHelper(options) {
  this.options = options
}

EntityHelper.prototype.bind = function(ctx, start) {
  if (start == null) {
    start = {}
  }
  var readFuncs = this.options.properties.map(function(field) { return field.read })

  return readFuncs.reduce(function(last, current) {
    return current(ctx, last)
  }, start)
}

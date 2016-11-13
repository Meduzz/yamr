/**
 * A mixin for riot custom tags to help out with rest calls.
 * It has a very opinionated view of rest, and cant do it all.
 * It uses jquery.ajax to do the heavy lifting.
 */
var RestMixin = {
  initRest:function(apiUrl, entity, headers) {
    return new RestHelper(apiUrl, entity, headers)
  }
}

function RestHelper(apiUrl, entityName, headers) {
  this.api = apiUrl
  this.entity = entityName
  this.headers = headers || {}
}

RestHelper.prototype.list = function(skip, limit, callback) {
  if (callback == null && typeof(skip) === 'function') {
    callback = skip
    skip = 0
    limit = 10
  }

  var setting = {
    contentType:'application/json',
    dataType:'json',
    url: this.api+this.entity+'?skip='+skip+'&limit='+limit,
    success: function(result) {
      callback(result)
    },
    method: 'GET',
    headers:this.headers
  }

  $.ajax(setting)
}

RestHelper.prototype.remove = function(id, callback) {
  var setting = {
    contentType:'application/json',
    dataType:'json',
    url: this.api+this.entity+'/'+id,
    complete: function(result) {
      callback(result)
    },
    method: 'DELETE'
  }

  $.ajax(setting)
}

RestHelper.prototype.create = function(data, callback) {
  var setting = {
    contentType:'application/json',
    data:JSON.stringify(data),
    dataType:'json',
    url: this.api+this.entity,
    success: function(result) {
      callback(result)
    },
    method: 'POST',
    headers:this.headers
  }

  $.ajax(setting)
}

RestHelper.prototype.store = function(id, data, callback) {
  var setting = {
    contentType:'application/json',
    data:JSON.stringify(data),
    dataType:'json',
    url: this.api+this.entity+'/'+id,
    complete:function(xhr, status) {
      callback(status)
    },
    method: 'PUT'
  }

  $.ajax(setting)
}

RestHelper.prototype.load = function(id, callback) {
  var setting = {
    contentType:'application/json',
    dataType:'json',
    url: this.api+this.entity+'/'+id,
    success:function(response) {
      callback(response)
    },
    method: 'GET'
  }

  $.ajax(setting)
}

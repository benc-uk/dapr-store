if (process.env.NODE_ENV == 'test') {
  console.log = function() {}
  console.error = function() {}
  console.warn = function() {}
}
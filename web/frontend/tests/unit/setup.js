// disable console noise
if (process.env.NODE_ENV == 'test') {
  console.log = function() {}
  console.error = function() {}
  console.warn = function() {}
}

createAppDiv()
function createAppDiv() {
  let app = document.createElement('div')
  app.setAttribute('id', 'app')
  document.body.appendChild(app)
}
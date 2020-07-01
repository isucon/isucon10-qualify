const agreedServer = require('agreed-server')
const morgan = require('morgan')
const cors = require('cors')

const path = './.dist.json'
const port = 1323
const middlewares = [
  cors(),
  morgan('dev')
]

const { app, createServer } = agreedServer({ path, port, middlewares })
createServer(app)

console.log(`listen on http://localhost:${port}`)

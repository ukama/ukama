const express = require('express')
const expressWinston = require('express-winston')

function configureExpress(logger) {
 const app = express()
 app.use(expressWinston.logger({ winstonInstance: logger }))
 return app
}

module.exports = {
 configureExpress,
}

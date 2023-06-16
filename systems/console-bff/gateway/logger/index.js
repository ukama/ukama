'use strict'
const winston = require('winston')

const logger = winston.createLogger({
 transports: [new winston.transports.Console()],
 format: winston.format.combine(
  winston.format.colorize({ all: true }),
  winston.format.simple()
 ),
})

module.exports = {
 logger,
}

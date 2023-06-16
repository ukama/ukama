const dotenv = require('dotenv')
dotenv.config({ path: '.env' })
const { ApolloServer } = require('@apollo/server')
const { expressMiddleware } = require('@apollo/server/express4')
const { ApolloServerPluginInlineTrace } = require('apollo-server-core')
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway')
const {
 ApolloServerPluginDrainHttpServer,
} = require('@apollo/server/plugin/drainHttpServer')

const cors = require('cors')
const http = require('http')
const { logger } = require('./logger')
const { json } = require('body-parser')
const { configureExpress } = require('./configureExpress')

const app = configureExpress(logger)
const subgraphUrls = [process.env.SUBGRAPH_PLANNING_TOOL]
const httpServer = http.createServer(app)
const gateway = new ApolloGateway({
 //  supergraphSdl: new IntrospectAndCompose({
 //   subgraphs: [{ name: 'subgraph', url: '' }],
 // }),
 serviceList: subgraphUrls.map((url) => ({ name: 'subgraph', url })),
})

const startServer = async () => {
 const server = new ApolloServer({
  gateway,
  plugins: [
   ApolloServerPluginInlineTrace(),
   ApolloServerPluginDrainHttpServer({ httpServer }),
  ],
 })

 await server.start()

 app.use(
  '/graphql',

  cors({
   origin: [
    process.env.AUTH_APP_URL,
    process.env.PLAYGROUND_URL,
    process.env.CONSOLE_APP_URL,
   ],
   credentials: true,
  }),

  json(),
  expressMiddleware(server)
 )
 await new Promise((resolve) =>
  httpServer.listen({ port: process.env.PORT }, resolve)
 )
 logger.info(`🚀 Server ready at http://localhost:${process.env.PORT}/graphql`)
}

startServer()

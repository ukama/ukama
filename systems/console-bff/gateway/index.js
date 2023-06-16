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

const subgraphUrls = ['http://localhost:4041']
const httpServer = http.createServer(app)
const gateway = new ApolloGateway({
 //  supergraphSdl: new IntrospectAndCompose({
 //   subgraphs: [{ name: 'subgraph', url: 'http://localhost:4041' }],
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
    'https://localhost:4455',
    'http://localhost:3000',
    'https://studio.apollographql.com',
   ],
   credentials: true,
  }),

  json(),
  expressMiddleware(server)
 )
 await new Promise((resolve) => httpServer.listen({ port: 4001 }, resolve))
 logger.info(`🚀 Server ready at http://localhost:4001/graphql`)
}

startServer()

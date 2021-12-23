# Ukama-BFF 

## BFF Structure
- 📂 __ukama\-bff__
   - 📄 [Dockerfile](Dockerfile) for Docker
   - 📄 [Makefile](Makefile)
   - 📄 [README.md](README.md)
   - 📂 __bff\_integration__  for integration test
     - 📄 [Dockerfile](bff_integration/Dockerfile)
     - 📄 [Int.Dockerfile](bff_integration/Int.Dockerfile)
     - 📄 [Makefile](bff_integration/Makefile)
     - 📂 __test__
   - 📄 [jest.config.ts](jest.config.ts)
   - 📄 [list.md](list.md)
   - 📂 __logs__ to maintain all logs
     - 📄 [app.log](logs/app.log)
     - 📄 [errors.log](logs/errors.log)
     - 📄 [exceptions.log](logs/exceptions.log)
   - 📄 [node\_modules](node_modules)
   - 📄 [nodemon.json](nodemon.json)
   - 📄 [package.json](package.json)
   - 📂 __src__
     - 📂 __api__ for axios fetch
       - 📄 [index.ts](src/api/index.ts)
     - 📂 __common__ 
       - 📄 [Authentication.ts](src/common/Authentication.ts)
       - 📄 [createSchema.ts](src/common/createSchema.ts)
       - 📄 [graphql.ts](src/common/graphql.ts)
       - 📄 [index.ts](src/common/index.ts)
       - 📄 [types.ts](src/common/types.ts)
       - 📄 [utils.ts](src/common/utils.ts)
     - 📂 __config__
       - 📄 [configureApolloServer.ts](src/config/configureApolloServer.ts)
       - 📄 [configureExpress.ts](src/config/configureExpress.ts)
       - 📄 [setupLogger.ts](src/config/setupLogger.ts)
     - 📂 __constants__
       - 📄 [endpoints.ts](src/constants/endpoints.ts)
       - 📄 [index.ts](src/constants/index.ts)
     - 📂 __errors__
       - 📄 [base.error.ts](src/errors/base.error.ts)
       - 📄 [codes.ts](src/errors/codes.ts)
       - 📄 [http400.error.ts](src/errors/http400.error.ts)
       - 📄 [http401.error.ts](src/errors/http401.error.ts)
       - 📄 [http404.error.ts](src/errors/http404.error.ts)
       - 📄 [http500.error.ts](src/errors/http500.error.ts)
       - 📄 [index.ts](src/errors/index.ts)
       - 📄 [messages.ts](src/errors/messages.ts)
     - 📄 [index.ts](src/index.ts)
     - 📂 __jobs__
       - 📄 [subscriptionJob.ts](src/jobs/subscriptionJob.ts)
     - 📂 __mockServer__
       - 📄 [index.ts](src/mockServer/index.ts)
       - 📂 __mockData__ to generate random data
       - 📄 [utils.ts](src/mockServer/utils.ts)
     - 📂 __modules__ 
       - 📂 __alert__ 
         - 📄 [interface.ts](src/modules/alert/interface.ts)
         - 📄 [mapper.ts](src/modules/alert/mapper.ts)
         - 📂 __resolver__ for all alert resolvers
         - 📄 [service.ts](src/modules/alert/service.ts)
         - 📂 __tests__ for all alert tests
         - 📄 [types.ts](src/modules/alert/types.ts)
       - 📂 __billing__
         - 📄 [interface.ts](src/modules/billing/interface.ts)
         - 📄 [mapper.ts](src/modules/billing/mapper.ts)
         - 📂 __resolver__ for all billing resolvers
         - 📄 [service.ts](src/modules/billing/service.ts)
         - 📂 __tests__ for all billing tests
         - 📄 [types.ts](src/modules/billing/types.ts)
       - 📂 __data__
         - 📄 [interface.ts](src/modules/data/interface.ts)
         - 📄 [mapper.ts](src/modules/data/mapper.ts)
         - 📂 __resolver__ for all data resolvers
         - 📄 [service.ts](src/modules/data/service.ts)
         - 📂 __tests__ for all data tests
         - 📄 [types.ts](src/modules/data/types.ts)
       - 📂 __esim__
         - 📄 [interface.ts](src/modules/esim/interface.ts)
         - 📄 [mapper.ts](src/modules/esim/mapper.ts)
         - 📂 __resolver__ for all esim resolvers
         - 📄 [service.ts](src/modules/esim/service.ts)
         - 📂 __tests__ for all esim tests
         - 📄 [types.ts](src/modules/esim/types.ts)
       - 📂 __network__
         - 📄 [interface.ts](src/modules/network/interface.ts)
         - 📄 [mapper.ts](src/modules/network/mapper.ts)
         - 📂 __resolver__ for all network resolvers
         - 📄 [service.ts](src/modules/network/service.ts)
         - 📂 __tests__ for all network tests
         - 📄 [types.ts](src/modules/network/types.ts)
       - 📂 __node__
         - 📄 [interface.ts](src/modules/node/interface.ts)
         - 📄 [mapper.ts](src/modules/node/mapper.ts)
         - 📂 __resolver__ for all node resolvers
         - 📄 [service.ts](src/modules/node/service.ts)
         - 📂 __tests__ for all node tests
         - 📄 [types.ts](src/modules/node/types.ts)
       - 📂 __subscriptions__
         - 📂 __resolver__ for all subscriptions
       - 📂 __user__
         - 📄 [interface.ts](src/modules/user/interface.ts)
         - 📄 [mapper.ts](src/modules/user/mapper.ts)
         - 📂 __resolver__ for all user resolvers
         - 📄 [service.ts](src/modules/user/service.ts)
         - 📂 __tests__ for all user tests
         - 📄 [types.ts](src/modules/user/types.ts)
     - 📂 __utils__
       - 📄 [index.ts](src/utils/index.ts)
   - 📄 [tsconfig.json](tsconfig.json)
   - 📄 [yarn.lock](yarn.lock)

## Tools 
1. Node Js (express)
2. Typescript
3. GraphQL
4. Eslint
5. Prettier

## Installation
`yarn install` in directory

## Lint
`yarn lint` in directory

## Build
`yarn build` in directory

## Test
1. `yarn build` in directory

2. `yarn test` in directory

## Production
1. `yarn build` in directory

2. `yarn start` in directory

## Development

1. Clone repository `https://github.com/ukama/ukama-bff.git`

2. `yarn install` in directory

3. `yarn dev` to spin up server in development

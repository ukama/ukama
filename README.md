# Ukama-BFF 

## Architecture Diagram


## How to try

Clone: `https://github.com/ukama/ukama-bff.git`

After cloning run below command:

    yarn 


## Project Folder Structure

    .
    ├── .github                 # Includes the Workflows YML file
    ├── .vscode                 # Includes the VS editor settings file
    ├── bff_integration         # Includes the integration tests
    ├── src                     # Main App folder
        ├── api                 # Contains the axios api fetch method for API Gateway
        ├── common              # Contains the common services and types required in project
        ├── config              # Contains the configuration files for server
        ├── constants           # Contains the constants (enums , variables) used in project
        ├── errors              # Contains the error classes , messages and hanlers
        ├── jobs                # Contains the GraphQL susbscriptions job required in project
        ├── mockServer          # Contains the mocked REST Api server that send dummy data.
        ├── modules             # Contains all the GraphQL resolvers and their services respectively
        ├── utils               # Includes all utility methods.
    ├── eslintrc                # Eslint configrations and Rules.
    ├── prettierrc              # Prettier configrations.
    ├── Dockerfile              # Docker configrations rules.
    ├── jest.config             # jest configrations and Rules.
    ├── nodemon                 # nodemon configrations and Rules.
    ├── tsconfig.json           # typescript configrations and Rules.
    └── package.json            # Include app metadata, dependencies and scripts.

## Scripts

    yarn lint

This script will tell all linting issues.

    yarn dev

This script will start the project in development.

    yarn build

This script will build the project.

    yarn test

This script will run unit tests.

### Note: Build project before test.

    yarn start

This script will start the project in production.

### Note: Build project before test.

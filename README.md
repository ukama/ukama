# Ukama-BFF 

## Architecture Diagram

<img width="571" alt="Screenshot 2021-12-23 at 5 17 34 PM" src="https://user-images.githubusercontent.com/61826215/147475970-054a9f33-dd46-4d47-b96c-5f9a097b252f.png">



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

## Contribute

### Dev Container
The repository supports [dev container](https://code.visualstudio.com/docs/remote/containers). To start working, clone the repo and reopen it in VS Remote Container. Mase user you have [Remote Development extenstion installed](https://aka.ms/vscode-remote/download/extension)

[![VS Code Container](https://img.shields.io/static/v1?label=VS+Code&message=Container&logo=visualstudiocode&color=007ACC&logoColor=007ACC&labelColor=2C2C32)](https://open.vscode.dev/microsoft/vscode)


### Otherwise

Clone repository `git clone https://github.com/ukama/ukama-bff.git` and run `yarn`

### Scripts 
    `yarn lint`

This script will tell all linting issues.

    `yarn dev`

This script will start the project in development.

    `yarn build`

This script will build the project.

    `yarn test`

This script will run unit tests.


    `yarn start`

This script will start the project in production.


# Ukama-DASHBOARD 

## Architecture Diagram
## How to try
- Clone ukamaDashboard repo.
- ` git clone https://github.com/ukama/ukama-dashboard.git`
- Install the dependencies
- `yarn && yarn start`.
- Start the project 
- `yarn start`

### Project Folder Structure

    .
    ├── .github                 # Includes the Workflows YML file for github actions
    ├── .husky                  # Includes the Pre-commit hook script to validate our commits before pushing to remote.
    ├── .vscode                 # Includes the VS editor settings file.
    ├── nginx                   # 
    ├── public                  # Root folder that gets served up as our react app. This includes the App icon and [Fonts](https://github.com/ukama/ukama-dashboard/blob/37c2ff8b5f1749ed95c395af18383c05e0467275/public/index.html#L48-L51).
    ├── src                     # Main App folder
        ├── api                 # Contains the apollo client setup and .graphql in which we define all the queries and mutations.
        ├── assets              # Contains all the resouces(i.e Icons, SVGs).
        ├── components          # Include collection of UI components, which build the pages.
        ├── constants           # Include static data (Objects, Arrays and Columns).
        ├── generated           # This is codegen generated file (Don't edit it manually).
        ├── helpers             # This will include all the helpers method.
        ├── i18n                # Used for multi language support.
        ├── layout              # Main app layout with Header and Footer.
        ├── pages               # Include all the pages here.
        ├── recoil              # Recoil for state management.
        ├── router              # App routing.
        ├── styles              # Include global style.
        ├── theme               # App theming and components customizations here.
        ├── types               # Includes global types.
        ├── utils               # Includes all utility methods.
    └── eslintrc                # Eslint configrations and Rules.
    └── prettierrc              # Prettier configrations.
    └── codegen.yml             # Codegen configration file to create hooks based on our Queries and Mutations.
    └── package.json            # Include app metadata, dependencies and scripts.


### Scripts

    yarn generate

If you updated the `src/api/index.graphql` file, You need to run the above script to generate/update the hooks.

    yarn start

This script will be used to run the app.

    yarn pretty

This script go through all files defined in `src` and format them accoring to prettierrc configrations.

    yarn lint:fix

This script will fix all possibly fixable linting issues.

    yarn check-types

To validate typescript types run this script.

    yarn no-console

Checks if there's any console log statement in app which developer forgot to remove.

    yarn pre-commit-checks

Above script runs before commit to validate the commit standerds and possible build failures.

#!/bin/bash

# Setup up gitignore
wget https://gist.githubusercontent.com/phoobynet/a77e757bfe7e74c159e6161c3f580060/raw/.gitignore -O .gitignore

# Set engines to node 18 or greater (to shut eslint up)
npm pkg set engines.node=">=18"

# Due to some issues with legacy dependencies with stylelint, we need to create some additional configuration
touch .npmrc
echo "legacy-peer-deps=true" > .npmrc

# Install some basics, in this case Node types and support for SASS/SCSS files.
npm i -D @types/node sass

# Install ESLint dependencies
npm i -D eslint eslint-config-prettier eslint-config-standard-with-typescript eslint-config-stylelint eslint-plugin-import eslint-plugin-n eslint-plugin-promise eslint-plugin-vue @typescript-eslint/eslint-plugin

# Install prettier (with Tailwind plugin - remove if not required)
npm i -D prettier eslint-config-prettier @trivago/prettier-plugin-sort-imports prettier-plugin-tailwindcss 

# Install stylelint
npm i -D stylelint eslint-config-stylelint stylelint-config-standard-scss stylelint-config-recommended-scss stylelint-order stylelint-config-idiomatic-order postcss-html stylelint-config-recommended-vue stylelint-config-css-modules stylelint-config-prettier stylelint-prettier

# Change the configuration of vite.config.ts to include a default alias
wget https://gist.githubusercontent.com/phoobynet/c3c047042260d324841670d08c7dee53/raw/vite.config.ts -O vite.config.ts

# IMPORTANT - Update tsconfig.json to include "baseUrl" and "paths" to match the ./src -> @
npx --yes json -I -f tsconfig.json -e 'this.compilerOptions.baseUrl=".";this.compilerOptions.paths={"@/*":["src/*"]}'

# Enable ESLint prettier configuration
wget https://gist.githubusercontent.com/phoobynet/54f6ef8ff09750a1c17f45096a776af2/raw/.eslintrc.cjs -O .eslintrc.cjs

# Download prettier configuration
wget https://gist.githubusercontent.com/phoobynet/95356c313f2cd8c98059f1053dbda4c9/raw/.prettierrc.cjs

# Download Stylelint configuration
wget https://gist.githubusercontent.com/phoobynet/bc88d7b1a3c62354ea037602868f81a4/raw/.stylelintrc.cjs

rm src/style.css

touch src/style.scss

npx --yes replace-in-file "import App from './App.vue'" "import App from '@/App.vue'" src/main.ts
npx replace-in-file "import './style.css'" "import '@/style.scss'" src/main.ts

rm src/components/HelloWorld.vue
rm src/assets/vue.svg
rm src/App.vue

touch src/App.vue

echo "<script lang=\"scss\" setup></script><template><div>TODO</div></template><style lang=\"scss\" scoped></style>" > src/App.vue

npx prettier --write src/**

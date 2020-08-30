# StackHead Project validator

This tool validates your StackHead project definition files.

The schema is provided as [JSON Schema](https://json-schema.org/).

## Installation

The validator is a binary which you can install via various package managers.

### [Composer (PHP)](https://packagist.org/packages/getstackhead/project-validator)

```shell script
composer require getstackhead/project-validator
```

Binary is located at `vendor/bin/project-validator`.

### [NPM (NodeJS)](https://www.npmjs.com/package/@getstackhead/project-validator)

```shell script
npm i --save-dev @getstackhead/project-validator
```

Binary is located at `./node_modules/.bin/project-validator`.

## Usage

Simply add the path to the project definition file you want to validate:

```shell script
./bin/project-validator path/to/definition.yml
```

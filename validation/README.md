# StackHead Project validator

This tool validates your StackHead project definition files.

The schema is provided as [JSON Schema](https://json-schema.org/).

## How to use

The validator is a binary which is available for different package managers.

You'll find example project 

Find valid and invalid files in examples/ folder.

In case you're looking for the source code of the validator, it's in src.

### Standalone

```shell script
./bin/project-validator path/to/definition.yml
```

### [PyPI (Python)](https://pypi.org/project/stackhead-project-validator)

```shell script
pip install stackhead-project-validator
```

Find the location of the package with `pip show stackhead-project-validator`. (Binary path as above.)

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
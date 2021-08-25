# prew

## what is this?

prew is a simple tool to manage your python code.

## project init

`prew init`  

when prew init this path, create virtualenv.  

## package

### install

`prew install <package> [<version>]`  

install a package in the virtualenv.

### remove

`prew remove <package>`

remove the package in the virtualenv.

### spec.yaml

install and remove command is record or delete dependencies in spec.yaml.

## run

`prew run`

this command is run python code in current path.  

when run this command, it will install dependencies in spec.yaml.

## restore

`prew restore`

restore command is install and init virtualenv and install dependencies based on spec.yaml.

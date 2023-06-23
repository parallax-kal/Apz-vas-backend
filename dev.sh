#!/bin/bash

# check if there is go installed and install it if not

if ! [ -x "$(command -v go)" ]; then
  echo 'Error: go is not installed.' >&2
  echo 'Installing go'
  # if you are using Ubuntu or any other distribution that has apt-get, run this
  if [ -x "$(command -v apt-get)" ]; then
    sudo apt-get install golang-go
  fi

  #  macos run this

  if [ -x "$(command -v brew)" ]; then
    brew install golang
  fi

  # if on windows run this

  if [ -x "$(command -v choco)" ]; then
    choco install golang
  fi

fi

if ! [ -x "$(command -v node)" ]; then
  echo 'Error: node is not installed.' >&2
  echo 'Installing nodejs'
  # if you are using Ubuntu or any other distribution that has apt-get, run this
  if [ -x "$(command -v apt-get)" ]; then
    sudo apt-get install nodejs
  fi
  #  macos run this

  if [ -x "$(command -v brew)" ]; then
    brew install nodejs
  fi

  # if on windows run this

  if [ -x "$(command -v choco)" ]; then
    choco install nodejs
  fi

fi

if ! [ -x "$(command -v npm)" ]; then
  echo 'Error: npm is not installed.' >&2
  echo 'Installing npm'

  # if you are using Ubuntu or any other distribution that has apt-get, run this
  if [ -x "$(command -v apt-get)" ]; then
    sudo apt-get install npm
  fi

  #  macos run this

  if [ -x "$(command -v brew)" ]; then
    brew install npm
  fi

  # if on windows run this

  if [ -x "$(command -v choco)" ]; then
    choco install npm
  fi

fi

if ! [ -x "$(command -v nodemon)" ]; then
  echo 'Error: nodemon is not installed.' >&2
  echo 'Installing nodemon'
  # if you are using Ubuntu or any other distribution that has apt-get, run this
  if [ -x "$(command -v apt-get)" ]; then
    sudo npm install -g nodemon
  fi

  #  macos run this

  if [ -x "$(command -v brew)" ]; then
    brew install nodemon
  fi

  # if on windows run this

  if [ -x "$(command -v choco)" ]; then
    choco install nodemon
  fi

fi

# check if there is postgres installed and install it if not

if ! [ -x "$(command -v psql)" ]; then
  echo 'Error: postgres is not installed.' >&2
  echo 'Installing postgres'

  # if you are using Ubuntu or any other distribution that has apt-get, run this
  if [ -x "$(command -v apt-get)" ]; then
    sudo apt-get install postgresql postgresql-contrib
    sudo service postgresql start
  fi

  #  macos run this

  if [ -x "$(command -v brew)" ]; then
    brew install postgresql
    brew services start postgresql
  fi

  # if on windows run this

  if [ -x "$(command -v choco)" ]; then
    choco install postgresql
  fi

fi

nodemon --exec go run main.go --signal SIGTERM

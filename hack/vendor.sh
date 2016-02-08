#!/usr/bin/env bash
set -e

rm -rf vendor/
source 'hack/.vendor-helpers.sh'

clone git github.com/Sirupsen/logrus master
clone git github.com/codegangsta/cli master
clone git github.com/dancannon/gorethink master
clone git github.com/gorilla/context master
clone git github.com/gorilla/sessions master
clone git github.com/cenkalti/backoff master
clone git github.com/golang/protobuf master
clone git github.com/gorilla/mux master
clone git github.com/gorilla/securecookie master
clone git github.com/hailocab/go-hostpool master
clone git go.googlesource.com/crypto master
clone git gopkg.in/fatih/pool.v2 master

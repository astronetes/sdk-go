#!/bin/bash

mkdir -p tmp
#helm repo add bitnami https://charts.bitnami.com/bitnami
#helm pull bitnami/mysql --version 9.7.1
#mv mysql-9.7.1.tgz tmp/


helm repo add redis-operator https://spotahome.github.io/redis-operator
helm repo update
helm pull redis-operator --version 3.2.7
mv redis-operator-3.2.7 tmp/

helm install redis-operator redis-operator/redis-operator
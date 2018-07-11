# Syndesis Infrastructure Operator

An operator for installing and updating [Syndesis](https://github.com/syndesisio/syndesis).


## Building

```
dep ensure
operator-sdk build syndesis/syndesis-operator
```

## Running

```
minishift addons enable admin-user
minishift start
oc login -u system:admin
oc create -f deploy/syndesis-crd.yaml
eval $(minishift docker-env)
operator-sdk build syndesis/syndesis-operator
oc create -f deploy/syndesis-operator.yaml
```

#!/usr/bin/env bash

set -ex

export TARGET_NAMESPACE="${TARGET_NAMESPACE:-openshift-operators}"

for yml in operator.yaml  role_binding.yaml  role.yaml  service_account.yaml ; do
      curl -k https://raw.githubusercontent.com/AsherShoshan/cust0-pvc-operator/master/deploy/$yml | sed "s/TARGET_NAMESPACE/${TARGET_NAMESPACE}/g" | oc --namespace ${TARGET_NAMESPACE} delete -f -

done
  

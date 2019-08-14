# cust0-pvc-operator


Installation
------------
export TARGET_NAMESPACE=your-target-namespace     (default to openshift-operators)

curl -k https://raw.githubusercontent.com/AsherShoshan/cust0-pvc-operator/master/deploy.sh | bash

note: for Openshift 3.11, before installation, run this taint once on masters:
for node in $(oc get node -o name -l node-role.kubernetes.io/master); do oc adm taint node $node node-role.kubernetes.io/master=:NoSchedule ; done


Uinstall
--------
export TARGET_NAMESPACE=your-target-namespace     (default to openshift-operators)

curl -k https://raw.githubusercontent.com/AsherShoshan/cust0-pvc-operator/master/undeploy.sh | bash





# cust0-pvc-operator

Automatically annotate PVC's with node specific info, to support specific local-storage assignment requiremnts.


Install
-------
export TARGET_NAMESPACE=your-target-namespace     (default to openshift-operators)

curl -k https://raw.githubusercontent.com/AsherShoshan/cust0-pvc-operator/master/deploy.sh | bash

note: for Openshift 3.11, before installation, run this taint once on masters:

for node in $(oc get node -o name -l node-role.kubernetes.io/master); do oc adm taint node $node node-role.kubernetes.io/master=:NoSchedule ; done


Uninstall
---------
export TARGET_NAMESPACE=your-target-namespace     (default to openshift-operators)

curl -k https://raw.githubusercontent.com/AsherShoshan/cust0-pvc-operator/master/undeploy.sh | bash





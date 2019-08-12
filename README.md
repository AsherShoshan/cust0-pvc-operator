# cust0-pvc-operator


Installation
------------
1. clone repo
2. cd repo folder  
3. sed  's/namespace: TARGET_NAMESPACE/namespace: NAMESPACE/g'   ./deploy/*.yaml | kubectl create -f -
   "NAMESPACE" is traget namespace (!= default, !=kube-system)


note: for Openshift 3.11, run this taint on masters
for node in $(oc get node -o name -l node-role.kubernetes.io/master="true"); do oc adm taint node $node node-role.kubernetes.io/master=:NoSchedule --overwrite ; done


Uinstall
--------
2. cd repo folder
3. sed  's/namespace: TARGET_NAMESPACE/namespace: NAMESPACE/g'   ./deploy/*.yaml | kubectl delete -f -





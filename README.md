# cust0-pvc-operator


Installation
------------
1. clone repo
2. cd repo folder  
3. sed  's/namespace: TARGET_NAMESPACE/namespace: NAMESPACE/g'   ./deploy/*.yaml | kubectl create -f -
   "NAMESPACE" is traget namespace (!= default, !=kube-system)


Uinstall
--------
2. cd repo folder
3. sed  's/namespace: TARGET_NAMESPACE/namespace: NAMESPACE/g'   ./deploy/*.yaml | kubectl delete -f -


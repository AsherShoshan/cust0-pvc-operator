package pvc

import (
	"context"
	"math/rand"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("cust0-pvc-operator")

var hostPathStorageClass = "kubevirt-hostpath-provisioner"
var hostPathAnnotation = "kubevirt.io/provisionOnNode"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Ctl Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &Reconciler{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cust0-pvc-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	chck := func(obj runtime.Object) bool {
		pvc := obj.DeepCopyObject().(*corev1.PersistentVolumeClaim)
		return *pvc.Spec.StorageClassName == hostPathStorageClass
	}
	pred := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return chck(e.Object)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			//return chck(e.ObjectOld) || chck(e.ObjectNew)
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
	// Watch for changes to primary resource PVC
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForObject{}, pred)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCtl implements reconcile.Reconciler
var _ reconcile.Reconciler = &Reconciler{}

type Reconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Ctl object and makes changes based on the state read
// and what is in the Ctl.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PVC")

	// Fetch the PVC instance

	pvc := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(context.TODO(), request.NamespacedName, pvc)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Only with this storageClass
	if *pvc.Spec.StorageClassName != hostPathStorageClass {
		return reconcile.Result{}, nil
	}

	// If already there is an annotation - exit
	if value, found := pvc.Annotations[hostPathAnnotation]; found && value != "" {
		return reconcile.Result{}, nil
	}

	//build schedulable node name list
	schedNodelist, err := r.buildSchedNodeList()
	if err != nil {
		return reconcile.Result{}, err
	}
	//if no node available, requeue after 10Sec
	if len(schedNodelist) < 1 {
		return reconcile.Result{RequeueAfter: time.Second * 10}, nil
	}

	//handle if empty
	if pvc.Annotations == nil {
		pvc.Annotations = make(map[string]string)
	}
	//pick up random node from schedNodelist
	nodeName := schedNodelist[rand.Intn(len(schedNodelist))]
	// Annotate the PVC with the node
	reqLogger.Info("Annotate PVC " + hostPathAnnotation + "=" + nodeName)
	pvc.Annotations[hostPathAnnotation] = nodeName
	if err = r.client.Update(context.TODO(), pvc); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *Reconciler) buildSchedNodeList() ([]string, error) {

	var schedNodelist []string

	nodelist := &corev1.NodeList{}
	err := r.client.List(context.TODO(), &client.ListOptions{}, nodelist)
	if err == nil {
		for _, node := range nodelist.Items {
			sched := true
			for _, taint := range node.Spec.Taints {
				if taint.Effect == "NoSchedule" {
					sched = false
					break
				}
			}
			if sched {
				schedNodelist = append(schedNodelist, node.Name)
			}
		}
	}
	return schedNodelist, err
}


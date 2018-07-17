package cluster

import (
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// reconcileResource reconciles a runtime.Object to the state implied by the Custom Resource
// by building the resource on a deepcopied stub object and trying to create it and if that fails
// by getting the resource from the kubernetes cluster on a second deepcopied stub object and then synching
// it to the state implied by the Custom Resource and finally updating (patching) the resource inplace
// in the k8s cluster.
//
// It takes builder and syncer callbacks as arguments which implement the object manipulation of create and update,
// respectively using a type switches. Also a variadic indices parameter can be passed to content address (possibly nested)
// slices within the Custom Resource object through their indices.
//
// Passing indices helps flatten the resource building and synching while delegating looping entirely to the main
// reconcilation sequence. Althoug it seems complicated, this is a very transparent and maintainable code organization.
func reconcileResource(stub runtime.Object, cr *api.OdooCluster, build func(runtime.Object, *api.OdooCluster, ...int) error, sync func(runtime.Object, *api.OdooCluster, ...int) error, idx ...int) error {
	newObj := stub.DeepCopyObject()
	oldObj := stub.DeepCopyObject()

	if err := build(newObj, cr, idx...); err != nil {
		return err
	}
	if err := sdk.Create(newObj); err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("failed to create: %v", err)
		return err
	} else if errors.IsAlreadyExists(err) {
		if err := sdk.Get(oldObj); err != nil {
			return fmt.Errorf("failed to get: %v", err)
		}
		if err := sync(oldObj, cr, idx...); err != nil {
			return err
		}
		if err := sdk.Update(oldObj); err != nil {
			logrus.Errorf("failed to update: %v", err)
			return err
		}
	}
	return nil

}

var pgNamespaceMetaType metav1.TypeMeta = metav1.TypeMeta{
	Kind:       "PgNamespace",
	APIVersion: "odoo.k8s.io/v1alpha1",
}
var pvcMetaType metav1.TypeMeta = metav1.TypeMeta{
	Kind:       "PersistentVolumeClaim",
	APIVersion: "v1",
}
var cmMetaType metav1.TypeMeta = metav1.TypeMeta{
	Kind:       "ConfigMap",
	APIVersion: "v1",
}
var secMetaType metav1.TypeMeta = metav1.TypeMeta{
	Kind:       "Secret",
	APIVersion: "v1",
}
var deployMetaType metav1.TypeMeta = metav1.TypeMeta{
	Kind:       "Deployment",
	APIVersion: "apps/v1",
}
var svcMetaType metav1.TypeMeta = metav1.TypeMeta{
	Kind:       "Service",
	APIVersion: "v1",
}

// Reconcile reconciles the OdooCluster state to the spec specified by the Custom Resource
// by simulating initializers, in case default values are missing, and then reconciling resources
// to the spec given by the Custom Resource.
//
// Reconciliation sequence:
// - Reconcile PgNamespace
// - Loop until PgNamespace is ready
// - Reconcile 0, 1 or more PVCs  // TODO: Tear down of variadic resources
// - Reconcile ConfigMap
// - Reconcile Secret
// - Reconcile Deployment per Track & Tier  // TODO: Tear down of variadic resources
//
// Reconcile is idempotent. It is safe to retry on this functions.
// TODO: Though idempotent, missing tear down of variadic resources doesn't make this function truely converging
func Reconcile(c *api.OdooCluster) (err error) {
	c = c.DeepCopy()

	// Set defaults and loop through the event stack, in case any default applied
	changed := c.SetDefaults()
	if changed {
		return sdk.Update(c)
	}
	objectMeta := metav1.ObjectMeta{
		Name:      c.GetName(),
		Namespace: c.GetNamespace(),
		Labels:    selectorForOdooCluster(c.GetName()),
	}

	// TODO: Spin up PgCluster

	// Reconcile the PgNamespace for the cluster
	pgns := &api.PgNamespace{TypeMeta: pgNamespaceMetaType, ObjectMeta: objectMeta}
	logrus.Debugf("Reconciler (PgNamespace-Obj) ----- %+v", pgns)
	if err := reconcileResource(pgns, c, builder, syncer); err != nil {
		logrus.Errorf("Failed to reconcile %s (%s/%s): %v", pgns.Kind, c.Namespace, pgns.Name, err)
		return err
	}

	// Loop through the event stack until PgNamespace is finally ready.
	ready, err := isPgNamespaceReady(&c.Spec.PgSpec)
	if err != nil {
		return fmt.Errorf("failed to check if PgNamespace is ready: %v", err)
	}
	if !ready {
		logrus.Infof("Waiting for PgNamespace (%v) to become ready", c.Spec.PgSpec.User)
		return nil
	}

	// Reconcile PVC(s) for the cluster
	// TODO: Tear down of variadic resources
	for i, v := range c.Spec.Volumes {
		objectMetaPVC := metav1.ObjectMeta{
			Name:      volumeNameForOdoo(c, &v),
			Namespace: c.GetNamespace(),
			Labels:    selectorForOdooCluster(c.GetName()),
		}
		pvc := &v1.PersistentVolumeClaim{TypeMeta: pvcMetaType, ObjectMeta: objectMetaPVC}
		if err := reconcileResource(pvc, c, builder, syncer, i); err != nil {
			logrus.Errorf("Failed to reconcile %s (%s/%s): %v", pvc.Kind, c.Namespace, pvc.Name, err)
			return err
		}

	}

	// Reconcile the ConfigMap for the cluster
	cm := &v1.ConfigMap{TypeMeta: cmMetaType, ObjectMeta: objectMeta}
	logrus.Debugf("Reconciler (ConfigMap-Obj) ----- %+v", cm)
	if err := reconcileResource(cm, c, builder, syncer); err != nil {
		logrus.Errorf("Failed to reconcile %s (%s/%s): %v", cm.Kind, c.Namespace, cm.Name, err)
		return err
	}

	// Reconcile the Secret for the cluster
	se := &v1.Secret{TypeMeta: secMetaType, ObjectMeta: objectMeta}
	logrus.Debugf("Reconciler (Secret-Obj) ----- %+v", se)
	if err := reconcileResource(se, c, builder, syncer); err != nil {
		logrus.Errorf("Failed to reconcile %s (%s/%s): %v", se.Kind, c.Namespace, se.Name, err)
		return err
	}

	// Reconcile Deployment(s) for the cluster
	// One deployment per track *and* tier
	// TODO: Tear down of variadic resources
	for i, trs := range c.Spec.Tracks {
		for j, tis := range c.Spec.Tiers {
			selector := selectorForOdooCluster(c.GetName())
			objectMetaDeployment := metav1.ObjectMeta{
				Name:      getFullName(c, &trs, &tis),
				Namespace: c.GetNamespace(),
				Labels:    labelsWithTrackAndTier(selector, &trs, &tis),
			}
			d := &appsv1.Deployment{TypeMeta: deployMetaType, ObjectMeta: objectMetaDeployment}
			logrus.Debugf("Reconciler (Deployment-Obj) ----- %+v", d)
			if err := reconcileResource(d, c, builder, syncer, i, j); err != nil {
				logrus.Errorf("Failed to reconcile %s (%s/%s): %v", d.Kind, c.Namespace, d.Name, err)
				return err
			}

			switch tis.Name {
			case api.ServerTier, api.LongpollingTier:

				svc := &v1.Service{TypeMeta: svcMetaType, ObjectMeta: objectMetaDeployment}
				logrus.Debugf("Reconciler (Service-Obj) ----- %+v", svc)
				if err := reconcileResource(svc, c, builder, syncer, i, j); err != nil {
					logrus.Errorf("Failed to reconcile %s (%s/%s): %v", svc.Kind, c.Namespace, svc.Name, err)
					return err
				}
			}
		}
	}
	return nil
}

func isPgNamespaceReady(cr *api.PgNamespaceSpec) (bool, error) { return true, nil }

func ReconcileMigration(cr *api.ClusterMigration) (err error) { return nil }

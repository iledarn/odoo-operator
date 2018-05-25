package cluster

import (
	api "github.com/xoes/odoo-operator/pkg/apis/odoo/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// labelsForOdooCluster returns the labels for selecting the resources
// belonging to the given OdooCluster CR name.
func labelsForOdooCluster(name string) map[string]string {
	return map[string]string{"app": "odoo", "odoo_cluster": name}
}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(o metav1.Object, r metav1.OwnerReference) {
	o.SetOwnerReferences(append(o.GetOwnerReferences(), r))
}

// asOwner returns an OwnerReference set as the OdooCluster CR
func asOwner(oc *api.OdooCluster) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: oc.APIVersion,
		Kind:       oc.Kind,
		Name:       oc.Name,
		UID:        oc.UID,
		Controller: &trueVar,
	}

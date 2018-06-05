package cluster

import (
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

func crForPgNamespace(cr *api.OdooCluster) *api.PgNamespace {
	return &api.PgNamespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResource",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    selectorForOdooCluster(cr.GetName()),
		},
		Spec: cr.Spec.PgSpec,
	}
}

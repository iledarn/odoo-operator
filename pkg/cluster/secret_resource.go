package cluster

import (
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	"k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

func secretForOdooCluster(cr *api.OdooCluster) *v1.Secret {

	se := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels:    selectorForOdooCluster(cr.GetName()),
		},
		Data: map[string][]byte{
			"adminpwd": cr.Spec.AdminPassword,
			".pgpass":  cr.Spec.PgPassFile,
		},
	}
	addOwnerRefToObject(se, asOwner(cr))
	return se
}

package cluster

import (
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	"k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

func secretForOdooCluster(cr *api.OdooCluster) *v1.Secret { return nil }

package cluster

import (
	// "fmt"
	// "path/filepath"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	// apierrors "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

// func crForPgNamespace(cr *api.PgNamespace) *appsv1.Deployment         { return nil }
func configmapForOdooCluster(cr *api.OdooCluster) *v1.ConfigMap       { return nil }
func secretForOdooCluster(cr *api.OdooCluster) *v1.Secret             { return nil }
func deploymentForOdooCluster(cr *api.OdooCluster) *appsv1.Deployment { return nil }

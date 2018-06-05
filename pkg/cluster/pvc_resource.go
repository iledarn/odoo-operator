package cluster

import (
	// "fmt"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
	// appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

func getPVCsForOdooCluster(cr *api.OdooCluster) []*v1.PersistentVolumeClaim {
	var pvcs []*v1.PersistentVolumeClaim

	for s := range cr.Spec.PVCSpecs {

		pvc := &v1.PersistentVolumeClaim{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PersistentVolumeClaim",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetName(),
				Namespace: cr.GetNamespace(),
				Labels:    selectorForOdooCluster(cr.GetName()),
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
				Resources:        &s.Resources,
				VolumeName:       volumeNameForOdoo(cr, s),
				StorageClassName: &s.StorageClassName,
			},
		}
		addOwnerRefToObject(pvc, asOwner(cr))
		append(pvcs, pvc)
	}
	return pvcs
}

// volumeNameForOdoo is the volume name for the given odoo cluster.
func volumeNameForOdoo(cr *api.OdooCluster, s *api.PVCSpec) string {
	return cr.GetName() + fmt.Println(strings.ToLower(s.Name))
}

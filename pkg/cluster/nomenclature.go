package cluster

import (
	"fmt"
	"strings"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// selectorForOdooCluster returns the labels for selecting the resources
// belonging to the given OdooCluster CR name.
func selectorForOdooCluster(name string) map[string]string {
	return map[string]string{"app": "odoo", "odoo_cluster": name}
}

func getFullName(oc *api.OdooCluster, tr *api.TrackSpec, t *api.TierSpec) string {
	return fmt.Sprintf("%s-%s-%s", oc.GetName(), tr.Name, strings.ToLower(string(t.Name)))
}

func labelsWithTrackAndTier(selector map[string]string, tr *api.TrackSpec, t *api.TierSpec) map[string]string {
	labels := map[string]string{"track": tr.Name, "tier": string(t.Name)}
	for k, v := range selector {
		labels[k] = v
	}
	return labels
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
}

func getImageName(s *api.ImageSpec) string {
	return fmt.Sprintf("%s/%s:%s", s.Registry, s.Image, s.Tag)
}

// getVolumeName is the volume name for the given odoo cluster.
func getVolumeName(cr *api.OdooCluster, s string) string {
	return cr.GetName() + strings.ToLower(s)
}

func getVolumeNameFromConstant(cr *api.OdooCluster, s api.VolumeName) string {
	return getVolumeName(cr, fmt.Sprintf("%s", s))
}

func getMountPath(key string) string {
	return appMountPath + strings.ToLower(key) + "/"
}

func getMountPathFromConstant(key api.VolumeName) string {
	return getMountPath(fmt.Sprintf("%s", key))
}

func getSecretFile(key string) string {
	return appSecretsPath + strings.ToLower(key)
}

func getConfigFile(key string) string {
	return appConfigsPath + strings.ToLower(key)
}

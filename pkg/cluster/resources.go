package cluster

import (
	"path/filepath"
	"strings"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func builder(into runtime.Object, c *api.OdooCluster, i ...int) error {
	syncer(into, c, i...)
	switch o := into.(type) {

	case *api.PgNamespace:
		addOwnerRefToObject(o, asOwner(c))
		return nil

	case *v1.PersistentVolumeClaim:
		addOwnerRefToObject(o, asOwner(c))
		return nil

	case *v1.ConfigMap:
		addOwnerRefToObject(o, asOwner(c))
		return nil

	case *appsv1.Deployment:
		addOwnerRefToObject(o, asOwner(c))
		return nil

	case *v1.Service:
		addOwnerRefToObject(o, asOwner(c))
		return nil

	}

	return nil
}

func syncer(into runtime.Object, c *api.OdooCluster, i ...int) error {
	switch o := into.(type) {

	case *api.PgNamespace:
		o.Spec = c.Spec.PgSpec
		return nil

	case *v1.PersistentVolumeClaim:
		s := c.Spec.PVCSpecs[i[0]]
		o.Spec = v1.PersistentVolumeClaimSpec{
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources:        s.Resources,
			VolumeName:       volumeNameForOdoo(c, &s),
			StorageClassName: s.StorageClassName,
		}
		return nil

	case *v1.ConfigMap:
		var cfgDefaultData string
		var cfgCustomData string

		cfgDefaultData = newConfigWithDefaultParams(cfgDefaultData)
		o.Data = map[string]string{filepath.Base(odooDefaultConfig): cfgDefaultData}
		if len(c.Spec.ConfigMap) != 0 {
			cfgCustomData = c.Spec.ConfigMap
			o.Data[filepath.Base(odooCustomConfig)] = cfgCustomData
		}
		return nil

	case *appsv1.Deployment:
		volumes := []v1.Volume{
			{
				Name: configVolName,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: configMapNameForOdoo(c),
						},
					},
				},
			},
		}

		for _, s := range c.Spec.PVCSpecs {
			vol := v1.Volume{
				// kubernetes.io/pvc-protection
				Name: volumeNameForOdoo(c, &s),
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: volumeNameForOdoo(c, &s),
						ReadOnly:  false,
					},
				},
			}
			volumes = append(volumes, vol)

		}

		securityContext := &v1.PodSecurityContext{
			RunAsUser:    func(i int64) *int64 { return &i }(9001),
			RunAsNonRoot: func(b bool) *bool { return &b }(true),
			FSGroup:      func(i int64) *int64 { return &i }(9001),
		}

		trackSpec := c.Spec.Tracks[i[0]]
		tierSpec := c.Spec.Tiers[i[1]]

		podTempl := v1.PodTemplateSpec{
			ObjectMeta: o.ObjectMeta,
			Spec: v1.PodSpec{
				Containers: []v1.Container{odooContainer(c, &trackSpec, &tierSpec)},
				// Containers: []v1.Container{odooContainer(cr), odooMonitoringContainer(cr)},
				Volumes:         volumes,
				SecurityContext: securityContext,
			},
		}

		selector := selectorForOdooCluster(c.GetName())

		o.Spec = appsv1.DeploymentSpec{
			Replicas: &tierSpec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: podTempl,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
					MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
				},
			},
		}
		return nil

	case *v1.Service:
		selector := selectorForOdooCluster(c.GetName())
		var svcPorts []v1.ServicePort

		tierSpec := c.Spec.Tiers[i[1]]

		switch tierSpec.Name {
		case api.ServerTier:
			svcPorts = []v1.ServicePort{{
				Name:     clientPortName,
				Protocol: v1.ProtocolTCP,
				Port:     int32(clientPort),
			}}
		case api.LongpollingTier:
			svcPorts = []v1.ServicePort{{
				Name:     longpollingPortName,
				Protocol: v1.ProtocolTCP,
				Port:     int32(longpollingPort),
			}}
		}
		o.Spec = v1.ServiceSpec{
			Selector: selector,
			Ports:    svcPorts,
		}

		return nil

	}
	return nil
}

// volumeNameForOdoo is the volume name for the given odoo cluster.
func volumeNameForOdoo(cr *api.OdooCluster, s *api.PVCSpec) string {
	return cr.GetName() + strings.ToLower(string(s.Name))
}

// configMapNameForOdoo is the config volume name for the given odoo cluster.
func configMapNameForOdoo(cr *api.OdooCluster) string {
	return cr.GetName() + "config"
}

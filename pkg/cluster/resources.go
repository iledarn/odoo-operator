package cluster

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func builder(into runtime.Object, c *api.OdooCluster, i ...int) (string, error) {
	syncer(into, c, i...)
	switch o := into.(type) {

	case *api.PgNamespace:
		addOwnerRefToObject(o, asOwner(c))
		return o.GetName(), nil

	case *v1.PersistentVolumeClaim:
		addOwnerRefToObject(o, asOwner(c))
		return o.GetName(), nil

	case *v1.ConfigMap:
		addOwnerRefToObject(o, asOwner(c))
		return o.GetName(), nil

	case *appsv1.Deployment:
		addOwnerRefToObject(o, asOwner(c))
		return o.GetName(), nil

	case *v1.Service:
		addOwnerRefToObject(o, asOwner(c))
		return o.GetName(), nil

	case *v1.Secret:
		addOwnerRefToObject(o, asOwner(c))
		return o.GetName(), nil

	}

	return "", nil
}

func syncer(into runtime.Object, c *api.OdooCluster, i ...int) (bool, error) {
	changed := false
	switch o := into.(type) {

	case *api.PgNamespace:
		newSpec := c.Spec.PgSpec
		if !reflect.DeepEqual(o.Spec, newSpec) {
			changed = true
			o.Spec = newSpec
		}
		logrus.Debugf("Syncer (PgNamespace-Spec) +++++ %+v", o.Spec)
		return changed, nil

	case *v1.PersistentVolumeClaim:
		newSpec := c.Spec.Volumes[i[0]].Spec
		if !reflect.DeepEqual(o.Spec, newSpec) {
			changed = true
			o.Spec = newSpec
		}
		logrus.Debugf("Syncer (PVC-Spec) +++++ %+v", o.Spec)
		return changed, nil

	case *v1.ConfigMap:
		var cfgDefaultData string
		var cfgCustomData string

		cfgDefaultData = newConfigWithDefaultParams(cfgDefaultData)
		newSpec := map[string]string{appDefaultConfigKey: cfgDefaultData}
		if len(c.Spec.ConfigMap) != 0 {
			cfgCustomData = c.Spec.ConfigMap
			newSpec[appCustomConfigKey] = cfgCustomData
		}
		if !reflect.DeepEqual(o.Data, newSpec) {
			changed = true
			o.Data = newSpec
		}
		logrus.Debugf("Syncer (ConfigMap-Spec) +++++ %+v", o.Data)
		return changed, nil

	case *v1.Secret:
		var secPsql string
		var secAdmin string

		secPsqlBuf := newPsqlSecretWithParams(secPsql, &c.Spec.PgSpec)
		newSpec := map[string][]byte{appPsqlSecretKey: secPsqlBuf}
		secAdminBuf := newAdminSecretWithParams(secAdmin, &c.Spec.AdminPassword)
		newSpec[appAdminSecretKey] = secAdminBuf
		if !reflect.DeepEqual(o.Data, newSpec) {
			changed = true
			o.Data = newSpec
		}
		logrus.Debugf("Syncer (Secret-Spec) +++++ %+v", o.Data)
		return changed, nil

	case *appsv1.Deployment:
		volumes := []v1.Volume{
			{
				Name: getVolumeName(c, configVolName),
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							// We don't use suffixes on sinlgeton resources
							Name: c.GetName(),
						},
						DefaultMode: func(a int32) *int32 { return &a }(272), // octal 0420
					},
				},
			},
			{
				Name: getVolumeName(c, secretVolName),
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						// We don't use suffixes on sinlgeton resources
						SecretName:  c.GetName(),
						DefaultMode: func(a int32) *int32 { return &a }(256), // octal 0400
					},
				},
			},
		}

		for _, s := range c.Spec.Volumes {
			vol := v1.Volume{
				// kubernetes.io/pvc-protection
				Name: getVolumeName(c, s.Name),
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: getVolumeName(c, s.Name),
						ReadOnly:  false,
					},
				},
			}
			volumes = append(volumes, vol)

		}

		if !reflect.DeepEqual(o.Spec.Template.Spec.Volumes, volumes) {
			changed = true
			o.Spec.Template.Spec.Volumes = volumes
		}

		securityContext := &v1.PodSecurityContext{
			RunAsUser:    func(i int64) *int64 { return &i }(9001),
			RunAsNonRoot: func(b bool) *bool { return &b }(true),
			FSGroup:      func(i int64) *int64 { return &i }(9001),
		}

		if !reflect.DeepEqual(o.Spec.Template.Spec.SecurityContext, securityContext) {
			changed = true
			o.Spec.Template.Spec.SecurityContext = securityContext
		}

		trackSpec := c.Spec.Tracks[i[0]]
		tierSpec := c.Spec.Tiers[i[1]]
		newContainers := []v1.Container{odooContainer(c, &trackSpec, &tierSpec)}

		imagePullSecrets := []v1.LocalObjectReference{
			{
				Name: trackSpec.Image.Secret,
			},
		}

		if !reflect.DeepEqual(o.Spec.Template.Spec.ImagePullSecrets, imagePullSecrets) {
			changed = true
			o.Spec.Template.Spec.ImagePullSecrets = imagePullSecrets
		}

		if !reflect.DeepEqual(o.Spec.Template.Spec.Containers, newContainers) {
			// logrus.Errorf("OldContainers %+v", o.Spec.Template.Spec.Containers)
			// logrus.Errorf("NewContainers %+v", newContainers)
			// logrus.Error("NewContainers")
			changed = true
			o.Spec.Template.Spec.Containers = newContainers
		}
		o.Spec.Template.ObjectMeta = o.ObjectMeta

		selector := selectorForOdooCluster(c.GetName())

		if !reflect.DeepEqual(o.Spec.Selector, &metav1.LabelSelector{MatchLabels: selector}) {
			changed = true
			o.Spec.Selector = &metav1.LabelSelector{MatchLabels: selector}
		}
		if !reflect.DeepEqual(o.Spec.Replicas, &tierSpec.Replicas) {
			changed = true
			o.Spec.Replicas = &tierSpec.Replicas
		}

		newStrategy := appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
				MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
			},
		}
		if !reflect.DeepEqual(o.Spec.Strategy, newStrategy) {
			changed = true
			o.Spec.Strategy = newStrategy
		}
		logrus.Debugf("Syncer (Deployment-Spec) +++++ %+v", o.Spec)
		return changed, nil

	case *v1.Service:
		selector := selectorForOdooCluster(c.GetName())
		var svcPorts []v1.ServicePort

		tierSpec := c.Spec.Tiers[i[1]]

		switch tierSpec.Name {
		case api.ServerTier:
			svcPorts = []v1.ServicePort{{
				Name:       clientPortName,
				Protocol:   v1.ProtocolTCP,
				Port:       int32(clientPort),
				TargetPort: intstr.FromString(clientPortName),
			}}
		case api.LongpollingTier:
			svcPorts = []v1.ServicePort{{
				Name:       longpollingPortName,
				Protocol:   v1.ProtocolTCP,
				Port:       int32(longpollingPort),
				TargetPort: intstr.FromString(longpollingPortName),
			}}
		}

		if !reflect.DeepEqual(o.Spec.Selector, selector) {
			changed = true
			o.Spec.Selector = selector
		}
		if !reflect.DeepEqual(o.Spec.Ports, svcPorts) {
			changed = true
			o.Spec.Ports = svcPorts
		}
		logrus.Debugf("Syncer (Service-Spec) +++++ %+v", o.Spec)
		return changed, nil

	}
	return changed, nil
}

func newPsqlSecretWithParams(data string, p *api.PgNamespaceSpec) []byte {
	buf := bytes.NewBufferString(data)
	secret := fmt.Sprintf(odooPsqlSecretFmt,
		p.PgCluster.Host,
		p.PgCluster.Port,
		p.User,
		p.Password)
	buf.WriteString(secret)
	return []byte(buf.Bytes())
}

func newAdminSecretWithParams(data string, pwd *string) []byte {
	buf := bytes.NewBufferString(data)
	secret := fmt.Sprintf(odooAdminSecretFmt, &pwd)
	buf.WriteString(secret)
	return []byte(buf.Bytes())
}

func newConfigWithDefaultParams(data string) string {
	buf := bytes.NewBufferString(data)
	basicSection := fmt.Sprintf(odooBasicFmt,
		getMountPath(appPersistenceKey),
		odooWithoutDemo,
		odooServerWideModules,
		odooDbName,
		odooDbTemplate,
		odooListDb,
		odooDbFilter,
		odooPublisherWarrantyURL,
		getMountPath(appBackupKey),
		odooIntegratorWarrantyURL)
	buf.WriteString(basicSection)

	loggingSection := fmt.Sprintf(odooLoggingFmt,
		odooLogLevel)
	buf.WriteString(loggingSection)

	// multiprocSection := fmt.Sprintf(odooMultiprocFmt,
	// 	"")
	buf.WriteString(odooMultiprocFmt)

	SMTPSection := fmt.Sprintf(odooSMTPFmt,
		odooSMTPMail,
		odooSMTPServer,
		odooSMTPPort,
		odooSMTPSsl,
		odooSMTPUser,
		odooSMTPPassword)
	buf.WriteString(SMTPSection)

	return buf.String()
}

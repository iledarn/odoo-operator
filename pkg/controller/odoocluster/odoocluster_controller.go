/*
 * This file is part of the Odoo-Operator (R) project.
 * Copyright (c) 2018-2018 XOE Corp. SAS
 * Authors: David Arnold, et al.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * ALTERNATIVE LICENCING OPTION
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving the Odoo-Operator software without
 * disclosing the source code of your own applications. These activities
 * include: Offering paid services to a customer as an ASP, shipping Odoo-
 * Operator with a closed source product.
 *
 */

package odoocluster

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new OdooCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this cluster.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOdooCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("odoocluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to OdooCluster
	err = c.Watch(&source.Kind{Type: &clusterv1beta1.OdooCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in DBNamespace owned by OdooCluster
	err = c.Watch(&source.Kind{Type: &clusterv1beta1.DBNamespace{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterv1beta1.OdooCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes in PersistentVolumeClaims owned by OdooCluster
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterv1beta1.OdooCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes in ConfigMaps owned by OdooCluster
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterv1beta1.OdooCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes in Secrets owned by OdooCluster
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterv1beta1.OdooCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes in Deployment owned by OdooCluster
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterv1beta1.OdooCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes in Services owned by OdooCluster
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &clusterv1beta1.OdooCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileOdooCluster{}

const (

	// Ports and Port Names
	clientPortName      string = "client-port"
	clientPort                 = 8069
	longpollingPortName        = "lp-port"
	longpollingPort            = 8072

	// Environment Variables
	envPGHOST       = "PGHOST"
	envPGUSER       = "PGUSER"
	envPGPASSFILE   = "PGPASSFILE"
	envODOORC       = "ODOO_RC"
	envODOOPASSFILE = "ODOO_PASSFILE"

	// App paths
	appMountPath   = "/mnt/odoo/"
	appConfigsPath = "/run/configs/odoo/"
	appSecretsPath = "/run/secrets/odoo/"

	// ConfigMaps, Secrets & Volumes Keys
	appPsqlSecretKey  = "pgpass"
	appAdminSecretKey = "adminpwd"

	// Basic Config
	defaultServerTierMaxConn      = "16"
	defaultLongpollingTierMaxConn = "16"
	defaultWithoutDemo            = "True"
	defaultServerWideModules      = "base,web"
	defaultDbName                 = "False"
	defaultDbTemplate             = "template0"
	defaultListDb                 = "False"
	defaultDbFilter               = "^%h$"
	defaultPublisherWarrantyURL   = "http://services.openerp.com/publisher-warranty/"

	// Log Config
	defaultLogLevel = ":INFO"
)

// ReconcileOdooCluster reconciles a OdooCluster object
type ReconcileOdooCluster struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a OdooCluster object and makes changes based on the state read
// and what is in the OdooCluster.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistenvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.odoo.io,resources=dbnamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.odoo.io,resources=odooclusters,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileOdooCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Fetch the OdooCluster instance
	instance := &clusterv1beta1.OdooCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	/* --------------------------------- */
	/* --- Reconcile the DBNamespace --- */
	/* --------------------------------- */

	// Define the object name/namespace
	objDBNamespace := &clusterv1beta1.DBNamespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-dbnamespace",
			Namespace: instance.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(context.TODO(), r, objDBNamespace, func(existing runtime.Object) error {
		// mutate here the state of existing object to the desired state
		// note that at this point _objDBNamespace_ and _existing_ point to the same struct
		out := existing.(*clusterv1beta1.DBNamespace)

		out.Spec = instance.Spec.DBSpec // Desired state

		if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objDBNamespace.Namespace, objDBNamespace.Name, result, err)
		return reconcile.Result{}, err
	}
	log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objDBNamespace.Namespace, objDBNamespace.Name, result)

	/* ----------------------------- */
	/* --- Reconcile the Volumes --- */
	/* ----------------------------- */

	for _, volume := range instance.Spec.Volumes {

		// Define the object name/namespace
		objPVC := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(fmt.Sprintf("%s-%s-volume", instance.Name, volume.Name)),
				Namespace: instance.Namespace,
			},
		}

		result, err := controllerutil.CreateOrUpdate(context.TODO(), r, objPVC, func(existing runtime.Object) error {
			// mutate here the state of existing object to the desired state
			// note that at this point _objPVC_ and _existing_ point to the same struct
			out := existing.(*corev1.PersistentVolumeClaim)

			// PVS spec is immutable after creation
			if out.ObjectMeta.CreationTimestamp.IsZero() {
				out.Spec = volume.Spec // Desired state

				if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objPVC.Namespace, objPVC.Name, result, err)
			return reconcile.Result{}, err
		}
		log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objPVC.Namespace, objPVC.Name, result)
	}

	/* -------------------------------- */
	/* --- Reconcile the ConfigMaps --- */
	/* -------------------------------- */

	for _, track := range instance.Spec.Tracks {
		// Define the object name/namespace
		objConfigMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(fmt.Sprintf("%s-%s-config", instance.Name, track.Name)),
				Namespace: instance.Namespace,
			},
		}

		result, err := controllerutil.CreateOrUpdate(context.TODO(), r, objConfigMap, func(existing runtime.Object) error {
			// mutate here the state of existing object to the desired state
			// note that at this point _objConfigMap_ and _existing_ point to the same struct
			out := existing.(*corev1.ConfigMap)

			cfgDefaultData := newDefaultConfig()
			cfgOptionsData := newOptionsConfig(instance.Spec.Config, track.Config)
			cfgIntegratorData := newIntegratorConfig(instance.Spec.IntegratorConfig, track.IntegratorConfig)
			cfgCustomData := newCustomConfig(instance.Spec.CustomConfig, track.CustomConfig)

			out.Data = map[string]string{
				"01-DEFAULT":    cfgDefaultData,
				"02-options":    cfgOptionsData,
				"03-integrator": cfgIntegratorData,
				"04-custom":     cfgCustomData,
			} // Desired state

			if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objConfigMap.Namespace, objConfigMap.Name, result, err)
			return reconcile.Result{}, err
		}
		log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objConfigMap.Namespace, objConfigMap.Name, result)
	}

	/* ---------------------------- */
	/* --- Reconcile the Secret --- */
	/* ---------------------------- */

	// Define the object name/namespace
	objSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(fmt.Sprintf("%s-secret", instance.Name)),
			Namespace: instance.Namespace,
		},
	}

	result, err = controllerutil.CreateOrUpdate(context.TODO(), r, objSecret, func(existing runtime.Object) error {
		// mutate here the state of existing object to the desired state
		// note that at this point _objDBNamespace_ and _existing_ point to the same struct
		out := existing.(*corev1.Secret)

		secPsqlBuf := newPsqlSecretWithParams(&instance.Spec.DBSpec)
		secAdminBuf := newAdminSecretWithParams(instance.Spec.AdminPassword)

		out.Data = map[string][]byte{
			appPsqlSecretKey:  secPsqlBuf,
			appAdminSecretKey: secAdminBuf,
		} // Desired state

		if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objSecret.Namespace, objSecret.Name, result, err)
		return reconcile.Result{}, err
	}
	log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objSecret.Namespace, objSecret.Name, result)

	/* --------------------------------- */
	/* --- Reconcile the Deployments --- */
	/* --------------------------------- */

	for _, track := range instance.Spec.Tracks {
		for _, tier := range instance.Spec.Tiers {
			// Define the object name/namespace
			objDeployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.ToLower(fmt.Sprintf("%s-%s-%s-deployment", instance.Name, track.Name, tier.Name)),
					Namespace: instance.Namespace,
				},
			}

			result, err := controllerutil.CreateOrUpdate(context.TODO(), r, objDeployment, func(existing runtime.Object) error {
				// mutate here the state of existing object to the desired state
				// note that at this point _objConfigMap_ and _existing_ point to the same struct
				out := existing.(*appsv1.Deployment)

				setDeploymentSpec(out, instance, &track, &tier)

				if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objDeployment.Namespace, objDeployment.Name, result, err)
				return reconcile.Result{}, err
			}
			log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objDeployment.Namespace, objDeployment.Name, result)
		}
	}

	/* ------------------------------ */
	/* --- Reconcile the Services --- */
	/* ------------------------------ */

	for _, track := range instance.Spec.Tracks {
		for _, tier := range instance.Spec.Tiers {
			// Define the object name/namespace
			objService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.ToLower(fmt.Sprintf("%s-%s-%s-service", instance.Name, track.Name, tier.Name)),
					Namespace: instance.Namespace,
				},
			}
			var result controllerutil.OperationResult
			switch tier.Name {
			case clusterv1beta1.ServerTier:
				result, err = controllerutil.CreateOrUpdate(context.TODO(), r, objService, func(existing runtime.Object) error {
					// mutate here the state of existing object to the desired state
					// note that at this point _objConfigMap_ and _existing_ point to the same struct
					out := existing.(*corev1.Service)

					svcPorts := []corev1.ServicePort{{
						Name:       clientPortName,
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(clientPort),
						TargetPort: intstr.FromString(clientPortName),
					}}

					out.Spec.Ports = svcPorts // Desired state
					out.Spec.Selector = map[string]string{
						"tier": fmt.Sprintf("%s", clusterv1beta1.ServerTier),
					} // Desired state

					if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
						return err
					}

					return nil
				})

				if err != nil {
					log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objService.Namespace, objService.Name, result, err)
					return reconcile.Result{}, err
				}
				log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objService.Namespace, objService.Name, result)
			case clusterv1beta1.LongpollingTier:
				result, err = controllerutil.CreateOrUpdate(context.TODO(), r, objService, func(existing runtime.Object) error {
					// mutate here the state of existing object to the desired state
					// note that at this point _objConfigMap_ and _existing_ point to the same struct
					out := existing.(*corev1.Service)

					svcPorts := []corev1.ServicePort{{
						Name:       longpollingPortName,
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(longpollingPort),
						TargetPort: intstr.FromString(longpollingPortName),
					}}

					out.Spec.Ports = svcPorts // Desired state
					out.Spec.Selector = map[string]string{
						"tier": fmt.Sprintf("%s", clusterv1beta1.LongpollingTier),
					} // Desired state

					if err := controllerutil.SetControllerReference(instance, out, r.scheme); err != nil {
						return err
					}

					return nil
				})

				if err != nil {
					log.Printf("%s/%s Operation: %s. (Controller: OdooCluster) Error: %s\n", objService.Namespace, objService.Name, result, err)
					return reconcile.Result{}, err
				}
				log.Printf("%s/%s reconciled. Operation: %s. (Controller: OdooCluster)\n", objService.Namespace, objService.Name, result)
			}
		}
	}

	return reconcile.Result{}, nil
}

func setDeploymentSpec(
	out *appsv1.Deployment, instance *clusterv1beta1.OdooCluster,
	track *clusterv1beta1.TrackSpec, tier *clusterv1beta1.TierSpec) {

	// First, setup the Deployment

	out.Spec = appsv1.DeploymentSpec{
		// Set Deployment Strategy
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
				MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
			},
		},
		// Set Deployment Replicas
		Replicas: &tier.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"cluster": strings.ToLower(instance.Name),
				"tier":    strings.ToLower(fmt.Sprintf("%s", tier.Name)),
				"track":   strings.ToLower(track.Name),
			},
		},
	}

	// Second, setup the Deployment Template

	// Set Template Labels (to match deployment labels)
	out.Spec.Template.Labels = map[string]string{
		"cluster": strings.ToLower(instance.Name),
		"tier":    strings.ToLower(fmt.Sprintf("%s", tier.Name)),
		"track":   strings.ToLower(track.Name),
	}

	// Set Template Spec
	out.Spec.Template.Spec = corev1.PodSpec{
		// Set Template Volumes
		Volumes: []corev1.Volume{
			{
				Name: strings.ToLower(fmt.Sprintf("%s-config", instance.Name)),
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: strings.ToLower(fmt.Sprintf("%s-%s-config", instance.Name, track.Name)),
						},
						DefaultMode: func(a int32) *int32 { return &a }(272), // octal 0420
					},
				},
			},
			{
				Name: strings.ToLower(fmt.Sprintf("%s-secret", instance.Name)),
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName:  strings.ToLower(fmt.Sprintf("%s-secret", instance.Name)),
						DefaultMode: func(a int32) *int32 { return &a }(256), // octal 0400
					},
				},
			},
		},
		// Set Template SecurityContext
		SecurityContext: &corev1.PodSecurityContext{
			RunAsUser:          func(i int64) *int64 { return &i }(9001),
			RunAsNonRoot:       func(b bool) *bool { return &b }(true),
			FSGroup:            func(i int64) *int64 { return &i }(9001),
			SupplementalGroups: []int64{2000}, // Host volume group, with 770 access.
		},
		// Set Template ImagePullSecrets
		ImagePullSecrets: []corev1.LocalObjectReference{
			{
				Name: track.Image.Secret,
			},
		},
	}
	// Set containers
	out.Spec.Template.Spec.Containers = getContainerSpec(instance, track, tier)

	// Set additional volumes dynamically
	for _, s := range instance.Spec.Volumes {
		vol := corev1.Volume{
			// kubernetes.io/pvc-protection
			Name: strings.ToLower(fmt.Sprintf("%s-%s-volume", instance.Name, s.Name)),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: strings.ToLower(fmt.Sprintf("%s-%s-volume", instance.Name, s.Name)),
					ReadOnly:  false,
				},
			},
		}
		out.Spec.Template.Spec.Volumes = append(out.Spec.Template.Spec.Volumes, vol)

	}
}

func getContainerSpec(instance *clusterv1beta1.OdooCluster, track *clusterv1beta1.TrackSpec, tier *clusterv1beta1.TierSpec) []corev1.Container {

	containers := []corev1.Container{}

	container := corev1.Container{
		Name:  strings.ToLower(fmt.Sprintf("%s-%s-%s", instance.Name, track.Name, tier.Name)),
		Image: strings.ToLower(fmt.Sprintf("%s/%s:%s", track.Image.Registry, track.Image.Image, track.Image.Tag)),
		Lifecycle: &corev1.Lifecycle{
			PostStart: &corev1.Handler{
				Exec: &corev1.ExecAction{
					// TODO: until proper fix of https://github.com/kubernetes/kubernetes/issues/2630
					Command: []string{"sh", "-c",
						"mkdir /run/secrets/patched && cat " + appSecretsPath + appPsqlSecretKey +
							" > /run/secrets/patched/" + appPsqlSecretKey +
							" && chmod 400 /run/secrets/patched/" + appPsqlSecretKey},
				},
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  envPGHOST,
				Value: instance.Spec.DBSpec.Host,
			},
			{
				Name:  envPGUSER,
				Value: instance.Spec.DBSpec.User,
			},
			{
				Name: envPGPASSFILE,
				// TODO: until proper fix of https://github.com/kubernetes/kubernetes/issues/2630
				Value: "/run/secrets/patched/" + appPsqlSecretKey,
			},
			{
				Name:  envODOORC,
				Value: appConfigsPath,
			},
			{
				Name:  envODOOPASSFILE,
				Value: appSecretsPath + appAdminSecretKey,
			},
		},
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		ImagePullPolicy:          corev1.PullAlways,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      instance.Name + "-config",
				MountPath: filepath.Dir(appConfigsPath),
				ReadOnly:  true,
			},
			{
				Name:      instance.Name + "-secret",
				MountPath: filepath.Dir(appSecretsPath),
				ReadOnly:  true,
			},
		},
	}

	// Set additional VolumeMounts dynamically
	volumes := []corev1.VolumeMount{}
	for _, s := range instance.Spec.Volumes {
		volumes = append(volumes, corev1.VolumeMount{
			Name:      strings.ToLower(fmt.Sprintf("%s-%s", instance.Name, s.Name)),
			MountPath: strings.ToLower(fmt.Sprintf("%s%s", appMountPath, s.Name)) + "/",
		})
	}
	switch tier.Name {
	case clusterv1beta1.ServerTier:
		setServerTierContainerSpec(&container, tier)
	case clusterv1beta1.LongpollingTier:
		setLonpgollingTierContainerSpec(&container, tier)
	case clusterv1beta1.CronTier:
		setCronTierContainerSpec(&container, tier)
	}
	containers = append(containers, container)
	return containers

}

func setServerTierContainerSpec(container *corev1.Container, tier *clusterv1beta1.TierSpec) {
	var maxConn string
	if tier.DBConn != nil {
		maxConn = fmt.Sprintf("%v", tier.DBConn)
	} else {
		maxConn = defaultServerTierMaxConn
	}
	container.Args = []string{"--config", appConfigsPath, "--db_maxconn=" + maxConn, "--workers=0", "--max-cron-threads=0"}
	container.Ports = []corev1.ContainerPort{{
		Name:          clientPortName,
		ContainerPort: int32(clientPort),
		Protocol:      corev1.ProtocolTCP,
	}}

	container.LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"curl",
					"--connect-timeout", "5",
					"--max-time", "10",
					"-k", "-s",
					fmt.Sprintf("http://localhost:%d", clientPort),
				},
			},
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      10,
		PeriodSeconds:       60,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}
	container.ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Port:   intstr.FromInt(clientPort),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      10,
		PeriodSeconds:       10,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}
}

func setLonpgollingTierContainerSpec(container *corev1.Container, tier *clusterv1beta1.TierSpec) {

	var maxConn string
	if tier.DBConn != nil {
		maxConn = fmt.Sprintf("%v", tier.DBConn)
	} else {
		maxConn = defaultLongpollingTierMaxConn
	}
	container.Args = []string{"gevent", "--config", appConfigsPath, "--db_maxconn=" + maxConn}
	container.Ports = []corev1.ContainerPort{}

	container.LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"curl",
					"--connect-timeout", "5",
					"--max-time", "10",
					"-k", "-s",
					fmt.Sprintf("http://localhost:%d", longpollingPort),
				},
			},
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      10,
		PeriodSeconds:       60,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}
	container.ReadinessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Port:   intstr.FromInt(longpollingPort),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      10,
		PeriodSeconds:       10,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}

}

func setCronTierContainerSpec(container *corev1.Container, tier *clusterv1beta1.TierSpec) {
	container.Args = []string{"--config", appConfigsPath, "--db_maxconn=1", "--workers=0", "--max-cron-threads=1", "--no-xmlrpc"}
	container.Ports = []corev1.ContainerPort{{
		Name:          longpollingPortName,
		ContainerPort: int32(longpollingPort),
		Protocol:      corev1.ProtocolTCP,
	}}
}

// ---- Helper functions ---- //

func newPsqlSecretWithParams(p *clusterv1beta1.DBNamespaceSpec) []byte {
	var data string
	buf := bytes.NewBufferString(data)
	secret := fmt.Sprintf(odooPsqlSecretFmt,
		p.Host,
		p.Port,
		p.User,
		p.Password)
	buf.WriteString(secret)
	return []byte(buf.Bytes())
}

func newAdminSecretWithParams(pwd string) []byte {
	var data string
	buf := bytes.NewBufferString(data)
	secret := fmt.Sprintf(odooAdminSecretFmt, pwd)
	buf.WriteString(secret)
	return []byte(buf.Bytes())
}

func newDefaultConfig() string {
	var s string
	buf := bytes.NewBufferString(s)
	section := fmt.Sprintf(odooDefaultSection,
		defaultWithoutDemo,
		defaultServerWideModules,
		defaultDbName,
		defaultDbTemplate,
		defaultListDb,
		defaultDbFilter,
		defaultPublisherWarrantyURL,
		defaultLogLevel)
	buf.WriteString(section)
	return buf.String()
}

func newOptionsConfig(clusterOverrides *string, trackOverrides *string) string {
	var s string
	buf := bytes.NewBufferString(s)

	var cO string
	var tO string

	if clusterOverrides != nil {
		cO = *clusterOverrides
	} else {
		cO = ""
	}

	if trackOverrides != nil {
		tO = *trackOverrides
	} else {
		tO = ""
	}

	section := fmt.Sprintf(odooOptionsSection,
		strings.ToLower(fmt.Sprintf("%s%s", appMountPath, clusterv1beta1.PVCNameData))+"/",
		strings.ToLower(fmt.Sprintf("%s%s", appMountPath, clusterv1beta1.PVCNameBackup))+"/",
		cO, tO)
	buf.WriteString(section)
	return buf.String()
}

func newIntegratorConfig(clusterOverrides *string, trackOverrides *string) string {
	var s string
	buf := bytes.NewBufferString(s)

	var cO string
	var tO string

	if clusterOverrides != nil {
		cO = *clusterOverrides
	} else {
		cO = ""
	}

	if trackOverrides != nil {
		tO = *trackOverrides
	} else {
		tO = ""
	}

	section := fmt.Sprintf(odooIntegratorSection,
		cO, tO)
	buf.WriteString(section)
	return buf.String()
}

func newCustomConfig(clusterOverrides *string, trackOverrides *string) string {
	var s string
	buf := bytes.NewBufferString(s)

	var cO string
	var tO string

	if clusterOverrides != nil {
		cO = *clusterOverrides
	} else {
		cO = ""
	}

	if trackOverrides != nil {
		tO = *trackOverrides
	} else {
		tO = ""
	}

	section := fmt.Sprintf(odooCustomSection,
		cO, tO)
	buf.WriteString(section)
	return buf.String()
}

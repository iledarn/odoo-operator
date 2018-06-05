package cluster

import (
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func servicesForOdooTrack(cr *api.OdooCluster, trck *api.TrackSpec) []*v1.Service {
	selector := selectorForOdooCluster(cr.GetName())

	var services []*v1.Service
	var svcPorts []v1.ServicePort

	for _, tierSpec := range cr.Spec.Tiers {
		// Construct the ServicePorts according to the tier or continue
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
		default:
			continue
		}

		objectMeta := metav1.ObjectMeta{
			Name:      getFullName(cr, trck, &tierSpec),
			Namespace: cr.GetNamespace(),
			Labels:    labelsWithTrackAndTier(selector, trck, &tierSpec),
		}
		svc := &v1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: objectMeta,
			Spec: v1.ServiceSpec{
				Selector: selector,
				Ports:    svcPorts,
			},
		}
		addOwnerRefToObject(svc, asOwner(cr))
		services = append(services, svc)
	}
	return services

}

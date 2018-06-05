package cluster

import (
	"fmt"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
)

// Reconcile reconciles the OdooCluster state to the spec specified by cr
// by reconciling the "PG Namespace" (pg user) and then deply the OdooCluster
// and finally update the OdooCluster if needed.
//
// Reconcile is idempotent. It is safe to retry on this functions.
func Reconcile(cr *api.OdooCluster) (err error) {
	cr = cr.DeepCopy()

	// TODO: Spin up PgCluster

	// Create or update PgNamespace for OdooCluster
	err = sdk.Create(crForPgNamespace(cr))
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create odoo cluster PgNamespace: %v", err)
		return err
	} else if errors.IsAlreadyExists(err) {
		err := sdk.Update(crForPgNamespace(cr))
		if err != nil {
			logrus.Errorf("Failed to update odoo cluster PgNamespace: %v", err)
			return err
		}
	}

	// Check if PgNamespace is ready.
	// If not, we need to wait until it is provisioned before proceeding;
	// Hence, we return from here and let the Watch triggers the handler again.
	ready, err := isPgNamespaceReady(cr.Spec.PgSpec)
	if err != nil {
		return fmt.Errorf("failed to check if PgNamespace is ready: %v", err)
	}
	if !ready {
		logrus.Infof("Waiting for PgNamespace (%v) to become ready", cr.Spec.PgSpec.User)
		return nil
	}

	// Create or update persistence PVC for OdooCluster
	for _, pvc := range getPVCsForOdooCluster(cr) {
		err = sdk.Create(pvc)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create odoo cluster persistence PVC: %v", err)
			return err
		} else if errors.IsAlreadyExists(err) {
			err := sdk.Update(pvc)
			if err != nil {
				logrus.Errorf("Failed to update odoo cluster persistence PVC: %v", err)
				return err
			}
		}
	}

	// Create or update ConfigMap for OdooCluster
	err = sdk.Create(configmapForOdooCluster(cr))
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create odoo cluster ConfigMap: %v", err)
		return err
	} else if errors.IsAlreadyExists(err) {
		err := sdk.Update(configmapForOdooCluster(cr))
		if err != nil {
			logrus.Errorf("Failed to update odoo cluster ConfigMap: %v", err)
			return err
		}
	}

	// Create or update Secret for OdooCluster
	err = sdk.Create(secretForOdooCluster(cr))
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Failed to create odoo cluster Secret: %v", err)
		return err
	} else if errors.IsAlreadyExists(err) {
		err := sdk.Update(secretForOdooCluster(cr))
		if err != nil {
			logrus.Errorf("Failed to update odoo cluster Secret: %v", err)
			return err
		}
	}

	// TODO Delete deployments
	// Create or update OdooCluster for every Track
	for _, t := range cr.Spec.Tracks {
		for _, d := range deploymentsForOdooTrack(cr, t) {
			err = sdk.Create(d)
			if err != nil && !errors.IsAlreadyExists(err) {
				logrus.Errorf("Failed to create odoo cluster : %v", err)
				return err
			} else if errors.IsAlreadyExists(err) {
				err = sdk.Update(d)
				if err != nil {
					logrus.Errorf("Failed to update odoo cluster : %v", err)
					return err
				}
			}

		}
		for _, s := range servicesForOdooTrack(cr, t) {
			err = sdk.Create(s)
			if err != nil && !errors.IsAlreadyExists(err) {
				logrus.Errorf("Failed to create odoo cluster : %v", err)
				return err
			} else if errors.IsAlreadyExists(err) {
				err = sdk.Update(s)
				if err != nil {
					logrus.Errorf("Failed to update odoo cluster : %v", err)
					return err
				}
			}

		}
	}
	return nil
}

func isPgNamespaceReady(cr *api.PgNamespaceSpec) (bool, error) { return true, nil }

func ReconcileMigration(cr *api.ClusterMigration) (err error) { return nil }

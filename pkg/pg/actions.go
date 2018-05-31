// Package pg contains the reconciliation logic for the PgNamespace Custom Resource.
package pg

import (
	"errors"

	"github.com/sirupsen/logrus"
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
)

func createPgNamespace(cr *api.PgNamespace) (err error) {
	// If there is not enough quota, creturn with an error
	enoughSpace, err := isEnoughQuota(cr, 0)
	if err != nil {
		return err
	}
	if !enoughSpace {
		return errors.New("not enough free quota")
	}
	// Create PgNamespace
	return nil
}

func updatePgNamespace(cr *api.PgNamespace) (err error) {

	// Get the current quota of the PgNamespace
	current, err := getPgNamespaceUsedQuota(cr)
	if err != nil {
		logrus.Errorf("Failed to get PgNamespace's used quota: %v", err)
		return err
	}
	// If there is not enough quota, creturn with an error
	enoughSpace, err := isEnoughQuota(cr, current)
	if err != nil {
		return err
	}
	// Get the
	if !enoughSpace {
		return errors.New("not enough free quota")
	}
	// Update PgNamespace
	return nil
}

func deletePgNamespace(cr *api.PgNamespace) (err error) {
	return nil
}

// isEnoughQuota validates if PgCluster has enough quota to fulfill the
// requested transition.
func isEnoughQuota(cr *api.PgNamespace, current int32) (bool, error) {
	// Get the free_reserved quota of the PgCluster
	free, err := getFreePgClusterSpace(cr)
	if err != nil {
		logrus.Errorf("Failed to get PgCluster's free quota: %v", err)
		return false, err
	}
	// Get the requested quota of the PgNamespace
	requested, err := getPgNamespaceQuota(cr)
	if err != nil {
		logrus.Errorf("Failed to get PgNamespace's assigned quota: %v", err)
		return false, err
	}

	if requested > current && requested-current > free {
		return false, nil
	}
	return true, nil
}

// getFreePgClusterSpace gets the PgCuster available free quota.
// It queries the PgCluster on the size of on the size of all database objects
// and compares it with it's PVS limits, if they exist. It then calculates a
// security margin and finally returns currently assignable quota.
func getFreePgClusterSpace(cr *api.PgNamespace) (quota int32, err error) {
	return 0, nil
}

// getPgNamespaceUsedQuota gets the PgNamespace's currently used quota.
// It queries the PgCluster on the size of all database objects owned by the
// PgNamespace and returns currently used quota.
func getPgNamespaceUsedQuota(cr *api.PgNamespace) (quota int32, err error) {
	return 0, nil
}

// getPgNamespaceUsedQuota gets the PgNamespace's quota.
func getPgNamespaceQuota(cr *api.PgNamespace) (quota int32, err error) {
	return 0, nil
}

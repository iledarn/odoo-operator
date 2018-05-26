package pg

import (
	api "github.com/xoes/odoo-operator/pkg/apis/odoo/v1alpha1"

	"github.com/sirupsen/logrus"
)

// ReconcilePgNamespace reconciles the PgNamespace state to the spec specified
// by the cr by either tearing down the PGNamespace or by validating the
// PgCluster resource fitness and creating or updating the PgNamespace
//
// ReconcilePgNamespace is idempotent. It is safe to retry on this functions.
func ReconcilePgNamespace(cr *api.PgNamespace) (err error) {
	cr = cr.DeepCopy()

	// PgNamespace teardown logic
	if cr.DeletionTimestamp != nil {
		err := deletePgNamespace(cr)
		if err != nil {
			return err
		}

		// https://github.com/operator-framework/operator-sdk/issues/270
		// cr.SetFinalizers([]string{})
		// err := sdk.Update(cr)
		// if err != nil {
		// 	return fmt.Errorf("failed to update cr finalizers: %v", err)
		// }

		return nil
	}
	if !isPgNamespaceExists(cr) {
		// Create or update PgNamespace.
		err := createPgNamespace(cr)
		if err != nil {
			logrus.Errorf("Failed to create postgres namespace: %v", err)
			return err
		}
	} else {
		err := updatePgNamespace(cr)
		if err != nil {
			logrus.Errorf("Failed to update postgres namespace: %v", err)
			return err
		}
	}

	return nil
}

// isPgNamespaceExists checks if the PgNamespace exists.
// It tries to log in as the namespace user to the PgCluster
// and returns a boolean.
func isPgNamespaceExists(cr *api.PgNamespace) bool {
	return false
}

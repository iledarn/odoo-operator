package pg

import (
	"database/sql"
	"fmt"

	// load postgres driver
	_ "github.com/lib/pq"

	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"

	"github.com/sirupsen/logrus"
)

const (
	DBMgtName = "postgres"
)

// ReconcilePgNamespace reconciles the PgNamespace state to the spec specified
// by the cr by either tearing down the PGNamespace or by validating the
// PgCluster resource fitness and creating or updating the PgNamespace
//
// ReconcilePgNamespace is idempotent. It is safe to retry on this functions.
func ReconcilePgNamespace(cr *api.PgNamespace) (err error) {
	cr = cr.DeepCopy()

	exists, err := IsPgNamespaceExists(&cr.Spec)
	if err != nil {
		return fmt.Errorf("failed to check if PgNamespace is ready: %v", err)
	}

	// PgNamespace teardown logic
	if cr.DeletionTimestamp != nil {
		err := deletePgNamespace(cr)
		if err != nil {
			return err
		}

		// https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#advanced-topics
		// cr.SetFinalizers([]string{})
		// err := sdk.Update(cr)
		// if err != nil {
		// 	return fmt.Errorf("failed to update cr finalizers: %v", err)
		// }

		return nil
	}
	if !exists {
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

// IsPgNamespaceExists checks if the PgNamespace exists.
// Logs in as DB admin and queries the role table for existence
// of the namespace user.
// Returns a boolean and an error, if the connection did not succeed
func IsPgNamespaceExists(cr *api.PgNamespaceSpec) (bool, error) {

	dbNamespaceUser := cr.User
	db, err := getDbClusterConnection(&cr.PgCluster)
	if err != nil {
		return false, err
	}
	query := fmt.Sprintf("SELECT 1 FROM pg_roles WHERE rolname='%s'", dbNamespaceUser)
	row := db.QueryRow(query)
	db.Close()

	if row.Scan() != sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

// Package pg contains the reconciliation logic for the PgNamespace Custom Resource.
package pg

import (
	"fmt"

	"github.com/sirupsen/logrus"
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
)

func createPgNamespace(cr *api.PgNamespace) (err error) {
	// Get Cluster Connection
	db, err := getDbClusterConnection(&cr.Spec.PgCluster)
	if err != nil {
		logrus.Errorf("Failed to establish db connection: %v", err)
		return err
	}
	// Create PgNamespace
	dbNamespaceUser := cr.Spec.User
	dbNamespacePassword := cr.Spec.Password
	query := fmt.Sprintf("CREATE ROLE '%s' WITH CREATEDB PASSWORD '%s';", dbNamespaceUser, dbNamespacePassword)
	_, err2 := db.Exec(query)
	if err2 != nil {
		logrus.Errorf("Failed to execute query: %v", err2)
		return err2
	}
	return nil
}

func updatePgNamespace(cr *api.PgNamespace) (err error) {
	// Get Cluster Connection
	db, err := getDbClusterConnection(&cr.Spec.PgCluster)
	if err != nil {
		logrus.Errorf("Failed to establish db connection: %v", err)
		return err
	}
	// Update PgNamespace
	dbNamespaceUser := cr.Spec.User
	dbNamespacePassword := cr.Spec.Password
	query := fmt.Sprintf("ALTER ROLE '%s' WITH PASSWORD '%s';", dbNamespaceUser, dbNamespacePassword)
	_, err2 := db.Exec(query)
	if err2 != nil {
		logrus.Errorf("Failed to execute query: %v", err2)
		return err2
	}
	return nil
}

func deletePgNamespace(cr *api.PgNamespace) (err error) {
	return nil
}

// getPgNamespaceUsedQuota gets the PgNamespace's currently used quota.
// It queries the PgCluster on the size of all database objects owned by the
// PgNamespace and returns currently used quota.
func getPgNamespaceUsedQuota(cr *api.PgNamespace) (quota *int64, err error) {
	// Get Cluster Connection
	db, err := getDbClusterConnection(&cr.Spec.PgCluster)
	if err != nil {
		logrus.Errorf("Failed to establish db connection: %v", err)
		return nil, err
	}
	// Create PgNamespace
	dbNamespaceUser := cr.Spec.User
	query := fmt.Sprintf(`
		SELECT SUM(pg_database_size(datname))::bigint ) AS usage
		FROM pg_database
		JOIN pg_authid
		ON pg_database.datdba = pg_authid.oid
		WHERE rolname = '%s'`, dbNamespaceUser)
	row := db.QueryRow(query)
	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return nil, err
	}
	var usage int64
	row.Scan(&usage)
	return &usage, nil
}

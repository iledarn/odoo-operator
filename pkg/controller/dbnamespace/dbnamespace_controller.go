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

package dbnamespace

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	// Load pg lib
	_ "github.com/lib/pq"

	"github.com/xoe-labs/odoo-operator/pkg/finalizer"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new DBNamespace Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this cluster.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDBNamespace{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dbnamespace-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to DBNamespace
	err = c.Watch(&source.Kind{Type: &clusterv1beta1.DBNamespace{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDBNamespace{}

const (
	// DBMgtName ...
	DBMgtName = "postgres"
	// FinalizerKey ...
	FinalizerKey = "dbnamespace.odoo.io"
)

// ReconcileDBNamespace reconciles a DBNamespace object
type ReconcileDBNamespace struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DBNamespace object and makes changes based on the state read
// and what is in the DBNamespace.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=cluster.odoo.io,resources=dbnamespaces,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileDBNamespace) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Fetch the DBNamespace instance
	instance := &clusterv1beta1.DBNamespace{}
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
	exists, err := IsDBNamespaceExists(&instance.Spec)
	if err != nil {
		return reconcile.Result{}, err
	}
	hasFinalizer, err := finalizers.HasFinalizer(instance, FinalizerKey)
	if err != nil {
		return reconcile.Result{}, err
	}
	var operation string
	// Marked for deletion: tear down!
	if instance.GetDeletionTimestamp() != nil {
		operation = "delete"
		if hasFinalizer {
			if exists {
				err := deleteDBNamespace(&instance.Spec)
				if err != nil {
					log.Printf("%s/%s Operation: %s. (Controller: DBNameSpace) Error: %s\n", instance.Namespace, instance.Name, operation, err)
					return reconcile.Result{}, err
				}
			}
			finalizers.RemoveFinalizers(instance, sets.NewString(FinalizerKey))
		}
		log.Printf("%s/%s reconciled. Operation: %s. (Controller: DBNameSpace)\n", instance.Namespace, instance.Name, operation)
		return reconcile.Result{}, nil
	}

	if !exists {
		operation = "create"
		// Create or update DBNamespace.
		err := createDBNamespace(&instance.Spec)
		if err != nil {
			log.Printf("%s/%s Operation: %s. (Controller: DBNameSpace) Error: %s\n", instance.Namespace, instance.Name, operation, err)
			return reconcile.Result{}, err
		}
	} else {
		operation = "update"
		err := updateDBNamespace(&instance.Spec)
		if err != nil {
			log.Printf("%s/%s Operation: %s. (Controller: DBNameSpace) Error: %s\n", instance.Namespace, instance.Name, operation, err)
			return reconcile.Result{}, err
		}
	}

	log.Printf("%s/%s reconciled. Operation: %s. (Controller: DBNameSpace)\n", instance.Namespace, instance.Name, operation)
	return reconcile.Result{}, nil
}

// IsDBNamespaceExists checks if the DBNamespace exists.
// Logs in as DB admin and queries the role table for existence
// of the namespace user.
// Returns a boolean and an error, if the connection did not succeed
func IsDBNamespaceExists(spec *clusterv1beta1.DBNamespaceSpec) (bool, error) {

	dbNamespaceUser := spec.User
	db, err := getDbClusterConnection(spec)
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

func createDBNamespace(spec *clusterv1beta1.DBNamespaceSpec) (err error) {
	// Get Cluster Connection
	db, err := getDbClusterConnection(spec)
	if err != nil {
		return err
	}
	// Create PgNamespace
	query := fmt.Sprintf("CREATE ROLE '%s' WITH CREATEDB PASSWORD '%s';", spec.User, spec.Password)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func updateDBNamespace(spec *clusterv1beta1.DBNamespaceSpec) (err error) {
	// Get Cluster Connection
	db, err := getDbClusterConnection(spec)
	if err != nil {
		return err
	}
	// Update PgNamespace
	query := fmt.Sprintf("ALTER ROLE '%s' WITH PASSWORD '%s';", spec.User, spec.Password)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func deleteDBNamespace(spec *clusterv1beta1.DBNamespaceSpec) (err error) {
	return nil
}

// getPgNamespaceUsedQuota gets the PgNamespace's currently used quota.
// It queries the PgCluster on the size of all database objects owned by the
// PgNamespace and returns currently used quota.
func getDBNamespaceUsedQuota(spec *clusterv1beta1.DBNamespaceSpec) (quota *int64, err error) {
	// Get Cluster Connection
	db, err := getDbClusterConnection(spec)
	if err != nil {
		return nil, err
	}
	// Create PgNamespace
	query := fmt.Sprintf(`
		SELECT SUM(pg_database_size(datname))::bigint ) AS usage
		FROM pg_database
		JOIN pg_authid
		ON pg_database.datdba = pg_authid.oid
		WHERE rolname = '%s'`, spec.User)
	row := db.QueryRow(query)
	if err != nil {
		return nil, err
	}
	var usage int64
	row.Scan(&usage)
	return &usage, nil
}

func getDbClusterConnection(spec *clusterv1beta1.DBNamespaceSpec) (*sql.DB, error) {

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		spec.Host, spec.Port, spec.DBCluster.User, spec.DBCluster.Password, DBMgtName)
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

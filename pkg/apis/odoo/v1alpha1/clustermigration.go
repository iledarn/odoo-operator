package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ClusterMigration `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterMigration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ClusterMigrationSpec   `json:"spec"`
	Status            ClusterMigrationStatus `json:"status,omitempty"`
}

type ClusterMigrationSpec struct {
	OdooCluster       string    `json:odooCluster`
	MigratorImageSpec ImageSpec `json:migratorImageSpec`
	NewImageSpec      ImageSpec `json:newImageSpec`
}

type ClusterMigrationStatus struct {
	State   ClusterMigrationState `json:"state,omitempty"`
	Message string                `json:"message,omitempty"`
	// Additional Status
	OldInstances int32 `json:"oldInstances,omitempty"`
	NewInstances int32 `json:"newInstances,omitempty"`
}

// ClusterMigrationState ...
type ClusterMigrationState string

const (
	// ClusterMigrationStateCreated ...
	ClusterMigrationStateCreated ClusterMigrationState = "Created"
	// ClusterMigrationStateFinished ...
	ClusterMigrationStateFinished ClusterMigrationState = "Processed"
	// ClusterMigrationStateError ...
	ClusterMigrationStateError ClusterMigrationState = "Error"
	// ClusterMigrationStateMigrating ...
	ClusterMigrationStateMigrating ClusterMigrationState = "Migrating"
)

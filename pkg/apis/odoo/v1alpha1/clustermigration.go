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
	OdooCluster       string               `json:odooCluster`
	MigratorImageSpec OdooClusterSpecImage `json:migratorImageSpec`
	NewImageSpec      OdooClusterSpecImage `json:newImageSpec`
}

type ClusterMigrationStatus struct {
	OldInstances int                   `json:"oldInstances,omitempty"`
	NewInstances int                   `json:"newInstances,omitempty"`
	State        ClusterMigrationState `json:"state,omitempty"`
}

// ClusterMigrationState ...
type ClusterMigrationState string

const (
	// ClusterMigrationStateCreated ...
	ClusterMigrationStateCreated ClusterMigrationState = "Created"
	// ClusterMigrationStateProcessed ...
	ClusterMigrationStateProcessed ClusterMigrationState = "Processed"
)

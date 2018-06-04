package v1alpha1

import (
	"k8s.io/api/core/v1"
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
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ClusterMigrationCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// Additional Status
	// +optional
	OldInstances int32 `json:"oldInstances,omitempty"`
	// +optional
	NewInstances int32 `json:"newInstances,omitempty"`
}

type ClusterMigrationCondition struct {
	// Type is the type of the condition.
	Type            ClusterMigrationConditionType `json:"type"`
	StatusCondition `json:",inline"`
}

// ClusterMigrationConditionType ...
type ClusterMigrationConditionType string

const (
	// ClusterMigrationConditionTypeCreated ...
	ClusterMigrationConditionTypeCreated ClusterMigrationConditionType = "Created"
	// ClusterMigrationConditionTypeProcessed ...
	ClusterMigrationConditionTypeProcessed ClusterMigrationConditionType = "Processed"
	// ClusterMigrationConditionTypeErrored ...
	ClusterMigrationConditionTypeErrored ClusterMigrationConditionType = "Errored"
)

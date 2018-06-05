package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []InstanceMigration `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceMigration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              InstanceMigrationSpec   `json:"spec"`
	Status            InstanceMigrationStatus `json:"status,omitempty"`
}

type InstanceMigrationSpec struct {
	OdooInstance     string `json:odooInstance`
	ClusterMigration string `json:clusterMigration`
	RedirectInfoPage string `json:redirectInfoPage`
}

type InstanceMigrationStatus struct {
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []InstanceMigrationCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type InstanceMigrationCondition struct {
	// Type is the type of the condition.
	Type            InstanceMigrationConditionType `json:"type"`
	StatusCondition `json:",inline"`
}

// InstanceMigrationConditionType ...
type InstanceMigrationConditionType string

const (
	// InstanceMigrationConditionTypeCreated ...
	InstanceMigrationConditionTypeCreated InstanceMigrationConditionType = "Created"
	// InstanceMigrationConditionTypeProcessed ...
	InstanceMigrationConditionTypeProcessed InstanceMigrationConditionType = "Processed"
	// InstanceMigrationConditionTypeErrored ...
	InstanceMigrationConditionTypeErrored InstanceMigrationConditionType = "Errored"
)

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PgNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PgNamespace `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PgNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PgNamespaceSpec   `json:"spec"`
	Status            PgNamespaceStatus `json:"status,omitempty"`
}

type PgNamespaceSpec struct {
	User      string              `json:"user"`
	Password  string              `json:"password"`
	UserQuota v1.ResourceList     `json:"userQuota"`
	PgCluster PgClusterConnection `json:"pgCluster"`
}

type PgClusterConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PgNamespaceStatus struct {
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []PgNamespaceCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// Additional Status
	// +optional
	UsedQuota v1.ResourceList `json:"usedQuota,omitempty"`
}

type PgNamespaceCondition struct {
	// Type is the type of the condition.
	Type            PgNamespaceConditionType `json:"type"`
	StatusCondition `json:",inline"`
}

// PgNamespaceConditionType ...
type PgNamespaceConditionType string

const (
	// PgNamespaceConditionTypeCreated ...
	PgNamespaceConditionTypeCreated PgNamespaceConditionType = "Created"
	// PgNamespaceConditionTypeReconciled ...
	PgNamespaceConditionTypeReconciled PgNamespaceConditionType = "Reconciled"
	// PgNamespaceConditionTypeErrored ...
	PgNamespaceConditionTypeErrored PgNamespaceConditionType = "Errored"
)

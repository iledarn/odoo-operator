package v1alpha1

import (
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
	UserQuota int                 `json:"userQuota"`
	PgCluster PgClusterConnection `json:"pgCluster"`
}

type PgClusterConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PgNamespaceStatus struct {
	State   PgNamespaceState `json:"state,omitempty"`
	Message string           `json:"message,omitempty"`
	// Additional Status
	UsedQuota int `json:"usedQuota"`
}

// PgNamespaceState ...
type PgNamespaceState string

const (
	// PgNamespaceStateCreated ...
	PgNamespaceStateCreated PgNamespaceState = "Created"
	// PgNamespaceStateReconciled ...
	PgNamespaceStateReconciled PgNamespaceState = "Reconciled"
	// PgNamespaceStateError ...
	PgNamespaceStateError PgNamespaceState = "Error"
)

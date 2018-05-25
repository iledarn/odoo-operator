package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OdooInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []OdooInstance `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OdooInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              OdooInstanceSpec   `json:"spec"`
	Status            OdooInstanceStatus `json:"status,omitempty"`
}

type OdooInstanceSpec struct {
	OdooCluster string `json:odooCluster`
	// DbName         string      `json:dbName`
	HostName       string `json:hostName`
	DbSeedCfgMap   string `json:"dbSeedCfgMap"`
	DbQuota        int    `json:"dbQuota"`
	FilestoreQuota int    `json:"fsQuota"`
}

type OdooInstanceStatus struct {
	State   OdooInstanceState `json:"state,omitempty"`
	Message string            `json:"message,omitempty"`
	// Additional Status
	UsedDbQuota int `json:"usedDbQuota,omitempty"`
	UsedFsQuota int `json:"usedFsQuota,omitempty"`
}

// OdooInstanceState ...
type OdooInstanceState string

const (
	// OdooInstanceStateCreated ...
	OdooInstanceStateCreated OdooInstanceState = "Created"
	// OdooInstanceStateReconciled ...
	OdooInstanceStateReconciled OdooInstanceState = "Reconciled"
	// OdooInstanceStateError ...
	OdooInstanceStateError OdooInstanceState = "Error"
	// OdooInstanceStateSuspended ...
	OdooInstanceStateSuspended OdooInstanceState = "Suspended"
	// OdooInstanceStateMaintenance ...
	OdooInstanceStateMaintenance OdooInstanceState = "Maintenance"
)

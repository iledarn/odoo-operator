package v1alpha1

import (
	"k8s.io/api/core/v1"
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
	DbQuota        int32  `json:"dbQuota"`
	FilestoreQuota int32  `json:"fsQuota"`
}

type OdooInstanceStatus struct {
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []OdooInstanceCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// Additional Status
	// +optional
	UsedDbQuota int32 `json:"usedDbQuota,omitempty"`
	// +optional
	UsedFsQuota int32 `json:"usedFsQuota,omitempty"`
}

type OdooInstanceCondition struct {
	// Type is the type of the condition.
	Type            OdooInstanceConditionType `json:"type"`
	StatusCondition `json:",inline"`
}

// OdooInstanceConditionType ...
type OdooInstanceConditionType string

const (
	// OdooInstanceConditionTypeCreated ...
	OdooInstanceConditionTypeCreated OdooInstanceConditionType = "Created"
	// OdooInstanceConditionTypeReconciled ...
	OdooInstanceConditionTypeReconciled OdooInstanceConditionType = "Reconciled"
	// OdooInstanceConditionTypeErrored ...
	OdooInstanceConditionTypeErrored OdooInstanceConditionType = "Errored"
	// OdooInstanceConditionTypeSuspended ...
	OdooInstanceConditionTypeSuspended OdooInstanceConditionType = "Suspended"
	// OdooInstanceConditionTypeMaintenance ...
	OdooInstanceConditionTypeMaintenance OdooInstanceConditionType = "Maintenance"
)

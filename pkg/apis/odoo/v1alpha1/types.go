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
	Image            OdooInstanceSpecImage    `json:image`
	DbSpec           OdooInstanceSpecDbSpec   `json:dbSpec`
	ResourceSpec     OdooInstanceResourceSpec `json:resourceSpec`
	AdminPassword    string                   `json:"adminPassword"`
	ConfigMap        string                   `json:"configMap"`
	SeedingConfigMap string                   `json:"seedingConfigMap"`
	// Replicas         int                      `json:"replicas"`

	// MailServer  bool `json:"mailServer"`
	// OnlyOffice  bool `json:"onlyOffice"`
	// Mattermost  bool `json:"mattermost"`
	// Nuxeo       bool `json:"nuxeo"`
	// BpmnEngine  bool `json:"bpmnEngine"`
	// OpenProject bool `json:"openProject"`
}

type OdooInstanceSpecImage struct {
	Registry string `json:"registry"`
	Name     string `json:"image"`
	Tag      string `json:"tag"`
}

type OdooInstanceSpecDbSpec struct {
	// TODO: Enforce DbQuota compiled by a db Cronjob
	// Using mtDatabase + pg_cron and dbQuota +
	// https://stackoverflow.com/a/37822365
	User        string `json:"user"`
	Password    string `json:"password"`
	MgtDatabase string `json:"mgtDatabase"`
	UserQuota   int    `json:"userQuota"`
}

type OdooInstanceResourceSpec struct {
	Cpu  int `json:"cpu"`
	Ram  int `json:"ram"`
	Disk int `json:"disk"`
}
type OdooInstanceStatus struct {
	DbQuotaUsage  string            `json:"dbQuotaUsage,omitempty"`
	DiskUsage     string            `json:"diskUsage,omitempty"`
	State         OdooInstanceState `json:"state,omitempty"`
	Message       string            `json:"message,omitempty"`
	ImageVersions []string          `json:"imageVersions,omitempty"`
	// Replicas     int               `json:"replicas,omitempty"`
}

// OdooInstanceState ...
type OdooInstanceState string

const (
	// OdooInstanceStateReconciled ...
	OdooInstanceStateReconciled OdooInstanceState = "Reconciled"
	// OdooInstanceStateMigrating ...
	OdooInstanceStateMigrating OdooInstanceState = "Migrating"
)

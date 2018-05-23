package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OdooClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []OdooCluster `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OdooCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              OdooClusterSpec   `json:"spec"`
	Status            OdooClusterStatus `json:"status,omitempty"`
}

type OdooClusterSpec struct {
	ImageSpec        OdooClusterSpecImage    `json:imageSpec`
	DbSpec           OdooClusterSpecDbSpec   `json:dbSpec`
	ResourceSpec     OdooClusterResourceSpec `json:resourceSpec`
	AdminPassword    string                  `json:"adminPassword"`
	ConfigMap        string                  `json:"configMap"`
	SeedingConfigMap string                  `json:"seedingConfigMap"`
	// Replicas         int                      `json:"replicas"`

	// MailServer  bool `json:"mailServer"`
	// OnlyOffice  bool `json:"onlyOffice"`
	// Mattermost  bool `json:"mattermost"`
	// Nuxeo       bool `json:"nuxeo"`
	// BpmnEngine  bool `json:"bpmnEngine"`
	// OpenProject bool `json:"openProject"`
}

type OdooClusterSpecImage struct {
	Registry string `json:"registry"`
	Name     string `json:"image"`
	Tag      string `json:"tag"`
}

type OdooClusterSpecDbSpec struct {
	// TODO: Enforce DbQuota compiled by a db Cronjob
	// Using mtDatabase + pg_cron and dbQuota +
	// https://stackoverflow.com/a/37822365
	User        string `json:"user"`
	Password    string `json:"password"`
	MgtDatabase string `json:"mgtDatabase"`
	UserQuota   int    `json:"userQuota"`
}

type OdooClusterResourceSpec struct {
	Cpu  int `json:"cpu"`
	Ram  int `json:"ram"`
	Disk int `json:"disk"`
}
type OdooClusterStatus struct {
	DbQuotaUsage  string           `json:"dbQuotaUsage,omitempty"`
	DiskUsage     string           `json:"diskUsage,omitempty"`
	State         OdooClusterState `json:"state,omitempty"`
	Message       string           `json:"message,omitempty"`
	ImageVersions []string         `json:"imageVersions,omitempty"`
	// Replicas     int               `json:"replicas,omitempty"`
}

// OdooClusterState ...
type OdooClusterState string

const (
	// OdooClusterStateReconciled ...
	OdooClusterStateReconciled OdooClusterState = "Reconciled"
	// OdooClusterStateMigrating ...
	OdooClusterStateMigrating OdooClusterState = "Migrating"
)

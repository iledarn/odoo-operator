package v1alpha1

import (
	"k8s.io/api/core/v1"
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
	Tracks            []TrackSpec           `json:tracks`
	Tiers             []TierSpec            `json:tiers`
	PgSpec            PgNamespaceSpec       `json:pgNsSpec`
	ResourceQuotaSpec *v1.ResourceQuotaSpec `json:resourceQuotaSpec,omitempty`
	AdminPassword     string                `json:"adminPassword"`
	PgPassFile        string                `json:"pgPassFile"`
	ConfigMap         string                `json:"configMap"`
	DeployModel       OdooClusterMode       `json:deployModel`
	NodeSelector      string                `json:"nodeSelector"`
	StorageClassName  *string               `json:"storageClassNam"`

	// MailServer  bool `json:"mailServer"`
	// OnlyOffice  bool `json:"onlyOffice"`
	// Mattermost  bool `json:"mattermost"`
	// Nuxeo       bool `json:"nuxeo"`
	// BpmnEngine  bool `json:"bpmnEngine"`
	// OpenProject bool `json:"openProject"`
}

type TrackSpec struct {
	Name  string    `json:"name"`
	Image ImageSpec `json:"image"`
}

type TierSpec struct {
	Name     Tier  `json:"name"`
	Replicas int32 `json:"replicas"`
}

type Tier string

const (
	ServerTier      Tier = "Server"
	CronTier        Tier = "Cron"
	BackgroundTier  Tier = "Background"
	LongpollingTier Tier = "LongPolling"
)

// OdooClusterMode ...
type OdooClusterMode string

const (
	// OdooClusterModeRemote ...
	OdooClusterModeRemote OdooClusterMode = "Remote"
	// OdooClusterModeLocal ...
	OdooClusterModeLocal OdooClusterMode = "Local"
	// OdooClusterModeHybrid ...
	OdooClusterModeHybrid OdooClusterMode = "Hybrid"
)

type OdooClusterStatus struct {
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []OdooClusterCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type OdooClusterCondition struct {
	// Type is the type of the condition.
	Type            OdooClusterConditionType `json:"type"`
	StatusCondition `json:",inline"`
}

// OdooClusterConditionType ...
type OdooClusterConditionType string

const (
	// OdooClusterConditionTypeCreated ...
	OdooClusterConditionTypeCreated OdooClusterConditionType = "Created"
	// OdooClusterConditionTypeReconciled ...
	OdooClusterConditionTypeReconciled OdooClusterConditionType = "Reconciled"
	// OdooClusterConditionTypeErrored ...
	OdooClusterConditionTypeErrored OdooClusterConditionType = "Errored"
)

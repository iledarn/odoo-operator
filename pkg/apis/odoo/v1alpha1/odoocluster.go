package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultTrackName     = "stable"
	defaultRegistry      = "docker.io"
	defaultImage         = "xoe-labs/odoo"
	defaultTag           = "latest"
	defaultAdminPassword = "admin-password"
	defaultPgPassFile    = "db:*:*:odoo:odoo"
	defaultClusterMode   = OdooClusterModeRemote
)

var defultImage = ImageSpec{
	Registry: defaultRegistry,
	Image:    defaultImage,
	Tag:      defaultTag,
}

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

// SetDefaults sets the default values for the OdooCluster spec and returns true if the spec was changed
func (c *OdooCluster) SetDefaults() bool {
	changed := false
	cs := &c.Spec
	if len(cs.Tracks) == 0 {
		cs.Tracks = []TrackSpec{{
			Name:  defaultTrackName,
			Image: defultImage,
		}}
		changed = true
	}
	if len(cs.Tiers) == 0 {
		cs.Tiers = []TierSpec{{Name: ServerTier}}
		changed = true
	}
	if cs.PVCSpecs == nil {
		cs.PVCSpecs = []PVCSpec{
			{Name: PVCNamePersistence},
			{Name: PVCNameBackup},
		}
		changed = true
	}
	if len(cs.AdminPassword) == 0 {
		cs.AdminPassword = defaultAdminPassword
		changed = true
	}
	if len(cs.PgPassFile) == 0 {
		cs.PgPassFile = defaultPgPassFile
		changed = true
	}
	if cs.DeployModel == "" {
		cs.DeployModel = OdooClusterModeRemote
		changed = true
	}
	return changed
}

type OdooClusterSpec struct {
	Tracks            []TrackSpec           `json:tracks`
	Tiers             []TierSpec            `json:tiers`
	PVCSpecs          []PVCSpec             `json:pvcs,omitempty`
	PgSpec            PgNamespaceSpec       `json:pgNsSpec`
	ResourceQuotaSpec *v1.ResourceQuotaSpec `json:resourceQuotaSpec,omitempty`
	AdminPassword     string                `json:"adminPassword"`
	PgPassFile        string                `json:"pgPassFile"`
	ConfigMap         string                `json:"configMap"`
	DeployModel       OdooClusterMode       `json:deployModel,omitempty`
	NodeSelector      *string               `json:"nodeSelector,omitempy"`

	// MailServer  bool `json:"mailServer"`
	// OnlyOffice  bool `json:"onlyOffice"`
	// Mattermost  bool `json:"mattermost"`
	// Nuxeo       bool `json:"nuxeo"`
	// BpmnEngine  bool `json:"bpmnEngine"`
	// OpenProject bool `json:"openProject"`
}

type PVCSpec struct {
	Name PVCName `json:"name"`
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// +optional
	StorageClassName *string `json:"storageClassNam"`
}

type TrackSpec struct {
	Name  string    `json:"name"`
	Image ImageSpec `json:"image"`
}

type TierSpec struct {
	Name Tier `json:"name"`
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
	// +optional
	QOS *v1.PodQOSClass `json:"qos,omitempty"`
	// +optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`
}

type PVCName string

const (
	PVCNamePersistence PVCName = "Persistence"
	PVCNameBackup      PVCName = "Backup"
)

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

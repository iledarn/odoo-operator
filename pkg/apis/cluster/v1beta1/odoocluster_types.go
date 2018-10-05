/*
 * This file is part of the Odoo-Operator (R) project.
 * Copyright (c) 2018-2018 XOE Corp. SAS
 * Authors: David Arnold, et al.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * ALTERNATIVE LICENCING OPTION
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving the Odoo-Operator software without
 * disclosing the source code of your own applications. These activities
 * include: Offering paid services to a customer as an ASP, shipping Odoo-
 * Operator with a closed source product.
 *
 */

package v1beta1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OdooClusterSpec defines the desired state of OdooCluster
type OdooClusterSpec struct {
	Tracks        []TrackSpec     `json:"tracks"`
	Tiers         []TierSpec      `json:"tiers"`
	Volumes       []Volume        `json:"volumes,omitempty"`
	DBSpec        DBNamespaceSpec `json:"dbNamespaceSpec"`
	AdminPassword string          `json:"adminPassword"`
	// +optional
	ResourceQuotaSpec *v1.ResourceQuotaSpec `json:"resourceQuotaSpec,omitempty"`
	// +optional
	Config *string `json:"config,omitempty"`
	// +optional
	IntegratorConfig *string `json:"integratorConfig,omitempty"`
	// +optional
	CustomConfig *string `json:"customConfig,omitempty"`
	// +optional
	NodeSelector *v1.NodeSelector `json:"nodeSelector,omitempty"`

	// MailServer  bool `json:"mailServer"`
	// OnlyOffice  bool `json:"onlyOffice"`
	// Mattermost  bool `json:"mattermost"`
	// Nuxeo       bool `json:"nuxeo"`
	// BpmnEngine  bool `json:"bpmnEngine"`
	// OpenProject bool `json:"openProject"`
}

// TrackSpec defines the desired state of a Track
type TrackSpec struct {
	Name  string    `json:"name"`
	Image ImageSpec `json:"image"`
	// +optional
	Config *string `json:"config,omitempty"`
	// +optional
	IntegratorConfig *string `json:"integratorConfig,omitempty"`
	// +optional
	CustomConfig *string `json:"customConfig,omitempty"`
}

// TierSpec defines the desired state of a Tier
type TierSpec struct {
	Name Tier `json:"name"`
	// +optional
	DBConn *int32 `json:"dbConn,omitempty"`
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
	// +optional
	QOS *v1.PodQOSClass `json:"qos,omitempty"`
	// +optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`
}

// Volume defines the desired state of a Volume
type Volume struct {
	Name VolumeName                   `json:"name"`
	Spec v1.PersistentVolumeClaimSpec `json:"spec"`
}

// VolumeName can either be "Data" or "Backup"
type VolumeName string

const (
	// PVCNameData ...
	PVCNameData VolumeName = "Data"
	// PVCNameBackup ...
	PVCNameBackup VolumeName = "Backup"
)

// Tier can either be "Server", "Cron" or "LongPolling"
type Tier string

const (
	// ServerTier ...
	ServerTier Tier = "Server"
	// CronTier ...
	CronTier Tier = "Cron"
	// LongpollingTier ...
	LongpollingTier Tier = "LongPolling"
)

// ImageSpec defines an Image and (optionally) it's registry credentials
type ImageSpec struct {
	Registry string `json:"registry"`
	Image    string `json:"image"`
	Tag      string `json:"tag"`
	Secret   string `json:"secret,omitempty"`
}

// OdooClusterStatus defines the observed state of OdooCluster
type OdooClusterStatus struct {
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []OdooClusterStatusCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// OdooClusterStatusCondition defines an observable OdooClusterStatus condition
type OdooClusterStatusCondition struct {
	// Type of the OdooClusterStatus condition.
	Type OdooClusterStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=OdooClusterStatusConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// OdooClusterStatusConditionType ...
type OdooClusterStatusConditionType string

const (
	// OdooClusterStatusConditionTypeCreated ...
	OdooClusterStatusConditionTypeCreated OdooClusterStatusConditionType = "Created"
	// OdooClusterStatusConditionTypeReconciled ...
	OdooClusterStatusConditionTypeReconciled OdooClusterStatusConditionType = "Reconciled"
	// OdooClusterStatusConditionTypeErrored ...
	OdooClusterStatusConditionTypeErrored OdooClusterStatusConditionType = "Errored"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OdooCluster is the Schema for the odooclusters API
// +k8s:openapi-gen=true
type OdooCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OdooClusterSpec   `json:"spec,omitempty"`
	Status OdooClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OdooClusterList contains a list of OdooCluster
type OdooClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OdooCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OdooCluster{}, &OdooClusterList{})
}

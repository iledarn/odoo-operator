package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []InstanceBackup `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              InstanceBackupSpec   `json:"spec"`
	Status            InstanceBackupStatus `json:"status,omitempty"`
}

type InstanceBackupSpec struct {
	OdooInstance   string                 `json:odooInstance`
	BackupStrategy InstanceBackupStrategy `json:backupStrategy`
	BackupTarget   InstanceBackupTarget   `json:backupTarget`
}

type InstanceBackupStrategy struct {
	Name          InstanceBackupStrategyName `json:name`
	StorageTarget string                     `json:storageTarget` // package strategy only
}

// InstanceBackupStrategyName ...
type InstanceBackupStrategyName string

const (
	// InstanceBackupStrategyNameMandate ...
	InstanceBackupStrategyNameMandate InstanceBackupStrategyName = "mandate"
	// InstanceBackupStrategyNamePackage ...
	InstanceBackupStrategyNamePackage InstanceBackupStrategyName = "package"
)

// InstanceBackupTarget ...
type InstanceBackupTarget string

const (
	// InstanceBackupTargetDB ...
	InstanceBackupTargetDB InstanceBackupTarget = "DB"
	// InstanceBackupTargetFS ...
	InstanceBackupTargetFS InstanceBackupTarget = "FS"
	// InstanceBackupTargetAll ...
	InstanceBackupTargetAll InstanceBackupTarget = "all"
)

type InstanceBackupStatus struct {
	State InstanceBackupState `json:"state,omitempty"`
}

// InstanceBackupState ...
type InstanceBackupState string

const (
	// InstanceBackupStateStarted ...
	InstanceBackupStateStarted InstanceBackupState = "Started"
	// InstanceBackupStateVerified ...
	InstanceBackupStateVerified InstanceBackupState = "Verified"
)

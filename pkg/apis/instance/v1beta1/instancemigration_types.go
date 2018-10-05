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

// InstanceMigrationSpec defines the desired state of InstanceMigration
type InstanceMigrationSpec struct {
	OdooInstance     string `json:"odooInstance"`
	ClusterMigration string `json:"clusterMigration"`
	RedirectInfoPage string `json:"redirectInfoPage"`
}

// InstanceMigrationStatus defines the observed state of InstanceMigration
type InstanceMigrationStatus struct {
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []InstanceMigrationStatusCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// InstanceMigrationStatusCondition defines an observable InstanceMigrationStatus condition
type InstanceMigrationStatusCondition struct {
	// Type of the InstanceMigrationStatus condition.
	Type InstanceMigrationStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=InstanceMigrationStatusConditionType"`
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

// InstanceMigrationStatusConditionType ...
type InstanceMigrationStatusConditionType string

const (
	// InstanceMigrationStatusConditionTypeCreated ...
	InstanceMigrationStatusConditionTypeCreated InstanceMigrationStatusConditionType = "Created"
	// InstanceMigrationStatusConditionTypeProcessed ...
	InstanceMigrationStatusConditionTypeProcessed InstanceMigrationStatusConditionType = "Processed"
	// InstanceMigrationStatusConditionTypeErrored ...
	InstanceMigrationStatusConditionTypeErrored InstanceMigrationStatusConditionType = "Errored"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceMigration is the Schema for the instancemigrations API
// +k8s:openapi-gen=true
type InstanceMigration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceMigrationSpec   `json:"spec,omitempty"`
	Status InstanceMigrationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceMigrationList contains a list of InstanceMigration
type InstanceMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstanceMigration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InstanceMigration{}, &InstanceMigrationList{})
}

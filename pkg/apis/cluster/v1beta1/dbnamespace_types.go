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

// DBNamespaceSpec defines the desired state of DBNamespace
type DBNamespaceSpec struct {
	Host      string             `json:"host"`
	Port      string             `json:"port"`
	User      string             `json:"user"`
	Password  string             `json:"password"`
	DBCluster DBAdminCredentials `json:"dbAdmin"`
	UserQuota v1.ResourceList    `json:"userQuota,omitempty"`
}

// DBAdminCredentials defines the DB admin credentials
type DBAdminCredentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// DBNamespaceStatus defines the observed state of DBNamespace
type DBNamespaceStatus struct {
	// Current service state of apiService.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []DBNamespaceStatusCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// Additional Status
	// +optional
	UsedQuota v1.ResourceList `json:"usedQuota,omitempty"`
}

// DBNamespaceStatusCondition defines an observable DBNamespaceStatus condition
type DBNamespaceStatusCondition struct {
	// Type of the DBNamespaceStatus condition.
	Type DBNamespaceStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=DBNamespaceStatusConditionType"`
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

// DBNamespaceStatusConditionType ...
type DBNamespaceStatusConditionType string

const (
	// DBNamespaceStatusConditionTypeCreated ...
	DBNamespaceStatusConditionTypeCreated DBNamespaceStatusConditionType = "Created"
	// DBNamespaceStatusConditionTypeReconciled ...
	DBNamespaceStatusConditionTypeReconciled DBNamespaceStatusConditionType = "Reconciled"
	// DBNamespaceStatusConditionTypeErrored ...
	DBNamespaceStatusConditionTypeErrored DBNamespaceStatusConditionType = "Errored"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBNamespace is the Schema for the dbnamespaces API
// +k8s:openapi-gen=true
type DBNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBNamespaceSpec   `json:"spec,omitempty"`
	Status DBNamespaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBNamespaceList contains a list of DBNamespace
type DBNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBNamespace{}, &DBNamespaceList{})
}

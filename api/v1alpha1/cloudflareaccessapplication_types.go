/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/bojanzelic/cloudflare-zero-trust-operator/internal/cfcollections"
	cloudflare "github.com/cloudflare/cloudflare-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CloudflareAccessApplicationSpec defines the desired state of CloudflareAccessApplication.
type CloudflareAccessApplicationSpec struct {
	// Name of the Cloudflare Access Application
	Name string `json:"name"`

	// The domain and path that Access will secure.
	// ex: "test.example.com/admin"
	Domain string `json:"domain"`

	// The application type. Defaults to "self_hosted"
	// +optional
	// +kubebuilder:default=self_hosted
	Type cloudflare.AccessApplicationType `json:"type,omitempty"`

	// Displays the application in the App Launcher.
	// +optional
	// +kubebuilder:default=true
	AppLauncherVisible *bool `json:"appLauncherVisible,omitempty"`

	// The identity providers your users can select when connecting to this application. Defaults to all IdPs configured in your account.
	// ex: ["699d98642c564d2e855e9661899b7252"]
	// +optional
	AllowedIdps []string `json:"allowedIdps,omitempty"`

	// When set to true, users skip the identity provider selection step during login.
	// You must specify only one identity provider in allowed_idps.
	// +optional
	AutoRedirectToIdentity *bool `json:"autoRedirectToIdentity,omitempty"`

	// Policies is the ordered set of policies that should be applied to the application
	// Order determines precidence
	// +optional
	Policies CloudflareAccessPolicyList `json:"policies,omitempty"`

	// SessionDuration is the length of the session duration. Defaults to 24h
	// +optional
	// +kubebuilder:default=24h
	SessionDuration string `json:"sessionDuration,omitempty"`

	// Enables the binding cookie, which increases security against compromised authorization tokens and CSRF attacks.
	// +optional
	// +kubebuilder:default=false
	EnableBindingCookie *bool `json:"enableBindingCookie,omitempty"`

	// Enables the HttpOnly cookie attribute, which increases security against XSS attacks. Defaults to true
	// +optional
	// +kubebuilder:default=true
	HTTPOnlyCookieAttribute *bool `json:"httpOnlyCookieAttribute,omitempty"`
}

type CloudflareAccessPolicy struct {
	Name     string                      `json:"name"`
	Decision string                      `json:"decision"`
	Include  []CloudFlareAccessGroupRule `json:"include,omitempty"`
	Require  []CloudFlareAccessGroupRule `json:"require,omitempty"`
	Exclude  []CloudFlareAccessGroupRule `json:"exclude,omitempty"`

	// PurposeJustificationRequired *bool                 `json:"purpose_justification_required,omitempty"`
	// PurposeJustificationPrompt   *string               `json:"purpose_justification_prompt,omitempty"`
	// ApprovalRequired             *bool                 `json:"approval_required,omitempty"`
	// ApprovalGroups               []cloudflare.AccessApprovalGroup `json:"approval_groups"`
}

type CloudflareAccessPolicyList []CloudflareAccessPolicy

func (aps CloudflareAccessPolicyList) ToCloudflare() cfcollections.AccessPolicyCollection {
	ret := cfcollections.AccessPolicyCollection{}

	for i, policy := range aps {
		transformed := cloudflare.AccessPolicy{
			Name:       policy.Name,
			Precedence: i + 1,
			Decision:   policy.Decision,
		}

		managedCRFields := CloudFlareAccessGroupRuleGroups{
			policy.Include,
			policy.Exclude,
			policy.Require,
		}

		managedCFFields := []*[]interface{}{
			&transformed.Include,
			&transformed.Exclude,
			&transformed.Require,
		}

		managedCRFields.TransformCloudflareRuleFields(managedCFFields)

		ret = append(ret, transformed)
	}

	return ret
}

// CloudflareAccessApplicationStatus defines the observed state of CloudflareAccessApplication.
type CloudflareAccessApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	AccessApplicationID string      `json:"accessApplicationId,omitempty"`
	CreatedAt           metav1.Time `json:"createdAt,omitempty"`
	UpdatedAt           metav1.Time `json:"updatedAt,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CloudflareAccessApplication is the Schema for the cloudflareaccessapplications API.
type CloudflareAccessApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudflareAccessApplicationSpec   `json:"spec,omitempty"`
	Status CloudflareAccessApplicationStatus `json:"status,omitempty"`
}

func (c *CloudflareAccessApplication) CloudflareName() string {
	return c.ObjectMeta.Name + " [K8s]"
}

func (c *CloudflareAccessApplication) ToCloudflare() cloudflare.AccessApplication {
	app := cloudflare.AccessApplication{
		Name:                    c.CloudflareName(),
		ID:                      c.Status.AccessApplicationID,
		CreatedAt:               &c.Status.CreatedAt.Time,
		UpdatedAt:               &c.Status.UpdatedAt.Time,
		Domain:                  c.Spec.Domain,
		Type:                    c.Spec.Type,
		AppLauncherVisible:      c.Spec.AppLauncherVisible,
		AllowedIdps:             c.Spec.AllowedIdps,
		AutoRedirectToIdentity:  c.Spec.AutoRedirectToIdentity,
		SessionDuration:         c.Spec.SessionDuration,
		EnableBindingCookie:     c.Spec.EnableBindingCookie,
		HttpOnlyCookieAttribute: c.Spec.HTTPOnlyCookieAttribute,
	}

	return app
}

//+kubebuilder:object:root=true

// CloudflareAccessApplicationList contains a list of CloudflareAccessApplication.
type CloudflareAccessApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudflareAccessApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudflareAccessApplication{}, &CloudflareAccessApplicationList{})
}

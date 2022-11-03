/*
Copyright 2021 The Crossplane Authors.

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

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ghapps "github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v48/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfig gets the config.
func GetConfig(ctx context.Context, c client.Client, mg resource.Managed) ([]byte, error) {
	pc := &v1beta1.ProviderConfig{}
	if err := c.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced ProviderConfig")
	}

	t := resource.NewProviderConfigUsageTracker(c, &v1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	return resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, c, pc.Spec.Credentials.CommonCredentialSelectors)
}

// newClientFields helps to create new client by differents credentials
type newClientFields struct {
	AppID          int64  `json:"appId,omitempty"`
	InstallationID int64  `json:"installationId,omitempty"`
	PEMFile        string `json:"pemFile,omitempty"`
	PAT            string `json:"token,omitempty"`
}

// NewClient creates a new client.
func NewClient(token string) *github.Client {
	ctx := context.Background()
	creds := newClientFields{}
	var tc *http.Client

	if err := json.Unmarshal([]byte(token), &creds); err != nil {
		fmt.Println(err)
	}

	if creds.PAT != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: creds.PAT},
		)

		tc = oauth2.NewClient(ctx, ts)
	} else {
		ts, err := ghapps.New(http.DefaultTransport, creds.AppID, creds.InstallationID, []byte(creds.PEMFile))
		if err != nil {
			fmt.Println(err)
		}

		tc = &http.Client{Transport: ts}
	}

	return github.NewClient(tc)
}

// StringPtr converts the supplied string to a pointer to that string.
func StringPtr(p string) *string { return &p }

// StringValue converts the supplied pointer string to a string.
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// ConvertTimestamp converts *github.Timestamp into *metav1.Time
func ConvertTimestamp(t *github.Timestamp) *metav1.Time {
	if t == nil {
		return nil
	}
	return &metav1.Time{
		Time: t.Time,
	}
}

// Int64Value converts the supplied pointer int64 to a int64.
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// IntValue converts the supplied pointer int to a int.
func IntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// BoolValue converts the supplied pointer bool to a bool.
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// LateInitializeString implements late initialization for string type.
func LateInitializeString(s *string, from string) *string {
	if s != nil || from == "" {
		return s
	}
	return &from
}

// LateInitializeInt implements late initialization for int type.
func LateInitializeInt(i *int, from int) *int {
	if i != nil || from == 0 {
		return i
	}
	return &from
}

// LateInitializeBool implements late initialization for bool type.
func LateInitializeBool(b *bool, from bool) *bool {
	if b != nil || !from {
		return b
	}
	return &from
}

// LateInitializeStringSlice implements late initialization for
// string slice type.
func LateInitializeStringSlice(s []string, from []string) []string {
	if len(s) != 0 || len(from) == 0 {
		return s
	}
	return from
}

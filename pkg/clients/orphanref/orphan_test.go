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

package orphanref

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v48/github"

	"github.com/crossplane-contrib/provider-github/apis/gitdatabase/v1alpha1"
)

var (
	ref = "refs/heads/fake-branch"
	url = "example.com"
)

func TestGenerateObservation(t *testing.T) {
	type args struct {
		ref github.Reference
	}
	cases := map[string]struct {
		args
		out v1alpha1.OrphanRefObservation
	}{
		"Must generate an OrphanRefObservation based on the given github.Reference": {
			args: args{
				ref: github.Reference{
					URL: &url,
					Ref: &ref,
				},
			},
			out: v1alpha1.OrphanRefObservation{
				URL: url,
				Ref: ref,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateObservation(&tc.args.ref)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Errorf("GenerateObservation(...): -want, +got:\n%s", diff)
			}
		})
	}
}

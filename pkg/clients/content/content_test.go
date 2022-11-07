/*
Copyright 2021 The Crossplane Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package content

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v48/github"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
)

var (
	content = "test"
	url     = "example.com"
	sha     = "example-sha"
)

func TestOverrideParameters(t *testing.T) {
	type args struct {
		params  v1alpha1.ContentParameters
		content github.RepositoryContent
	}
	cases := map[string]struct {
		args
		out github.RepositoryContent
	}{
		"Must create a github.RepositoryContent from ContentParameters": {
			args: args{
				params: v1alpha1.ContentParameters{
					Content: content,
				},
				content: github.RepositoryContent{},
			},
			out: github.RepositoryContent{
				Content: &content,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := OverrideParameters(tc.args.params, tc.args.content)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Errorf("OverrideParameters(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsUpToDate(t *testing.T) {
	type args struct {
		params  *v1alpha1.ContentParameters
		content *github.RepositoryContent
	}
	cases := map[string]struct {
		args
		out bool
		err error
	}{
		"NotUpToDate": {
			args: args{
				params: &v1alpha1.ContentParameters{
					Content: content,
				},
				content: &github.RepositoryContent{},
			},
			out: false,
		},
		"UpToDate": {
			args: args{
				params: &v1alpha1.ContentParameters{
					Content: content,
				},
				content: &github.RepositoryContent{
					Content: &content,
				},
			},
			out: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, _ := IsUpToDate(tc.args.params, tc.args.content)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Errorf("IsUpToDate(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateObservation(t *testing.T) {
	type args struct {
		content github.RepositoryContent
	}
	cases := map[string]struct {
		args
		out v1alpha1.ContentObservation
	}{
		"Must generate an ContentObservation based on the given github.RepositoryContent": {
			args: args{
				content: github.RepositoryContent{
					URL:     &url,
					HTMLURL: &url,
					SHA:     &sha,
				},
			},
			out: v1alpha1.ContentObservation{
				URL:     url,
				HTMLURL: url,
				SHA:     sha,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateObservation(tc.args.content)
			if diff := cmp.Diff(tc.out, got); diff != "" {
				t.Errorf("GenerateObservation(...): -want, +got:\n%s", diff)
			}
		})
	}
}

// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudwatch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalConfig(t *testing.T) {
	testSuite := []struct {
		title     string
		jason     string
		yamele    string
		result    Config
		expectErr bool
	}{
		{
			title:  "default identity",
			jason:  `{"region": "us-east-1"}`,
			yamele: `region: us-east-1`,
			result: Config{Region: "us-east-1"},
		},
		{
			title: "assumed role with external ID",
			jason: `{"region": "eu-west-3", "roleArn": "arn:aws:iam::123456789012:role/perses-read", "externalIdSecret": "cw-external-id"}`,
			yamele: `
region: eu-west-3
roleArn: arn:aws:iam::123456789012:role/perses-read
externalIdSecret: cw-external-id
`,
			result: Config{Region: "eu-west-3", RoleARN: "arn:aws:iam::123456789012:role/perses-read", ExternalIDSecret: "cw-external-id"},
		},
		{
			title:  "role in another partition with a path",
			jason:  `{"region": "us-gov-west-1", "roleArn": "arn:aws-us-gov:iam::123456789012:role/monitoring/perses"}`,
			yamele: "region: us-gov-west-1\nroleArn: arn:aws-us-gov:iam::123456789012:role/monitoring/perses",
			result: Config{Region: "us-gov-west-1", RoleARN: "arn:aws-us-gov:iam::123456789012:role/monitoring/perses"},
		},
		{
			title:     "missing region",
			jason:     `{}`,
			yamele:    `roleArn: arn:aws:iam::123456789012:role/perses`,
			expectErr: true,
		},
		{
			title:     "invalid region",
			jason:     `{"region": "us_east_1"}`,
			yamele:    `region: us_east_1`,
			expectErr: true,
		},
		{
			title:     "user ARN instead of a role",
			jason:     `{"region": "us-east-1", "roleArn": "arn:aws:iam::123456789012:user/bob"}`,
			yamele:    "region: us-east-1\nroleArn: arn:aws:iam::123456789012:user/bob",
			expectErr: true,
		},
		{
			title:     "external ID without role",
			jason:     `{"region": "us-east-1", "externalIdSecret": "cw-external-id"}`,
			yamele:    "region: us-east-1\nexternalIdSecret: cw-external-id",
			expectErr: true,
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			jsonResult := Config{}
			jsonErr := json.Unmarshal([]byte(test.jason), &jsonResult)
			yamlResult := Config{}
			yamlErr := yaml.Unmarshal([]byte(test.yamele), &yamlResult)
			if test.expectErr {
				assert.Error(t, jsonErr)
				assert.Error(t, yamlErr)
				return
			}
			assert.NoError(t, jsonErr)
			assert.NoError(t, yamlErr)
			assert.Equal(t, test.result, jsonResult)
			assert.Equal(t, test.result, yamlResult)
		})
	}
}

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
	"errors"
	"fmt"
	"regexp"

	"github.com/perses/spec/go/datasource/proxy"
)

var (
	regionRegexp  = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-[0-9]{1,2}$`)
	roleARNRegexp = regexp.MustCompile(`^arn:aws(-[a-z]+)*:iam::[0-9]{12}:role/[A-Za-z0-9+=,.@_/-]{1,576}$`)
)

// Config is the configuration of the CloudWatch proxy.
// It never contains AWS credentials: the Perses server signs the requests with its own identity,
// optionally assuming the role defined here.
type Config struct {
	// region is the AWS region of the CloudWatch API, for example "us-east-1".
	Region string `json:"region" yaml:"region"`
	// roleArn is the ARN of the IAM role the Perses server assumes to query CloudWatch.
	// When it is empty, the server uses its own AWS identity directly.
	RoleARN string `json:"roleArn,omitempty" yaml:"roleArn,omitempty"`
	// externalIdSecret is the name of the secret holding the external ID used to assume the role.
	// The external ID is stored in the authorization credentials of the secret.
	// It can only be set together with roleArn.
	ExternalIDSecret string `json:"externalIdSecret,omitempty" yaml:"externalIdSecret,omitempty"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var tmp Config
	type plain Config
	if err := json.Unmarshal(data, (*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).Validate(); err != nil {
		return err
	}
	*c = tmp
	return nil
}

func (c *Config) UnmarshalYAML(unmarshal func(any) error) error {
	var tmp Config
	type plain Config
	if err := unmarshal((*plain)(&tmp)); err != nil {
		return err
	}
	if err := (&tmp).Validate(); err != nil {
		return err
	}
	*c = tmp
	return nil
}

// Validate checks the configuration. It is exported so the server can validate a configuration it builds itself.
func (c *Config) Validate() error {
	if !regionRegexp.MatchString(c.Region) {
		return fmt.Errorf("region %q is not a valid AWS region", c.Region)
	}
	if len(c.RoleARN) > 0 && !roleARNRegexp.MatchString(c.RoleARN) {
		return fmt.Errorf("roleArn %q is not a valid IAM role ARN", c.RoleARN)
	}
	if len(c.ExternalIDSecret) > 0 && len(c.RoleARN) == 0 {
		return errors.New("externalIdSecret can only be used with roleArn")
	}
	return nil
}

// Proxy is the CloudWatch proxy definition proposed in Perses.
// In case you are defining a datasource that will work with the Perses backend, then you will need to use this definition.
type Proxy proxy.Proxy[Config]

const (
	ProxyKindName = "cloudwatchproxy"
)

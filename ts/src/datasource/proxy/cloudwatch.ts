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

export interface CloudWatchProxy {
  kind: 'CloudWatchProxy';
  spec: CloudWatchProxySpec;
}

export interface CloudWatchProxySpec {
  // region is the AWS region of the CloudWatch API, for example "us-east-1".
  region: string;
  // roleArn is the ARN of the IAM role the Perses server assumes to query CloudWatch.
  // When it is empty, the server uses its own AWS identity directly.
  roleArn?: string;
  // externalIdSecret is the name of the secret holding the external ID used to assume the role.
  // The external ID is stored in the authorization credentials of the secret.
  // It can only be set together with roleArn.
  externalIdSecret?: string;
}

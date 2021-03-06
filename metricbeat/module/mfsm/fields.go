// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package mfsm

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "mfsm", asset.ModuleFieldsPri, AssetMfsm); err != nil {
		panic(err)
	}
}

// AssetMfsm returns asset data.
// This is the base64 encoded gzipped contents of module/mfsm.
func AssetMfsm() string {
	return "eJx8j0uOhCAARPecouLeC7CY3RwEh3JC5BfAdHv7DtJt1JB+O6jwqhixcJNwc3YCKKZYSgz1OAgg0VJlSkwsSgCa+S+ZWEzwEj8CwP4SLujVUgCzodVZ7skIrxwPd6VskRL/KazxfdMxXi0X00TljZ/DkfSUlfvwD926Rk9+H3Iew6dycf/0mTZn4fYISd+yL+WV3yZspeIVAAD//6CBbCE="
}

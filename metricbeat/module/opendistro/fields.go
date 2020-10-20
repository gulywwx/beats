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

package opendistro

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "opendistro", asset.ModuleFieldsPri, AssetOpendistro); err != nil {
		panic(err)
	}
}

// AssetOpendistro returns asset data.
// This is the base64 encoded gzipped contents of module/opendistro.
func AssetOpendistro() string {
	return "eJzEkE9OxyAQhfec4qX7XoCFO89hsDwNKQUC02hvb1r8g4juzO9bziTfezMzVh4aMTFYVyRHBYgTT43pazgpINPTFGo8UowCLMuSXRIXg8adAtBYsEW7eyrgydHboq/9jGA2dmknciRqPOe4p/fJwP7d1foWvxdhfihipHxuR9qT/pAPhpGV3wL6Qm0pvpotXS9oqZVWHi8x2273R4GT+yqsoT9+EKLl/z1gaL/h9W8BAAD///hruVU="
}

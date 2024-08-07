/*
Copyright 2021 The metaGraf Authors

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

package params

var (
	// Dryrun bool indicated if mg should do operations against a Kubernetes API.
	Dryrun bool

	// Output bool indicates if mg should output the generated objects.
	Output bool

	// Format value can be either json or yaml. Controls output format.
	Format string = "json"

	// Slice of strings to hold any labels added through --labels argument.
	Labels []string

	// Slice of strings to hold any annotations added through --annotations argument.
	Annotations []string
)

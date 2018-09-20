/*
Copyright 2018 The metaGraf Authors

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

package modules

import (
	"fmt"
	"metagraf/pkg/metagraf"
)

func InspectConfigMaps(mg *metagraf.MetaGraf) {
	for _,c := range mg.Spec.Config {
		fmt.Println(Name(mg),"ConfigMap",c.Name)
	}
}

func InspectSecrets(mg *metagraf.MetaGraf) {
	for _,r := range mg.Spec.Resources {
		if len(r.Secret) == 0 && len(r.User) > 0 {
			fmt.Println(Name(mg), "needs secret for user", r.User,"for resource", r.Name+".", "Secret name:", ResourceSecretName(&r))
		}

		if len(r.Secret) > 0 && r.SecretType == "cert" {
			fmt.Println(Name(mg), "needs cert(\""+r.Secret+"\") for", r.Name+".", "Secret name:", ResourceSecretName(&r))
		}
	}
}
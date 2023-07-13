/*
Copyright 2016 Skippbox, Ltd.
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

package event

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
)

// Event represent an event got from k8s api server
// Events from different endpoints need to be casted to KubewatchEvent
// before being able to be handled by handler
type Event struct {
	Namespace   string
	Kind        string
	ApiVersion  string
	Component   string
	Host        string
	Reason      string
	Status      string
	InfoMessage string
	Name        string
	Obj         runtime.Object
	OldObj      runtime.Object
}

// Message returns event message in standard format.
// included as a part of event package to enhance code reusability across handlers.
func (e *Event) Message() (msg string) {
	msg = fmt.Sprintf(
		"`%s` status for `%s` `%s` in namespace `%s`. Reason: `%s` Details: `%s`",
		e.Status,
		e.Kind,
		e.Name,
		e.Namespace,
		e.Reason,
		e.InfoMessage,
	)
	return msg
}

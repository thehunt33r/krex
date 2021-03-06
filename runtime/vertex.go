// Copyright © 2018 Kris Nova <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type Vertex struct {
	Name  string
	Above *Vertex
	//Down *Vertex
	Prompt         promptui.Select
	ListFunc       List
	PreviousOutput string
	Namespace      string
	Terminate      bool
}

func (v *Vertex) Select() (*Vertex, Action, error) {
	v2, err := v.ListFunc(v)
	if err != nil {
		return nil, ActionEmpty, err
	}
	if v2 == nil {

	}
	_, str, err := v2.Prompt.Run()
	if err != nil {
		return nil, ActionEmpty, err
	}
	if v2.Terminate {
		// 		Items: []string{"Edit", "Describe", "Logs", "Containers", "Shell Debug"},
		switch strings.ToLower(str) {
		case "edit":
			return nil, ActionEdit, nil
		case "describe":
			return nil, ActionDescribe, nil
		case "logs":
			return nil, ActionLogs, nil
		//case "containers":
		//	return nil, ActionContainers, nil
		case "shell debug":
			return nil, ActionShellDebug, nil
		default:
			return nil, ActionEmpty, fmt.Errorf("Invalid action: %s", str)
		}
	}
	v2.PreviousOutput = str
	return v2, ActionEmpty, nil
}

func (v *Vertex) RecursiveSelect() error {
	v2, action, err := v.Select()
	if err != nil {
		return err
	}
	if v2 == nil {
		err := action(&ActionParametes{
			PodName:   v.PreviousOutput,
			Namespace: v.Namespace,
		})
		if err != nil {
			return err
		}
		return nil
	}
	return v2.RecursiveSelect()
}

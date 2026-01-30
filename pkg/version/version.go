/**
 * Copyright 2025 Saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package version

import "fmt"

var (
	buildTime = ""
	gitTag    = ""
	gitHash   = ""
	version   = ""
)

func Print(service string) {
	fmt.Printf("%s\n", service)
	fmt.Printf("\tBuildTime:\t%s\n", buildTime)
	fmt.Printf("\tGitTag:\t\t%s\n", gitTag)
	fmt.Printf("\tGitHash:\t%s\n", gitHash)
	fmt.Printf("\tVersion:\t%s\n", version)
}

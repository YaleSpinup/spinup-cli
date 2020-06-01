/*
Copyright © 2020 Yale University

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import "github.com/YaleSpinup/spinup-cli/cmd"

var (
	// Version is the main version number
	Version = "0.0.0"

	// VersionPrerelease is a prerelease marker
	VersionPrerelease = ""

	// BuildStamp is the timestamp the binary was built, it should be set at buildtime with ldflags
	BuildStamp = ""

	// GitHash is the git sha of the built binary, it should be set at buildtime with ldflags
	GitHash = ""
)

func main() {
	cmd.Version = Version
	cmd.VersionPrerelease = VersionPrerelease
	cmd.BuildStamp = BuildStamp
	cmd.GitHash = GitHash

	cmd.Execute()
}

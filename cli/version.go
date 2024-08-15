// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package cli

var (
	version  = "0.0.0"
	buildNum = "unknown"
	commitID = "unknown"
)

type VersionInfo struct {
	Version     string
	BuildNumber string
	CommitID    string
}

var ProgramVersion = VersionInfo{
	Version:     version,
	BuildNumber: buildNum,
	CommitID:    commitID,
}

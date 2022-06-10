package versioninfo

import "fmt"

type VersionInfo struct {
	Version string `json:"version"`
	Time    string `json:"build-time"`
	Arch    string `json:"arch"`
	Build   string `json:"build"`
}

var version *VersionInfo

func Init(ver, buildNum, time, arch string) {
	version = &VersionInfo{
		Version: ver,
		Build:   buildNum,
		Time:    time,
		Arch:    arch}
}

func Get() *VersionInfo {
	if version != nil {
		return version
	}
	return &VersionInfo{}
}

func (v VersionInfo) String() string {
	return fmt.Sprintf("VERSION: %s\nBUILD: %s\nBUILD-TIME: %s\nARCH: %s\n",
		v.Version, v.Build, v.Time, v.Arch)
}

func (v VersionInfo) Output() {
	fmt.Println(v.String())
}

func GetUserAgent() string {
	version := Get()
	return fmt.Sprintf("PCTL/%s %s", version.Version, version.Arch)
}

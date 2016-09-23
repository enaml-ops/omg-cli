package utils

import "github.com/enaml-ops/enaml"

func GetVMTypeNames(vmTypes []enaml.VMType) []string {
	types := []string{}
	for _, vmType := range vmTypes {
		types = append(types, vmType.Name)
	}
	return types
}
func GetDiskTypeNames(diskTypes []enaml.DiskType) []string {
	types := []string{}
	for _, diskType := range diskTypes {
		types = append(types, diskType.Name)
	}
	return types
}

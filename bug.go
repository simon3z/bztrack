package main

const (
	STATUS_NEW = "NEW"
)

// Severity definitions
const (
	SEVERITY_UNDEFINED = "undefined"
	SEVERITY_LOW       = "low"
	SEVERITY_MEDIUM    = "medium"
	SEVERITY_HIGH      = "high"
)

// Priority definitions
const (
	PRIORITY_UNDEFINED = "undefined"
	PRIORITY_LOW       = "low"
	PRIORITY_MEDIUM    = "medium"
	PRIORITY_HIGH      = "high"
)

// BugInformation
type BugInformation struct {
	Product            string
	Component          string
	AssignedTo         string
	Status             string
	Summary            string
	InternalWhiteboard string
	Severity           string
	Priority           string
	TargetRelease      string
}

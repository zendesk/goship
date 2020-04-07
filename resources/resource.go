package resources

// Resource is an interface which all of the resources should be implementing
type Resource interface {
	ConnectIdentifier(bool, bool) string
	Name() string
	ResourceID() string
	GetTag(string) string
	GetZone() string
	RenderShortOutput() string
	RenderLongOutput() string
	SortKey() string
	String() string
}

// ResourceList is list of resources
type ResourceList []Resource

// Len returns length of resources list
func (r ResourceList) Len() int {
	return len(r)
}

// Less returns whether i is less than j
func (r ResourceList) Less(i, j int) bool {
	return r[i].SortKey() < r[j].SortKey()
}

// Swap swaps i and j
func (r ResourceList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

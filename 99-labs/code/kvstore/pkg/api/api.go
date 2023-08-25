package api

type VersionedValue struct {
	Value   string `json:"value"`
	Version int    `json:"version"`
}

type VersionedKeyValue struct {
	Key string `json:"key"`
	VersionedValue
}

package entity

type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

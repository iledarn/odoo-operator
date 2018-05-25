package v1alpha1

type ImageSpec struct {
	Registry string `json:"registry"`
	Name     string `json:"image"`
	Tag      string `json:"tag"`
}

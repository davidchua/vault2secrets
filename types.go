package main

type Placeholder struct {
	Message string
}

type VaultEvent struct {
	Type   string
	Object VaultObject `json:"object"`
}

type VaultObject struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Spec       CustomSecretSpec  `json:"spec"`
	Metadata   map[string]string `json:"metadata"`
}

type CustomSecretSpec struct {
	Url      string `json:"url"`
	Path     string `json:"path"`
	TokenRef string `json:"tokenRef"`
	Secret   string `json:"secret"`
}

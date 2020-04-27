package modelo

type Departamento struct {
	Nombre     string   `json:"nombre"`
	Municipios []string `json:"municipios"`
}

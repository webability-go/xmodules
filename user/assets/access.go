package assets

type AccessGroup struct {
	Key         string
	Name        string
	Description string
}

type Access struct {
	Key         string
	Name        string
	Group       string
	Description string
}

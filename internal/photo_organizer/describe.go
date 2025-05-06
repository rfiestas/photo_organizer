package photoorganizer

type Describer interface {
	Describe(imagePath string) (string, []string, error)
}

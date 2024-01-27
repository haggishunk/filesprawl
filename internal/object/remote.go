package object

type Remote struct {
	Hostname string
	Name     string // different hosts may have same remote with different name
}

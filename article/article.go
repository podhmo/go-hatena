package article

// Article :
type Article struct {
	Title Title
	Body  string
}

// Title
type Title struct {
	Title      string
	Categories []string
}

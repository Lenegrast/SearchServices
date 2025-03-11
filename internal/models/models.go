package models

type Services struct {
	Name     string
	Requests []string
}

type TFIDF struct {
	documents []map[string]int
	df        map[string]int
}

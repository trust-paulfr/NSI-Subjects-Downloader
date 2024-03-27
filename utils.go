package main

type File struct {
	ID string `json:"id"`
}

type Subject struct {
	ID       string `json:"id"`
	Libelle  string `json:"libelle"`
	Fichiers []File `json:"fichiers"`
}

type Response struct {
	Content []Subject `json:"content"`
}

func isPDF(content []byte) bool {
	return len(content) >= 4 && string(content[:5]) == "%PDF-"
}

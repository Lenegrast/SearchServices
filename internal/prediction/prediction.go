package prediction

import (
	"SearchServices/internal/data"
	"bufio"
	"fmt"
	"github.com/kljensen/snowball"
	"math"
	"os"
	"regexp"
	"strings"
)

var stopWords = map[string]struct{}{
	"и": {}, "в": {}, "не": {}, "на": {}, "с": {}, "по": {}, "для": {}, "Добрый": {},
}

type TFIDF struct {
	documents []map[string]int
	df        map[string]int
}

func ReadConsole() string {
	fmt.Println("Введите описание проблемы:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func cleanText(text string) []string {
	reg := regexp.MustCompile(`[^\p{L}\s]`)
	cleanText := reg.ReplaceAllString(strings.ToLower(text), "")
	words := strings.Fields(cleanText)
	var result []string
	for _, word := range words {
		if _, ok := stopWords[word]; ok {
			continue
		}
		stemmed, err := snowball.Stem(word, "russian", true)
		if err == nil {
			result = append(result, stemmed)
		}
	}
	return result
}

func (t *TFIDF) AddDocument(doc []string) {
	docMap := make(map[string]int)
	for _, term := range doc {
		docMap[term]++
	}
	t.documents = append(t.documents, docMap)
	for term := range docMap {
		t.df[term]++
	}
}

func (t *TFIDF) GetVector(doc []string) map[string]float64 {
	vector := make(map[string]float64)
	totalTerms := len(doc)
	if totalTerms == 0 {
		return vector
	}
	termCounts := make(map[string]int)
	for _, term := range doc {
		termCounts[term]++
	}
	for term, count := range termCounts {
		tf := float64(count) / float64(totalTerms)
		idf := math.Log((float64(len(t.documents)) + 1/(float64(t.df[term])+1)))
		// Исправленная формула
		vector[term] = tf * idf
	}
	return vector
}

func cosineSimilarity(a, b map[string]float64) float64 {
	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0
	for term, weightA := range a {
		magnitudeA += weightA * weightA
		if weightB, exists := b[term]; exists {
			dotProduct += weightA * weightB
		}
	}
	for _, weightB := range b {
		magnitudeB += weightB * weightB
	}
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}

func FinalResponce(personQuery string) string {
	services := data.AllData()
	tfidf := TFIDF{df: make(map[string]int)}
	for _, service := range services {
		combinedText := strings.Join(service.Requests, " ")
		processed := cleanText(combinedText)
		tfidf.AddDocument(processed)
	}

	serviceVectors := make([]map[string]float64, len(services))
	for i, service := range services {
		combinedText := strings.Join(service.Requests, " ")
		processed := cleanText(combinedText)
		serviceVectors[i] = tfidf.GetVector(processed)
	}

	query := personQuery
	processedQuery := cleanText(query)
	if len(processedQuery) == 0 {
		fmt.Println("Ошибка: запрос пуст после обработки.")
	}
	queryVector := tfidf.GetVector(processedQuery)

	bestScore := 0.0
	bestService := "Нет подходящих услуг"
	for i, vector := range serviceVectors {
		score := cosineSimilarity(queryVector, vector)
		if score > bestScore {
			bestScore = score
			bestService = services[i].Name
		}
	}

	return fmt.Sprintf("Рекомендуем: %s (схожесть: %.2f)\n", bestService, bestScore)
}

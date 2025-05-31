/*
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
*/

package prediction

import (
	"SearchServices/internal/data"
	"bufio"
	"fmt"
	"github.com/kljensen/snowball"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

var stopWords = map[string]struct{}{
	"и": {}, "в": {}, "не": {}, "на": {}, "с": {}, "по": {}, "для": {}, "добрый": {}, "день": {},
}

// TFIDF структура для вычисления TF-IDF векторов
type TFIDF struct {
	documents []map[string]int
	df        map[string]int
}

// LinearSVC упрощённая реализация линейного классификатора (One-vs-Rest)
type LinearSVC struct {
	weights      []map[string]float64 // Веса для каждого класса
	biases       []float64            // Смещения для каждого класса
	classes      []string             // Названия классов (услуги)
	learningRate float64              // Скорость обучения
	iterations   int                  // Количество итераций
	lambda       float64              // Параметр регуляризации
}

// ReadConsole читает ввод пользователя
func ReadConsole() string {
	fmt.Println("Введите описание проблемы:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// cleanText очищает и стеммирует текст
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

// AddDocument добавляет документ в TF-IDF модель
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

// GetVector вычисляет TF-IDF вектор для документа
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
		idf := math.Log(float64(len(t.documents)+1) / (float64(t.df[term]) + 1))
		vector[term] = tf * idf
	}
	// Нормализация вектора
	magnitude := 0.0
	for _, val := range vector {
		magnitude += val * val
	}
	if magnitude > 0 {
		magnitude = math.Sqrt(magnitude)
		for term := range vector {
			vector[term] /= magnitude
		}
	}
	return vector
}

// NewLinearSVC создаёт новый классификатор
func NewLinearSVC(classes []string) *LinearSVC {
	weights := make([]map[string]float64, len(classes))
	for i := range weights {
		weights[i] = make(map[string]float64)
	}
	return &LinearSVC{
		weights:      weights,
		biases:       make([]float64, len(classes)),
		classes:      classes,
		learningRate: 0.005,
		iterations:   500,
		lambda:       0.1,
	}
}

// Train обучает LinearSVC на данных (One-vs-Rest)
func (svc *LinearSVC) Train(X []map[string]float64, y []int) {
	for classIdx := range svc.classes {
		// Инициализация весов для текущего класса
		for _, vector := range X {
			for term := range vector {
				if _, exists := svc.weights[classIdx][term]; !exists {
					svc.weights[classIdx][term] = 0.001 * (float64(classIdx%2)*2 - 1)
				}
			}
		}

		// Обучение бинарного классификатора для текущего класса
		for iter := 0; iter < svc.iterations; iter++ {
			for i, vector := range X {
				// Метка: 1, если это текущий класс, иначе 0
				target := 0.0
				if y[i] == classIdx {
					target = 1.0
				}

				// Предсказание: w * x + b
				score := svc.biases[classIdx]
				for term, value := range vector {
					score += svc.weights[classIdx][term] * value
				}

				// Логистическая функция
				pred := 1 / (1 + math.Exp(-score))

				// Обновление весов и смещения с L2-регуляризацией
				gradient := pred - target
				for term, value := range vector {
					svc.weights[classIdx][term] -= svc.learningRate * (gradient*value + svc.lambda*svc.weights[classIdx][term])
				}
				svc.biases[classIdx] -= svc.learningRate * gradient
			}
		}
	}
}

// Predict возвращает вероятности для всех классов
func (svc *LinearSVC) Predict(vector map[string]float64) []float64 {
	scores := make([]float64, len(svc.classes))
	sumExp := 0.0

	// Вычисление сырых оценок для каждого класса
	for classIdx := range svc.classes {
		scores[classIdx] = svc.biases[classIdx]
		for term, value := range vector {
			if weight, exists := svc.weights[classIdx][term]; exists {
				scores[classIdx] += weight * value
			}
		}
		sumExp += math.Exp(scores[classIdx])
	}

	// Softmax для получения вероятностей
	probabilities := make([]float64, len(svc.classes))
	for classIdx, score := range scores {
		probabilities[classIdx] = math.Exp(score) / sumExp
	}

	return probabilities
}

// FinalResponce выполняет предсказание до 5 наиболее похожих услуг
func FinalResponce(personQuery string) string {
	services := data.AllData()
	tfidf := TFIDF{df: make(map[string]int)}

	// Подготовка данных
	X := make([]map[string]float64, 0)
	y := make([]int, 0)
	classes := make([]string, 0)

	for i, service := range services {
		combinedText := strings.Join(service.Requests, " ")
		processed := cleanText(combinedText)
		tfidf.AddDocument(processed)
		classes = append(classes, service.Name)
		X = append(X, tfidf.GetVector(processed))
		y = append(y, i)
	}

	// Обучение модели
	svc := NewLinearSVC(classes)
	svc.Train(X, y)

	// Обработка запроса пользователя
	query := personQuery
	processedQuery := cleanText(query)
	if len(processedQuery) == 0 {
		return "Ошибка: запрос пуст после обработки."
	}
	queryVector := tfidf.GetVector(processedQuery)

	// Предсказание
	probabilities := svc.Predict(queryVector)

	// Формирование списка до 5 наиболее похожих услуг
	type serviceProb struct {
		name string
		prob float64
	}
	serviceProbs := make([]serviceProb, len(classes))
	for i, prob := range probabilities {
		serviceProbs[i] = serviceProb{name: classes[i], prob: prob}
	}
	sort.Slice(serviceProbs, func(i, j int) bool {
		return serviceProbs[i].prob > serviceProbs[j].prob
	})

	// Ограничение до 5 услуг
	maxResults := 5
	if len(serviceProbs) < maxResults {
		maxResults = len(serviceProbs)
	}

	// Формирование вывода
	var result strings.Builder
	result.WriteString("Наиболее подходящие услуги:\n")
	for i := 0; i < maxResults; i++ {
		result.WriteString(fmt.Sprintf("%s\n", serviceProbs[i].name))
	}

	return result.String()
}

package convertors

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	ais "main/internal/ai_asker"
)

func calculate(birthDate time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()

	if today.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func parseDate(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, ".")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("неверный формат даты")
	}

	day := parts[0]
	month := parts[1]
	year := parts[2]
	normalDate := fmt.Sprintf("%s-%s-%s", year, month, day)

	birthDate, err := time.Parse("2006-01-02", normalDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка преобразования даты: %w", err)
	}
	return birthDate, nil
}

// итоговая функция
func CalculateAge(birthDateStr string) string {
	birthDate, err := parseDate(birthDateStr)
	if err != nil {
		return ""
	}

	age := calculate(birthDate)
	return fmt.Sprintf("%d", age)
}

// Функция для обработки строки и извлечения чисел
func extractNumbers(input string) string {
	var result strings.Builder
	for _, r := range input {
		if unicode.IsDigit(r) || r == '.' {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(result.String()), " ")
}

// Функция для подсчета суммы и количества чисел
func calculateSumAndCount(numbers []string) (float64, int) {
	var sum float64
	var count int
	for _, numStr := range numbers {
		num, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			sum += num
			count++
		}
	}
	return sum, count
}

// Функция для вычисления среднего значения
func calculateAverage(sum float64, count int) float64 {
	if count > 0 {
		return sum / float64(count)
	}
	return 0
}

// Основная функция для нахождения среднего значения
func Average(input string) string {
	trimmed := extractNumbers(input)
	numbers := strings.Fields(trimmed)
	sum, count := calculateSumAndCount(numbers)
	average := calculateAverage(sum, count)
	averageStr := fmt.Sprintf("%.2f", average)
	return averageStr
}

func CalcItems(input string) string {

	return "0"
}

func AskAI(input string) string {

	authorizationKey := "OGYwMTljNzYtYzEyMy00MjE4LWJmY2UtZTY2ZWE1ZGRlM2E4OmVmNjg4Zjk0LWE5ZGYtNDdkMS1hYzI2LTI5NWQyOGRlMzZlNA=="

	accessToken, err := ais.GetAccessToken(authorizationKey)
	if err != nil {
		log.Println("Error getting access token:", err)
		return "I am Groot"
	}
	// log.Println("Access Token:", accessToken)

	ans, err := ais.SendNeuralNetRequest(accessToken, input)
	if err != nil {
		log.Printf("Ошибка при запросе k AI: %v", err)
		return "I am Groooot"
	}

	return ans
}

func Insert(input string) string {
	return input
}

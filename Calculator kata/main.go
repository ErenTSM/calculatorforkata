package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var romanNumerals = map[rune]int{
	'I': 1,
	'V': 5,
	'X': 10,
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите выражение (например, 'X + V' или '2 - 3'):")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return
	}
	input = strings.TrimSpace(input)

	result, err := calculate(input)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Результат: %s\n", result)
	}
}

func calculate(expression string) (string, error) {
	parts := strings.Fields(expression)
	if len(parts) != 3 {
		return "", fmt.Errorf("формат выражения должен быть 'число оператор число'")
	}

	// Проверяем, не смешаны ли римские и арабские числа
	if isRoman(parts[0]) != isRoman(parts[2]) {
		return "", fmt.Errorf("нельзя смешивать римские и арабские цифры в одном выражении")
	}

	operand1, err := convertToArabic(parts[0])
	if err != nil {
		return "", err
	}

	operand2, err := convertToArabic(parts[2])
	if err != nil {
		return "", err
	}

	operator := parts[1]

	result := 0
	switch operator {
	case "+":
		result = operand1 + operand2
	case "-":
		result = operand1 - operand2
	case "*":
		result = operand1 * operand2
	case "/":
		if operand2 == 0 {
			return "", fmt.Errorf("деление на ноль")
		}
		result = operand1 / operand2
	default:
		return "", fmt.Errorf("недопустимый оператор: %s", operator)
	}

	// Если оба операнда были римскими числами
	if isRoman(parts[0]) && isRoman(parts[2]) {
		if result <= 0 {
			return "", fmt.Errorf("результат выражения с римскими цифрами не может быть меньше I")
		}
		return convertToRoman(result), nil
	}
	return strconv.Itoa(result), nil
}

func isRoman(s string) bool {
	match, _ := regexp.MatchString("^[IVX]+$", s)
	return match
}

func convertToArabic(num string) (int, error) {
	if isRoman(num) {
		return romanToArabic(num), nil
	}
	value, err := strconv.Atoi(num)
	if err != nil || value < 1 || value > 10 {
		return 0, fmt.Errorf("арабские числа должны быть от 1 до 10")
	}
	return value, nil
}

func romanToArabic(roman string) int {
	total := 0
	lastValue := 0

	for i := len(roman) - 1; i >= 0; i-- {
		value := romanNumerals[rune(roman[i])]
		if value < lastValue {
			total -= value
		} else {
			total += value
		}
		lastValue = value
	}

	return total
}

func convertToRoman(arabic int) string {
	type numeral struct {
		Value  int
		Symbol string
	}

	numerals := []numeral{
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	var result strings.Builder
	for _, numeral := range numerals {
		for arabic >= numeral.Value {
			arabic -= numeral.Value
			result.WriteString(numeral.Symbol)
		}
	}

	return result.String()
}

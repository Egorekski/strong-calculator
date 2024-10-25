package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "**":
		return 3
	}
	return 0
}

func applyOp(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("invalid expression") // Изменено на "invalid expression"
		}
		return a / b, nil
	case "**":
		return math.Pow(a, b), nil
	}
	return 0, errors.New("invalid expression")
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
	i := 0
	for i < len(expression) {
		if expression[i] == ' ' {
			i++
			continue
		}

		if unicode.IsDigit(rune(expression[i])) || expression[i] == '.' {
			start := i
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				i++
			}
			tokens = append(tokens, expression[start:i])
		} else if expression[i] == '(' || expression[i] == ')' {
			tokens = append(tokens, string(expression[i]))
			i++
		} else if expression[i] == '*' && i+1 < len(expression) && expression[i+1] == '*' {
			tokens = append(tokens, "**")
			i += 2
		} else {
			tokens = append(tokens, string(expression[i]))
			i++
		}
	}
	return tokens, nil
}

func validateTokens(tokens []string) error {
	if len(tokens) == 0 {
		return errors.New("empty expression")
	}

	brackets := 0
	var lastToken string

	for i, token := range tokens {
		if token == "(" {
			brackets++
		} else if token == ")" {
			brackets--
			if brackets < 0 {
				return errors.New("invalid expression")
			}
		} else if precedence(token) > 0 {
			if i == 0 || lastToken == "(" || lastToken == "" {
				return errors.New("invalid expression")
			}

			if i == len(tokens)-1 {
				return errors.New("invalid expression")
			}

			if i+1 < len(tokens) && precedence(tokens[i+1]) > 0 {
				return errors.New("invalid expression")
			}
		}
		lastToken = token
	}

	if brackets != 0 {
		return errors.New("invalid expression")
	}

	return nil
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	if err := validateTokens(tokens); err != nil {
		return 0, err
	}

	var values []float64
	var ops []string

	i := 0
	for i < len(tokens) {
		token := tokens[i]

		if num, err := strconv.ParseFloat(token, 64); err == nil {
			values = append(values, num)
		} else if token == "(" {
			ops = append(ops, token)
		} else if token == ")" {
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				val2 := values[len(values)-1]
				values = values[:len(values)-1]

				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				result, err := applyOp(val1, val2, op)
				if err != nil {
					return 0, errors.New("invalid expression")
				}

				values = append(values, result)
			}

			if len(ops) == 0 {
				return 0, errors.New("invalid expression")
			}
			ops = ops[:len(ops)-1]
		} else {

			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(token) {
				val2 := values[len(values)-1]
				values = values[:len(values)-1]

				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				result, err := applyOp(val1, val2, op)
				if err != nil {
					return 0, errors.New("invalid expression")
				}

				values = append(values, result)
			}
			ops = append(ops, token)
		}
		i++
	}

	for len(ops) > 0 {
		val2 := values[len(values)-1]
		values = values[:len(values)-1]

		val1 := values[len(values)-1]
		values = values[:len(values)-1]

		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		result, err := applyOp(val1, val2, op)
		if err != nil {
			return 0, errors.New("invalid expression")
		}

		values = append(values, result)
	}

	if len(values) != 1 {
		return 0, errors.New("invalid expression")
	}

	return values[0], nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	expression, _ := reader.ReadString('\n')
	expression = strings.TrimSpace(expression)

	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(result)
	}
}

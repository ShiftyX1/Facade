package conditions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/ShiftyX1/Facade/internal/config"
	"github.com/ShiftyX1/Facade/internal/template"
)

type Evaluator struct {
	templateEngine *template.Engine
}

func New(templateEngine *template.Engine) *Evaluator {
	return &Evaluator{
		templateEngine: templateEngine,
	}
}

func (e *Evaluator) EvaluateConditions(conditions []config.RouteCondition, ctx *template.Context) (*config.RouteCondition, error) {
	for _, condition := range conditions {
		if condition.Default {
			return &condition, nil
		}

		if condition.If != "" {
			match, err := e.evaluateCondition(condition.If, ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate condition: %w", err)
			}
			if match {
				return &condition, nil
			}
		}
	}

	for _, condition := range conditions {
		if condition.Default {
			return &condition, nil
		}
	}

	if len(conditions) > 0 {
		return &conditions[0], nil
	}

	return nil, fmt.Errorf("no matching condition found")
}

func (e *Evaluator) evaluateCondition(conditionStr string, ctx *template.Context) (bool, error) {
	rendered, err := e.templateEngine.Render(conditionStr, ctx)
	if err != nil {
		return false, fmt.Errorf("failed to render condition template: %w", err)
	}

	return e.parseCondition(rendered)
}

func (e *Evaluator) parseCondition(condition string) (bool, error) {
	condition = strings.TrimSpace(condition)

	if strings.Contains(condition, "==") {
		return e.evaluateComparison(condition, "==")
	}
	if strings.Contains(condition, "!=") {
		return e.evaluateComparison(condition, "!=")
	}
	if strings.Contains(condition, ">=") {
		return e.evaluateComparison(condition, ">=")
	}
	if strings.Contains(condition, "<=") {
		return e.evaluateComparison(condition, "<=")
	}
	if strings.Contains(condition, ">") {
		return e.evaluateComparison(condition, ">")
	}
	if strings.Contains(condition, "<") {
		return e.evaluateComparison(condition, "<")
	}

	return e.isTruthy(condition), nil
}

func (e *Evaluator) evaluateComparison(condition, operator string) (bool, error) {
	parts := strings.SplitN(condition, operator, 2)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid comparison: %s", condition)
	}

	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	left = strings.Trim(left, `"'`)
	right = strings.Trim(right, `"'`)

	switch operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case ">", ">=", "<", "<=":
		return e.evaluateNumericComparison(left, right, operator)
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

func (e *Evaluator) evaluateNumericComparison(left, right, operator string) (bool, error) {
	leftNum, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return false, fmt.Errorf("left operand is not a number: %s", left)
	}

	rightNum, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return false, fmt.Errorf("right operand is not a number: %s", right)
	}

	switch operator {
	case ">":
		return leftNum > rightNum, nil
	case ">=":
		return leftNum >= rightNum, nil
	case "<":
		return leftNum < rightNum, nil
	case "<=":
		return leftNum <= rightNum, nil
	default:
		return false, fmt.Errorf("unknown numeric operator: %s", operator)
	}
}

func (e *Evaluator) isTruthy(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	
	switch value {
	case "", "0", "false", "null", "undefined":
		return false
	default:
		return true
	}
}

func ExtractFromJSON(jsonStr, path string) string {
	result := gjson.Get(jsonStr, path)
	return result.String()
}

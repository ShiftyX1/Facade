package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"text/template"
	"time"

	"github.com/go-faker/faker/v4"
)


type Engine struct {
	tmpl *template.Template
}


type Context struct {
	PathParams  map[string]string      `json:"path_params"`
	QueryParams map[string]string      `json:"query_params"`
	Body        map[string]interface{} `json:"body"`
}


func New() *Engine {
	tmpl := template.New("facade").Funcs(template.FuncMap{
		
		"randomInt":    randomInt,
		"randomFloat":  randomFloat,
		"randomString": randomString,
		"randomBool":   randomBool,

		
		"fakerName":    func() string { return faker.Name() },
		"fakerEmail":   func() string { return faker.Email() },
		"fakerPhone":   func() string { return faker.Phonenumber() },
		"fakerAddress": func() string { return faker.GetRealAddress().Address },
		"fakerCompany": func() string { return faker.Word() },
		"fakerWord":    func() string { return faker.Word() },
		"fakerUUID":    func() string { return faker.UUIDHyphenated() },

		
		"now":        func() time.Time { return time.Now() },
		"formatTime": formatTime,

		
		"json": jsonEncode,
		"add":  add,
		"sub":  sub,
		"mul":  mul,
		"div":  div,
	})

	return &Engine{tmpl: tmpl}
}


func (e *Engine) Render(templateStr string, ctx *Context) (string, error) {
	tmpl, err := e.tmpl.Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}



func randomInt(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min+1) + min
}

func randomFloat(min, max float64) float64 {
	if min >= max {
		return min
	}
	return min + rand.Float64()*(max-min)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func formatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

func jsonEncode(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func add(a, b interface{}) (interface{}, error) {
	return mathOp(a, b, "+")
}

func sub(a, b interface{}) (interface{}, error) {
	return mathOp(a, b, "-")
}

func mul(a, b interface{}) (interface{}, error) {
	return mathOp(a, b, "*")
}

func div(a, b interface{}) (interface{}, error) {
	return mathOp(a, b, "/")
}

func mathOp(a, b interface{}, op string) (interface{}, error) {
	aFloat, err := toFloat64(a)
	if err != nil {
		return nil, err
	}
	bFloat, err := toFloat64(b)
	if err != nil {
		return nil, err
	}

	switch op {
	case "+":
		return aFloat + bFloat, nil
	case "-":
		return aFloat - bFloat, nil
	case "*":
		return aFloat * bFloat, nil
	case "/":
		if bFloat == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return aFloat / bFloat, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", op)
	}
}

func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

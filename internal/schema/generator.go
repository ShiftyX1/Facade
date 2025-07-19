package schema

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/ShiftyX1/Facade/internal/config"
)


type Generator struct {
	rand *rand.Rand
}


func New() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}


func (g *Generator) Generate(schema *config.SchemaDefinition) (interface{}, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	switch schema.Type {
	case "array":
		return g.generateArray(schema)
	case "object":
		return g.generateObject(schema)
	case "string":
		return g.generateString(schema)
	case "integer", "int":
		return g.generateInteger(schema)
	case "number", "float":
		return g.generateFloat(schema)
	case "boolean", "bool":
		return g.generateBoolean()
	default:
		return nil, fmt.Errorf("unsupported schema type: %s", schema.Type)
	}
}

func (g *Generator) generateArray(schema *config.SchemaDefinition) ([]interface{}, error) {
	count := schema.Count
	if count == 0 {
		count = 1 
	}

	var result []interface{}
	for i := 0; i < count; i++ {
		item, err := g.Generate(schema.Items)
		if err != nil {
			return nil, fmt.Errorf("failed to generate array item: %w", err)
		}
		result = append(result, item)
	}

	return result, nil
}

func (g *Generator) generateObject(schema *config.SchemaDefinition) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, prop := range schema.Properties {
		propSchema := &config.SchemaDefinition{
			Type:  prop.Type,
			Min:   prop.Min,
			Max:   prop.Max,
			Faker: prop.Faker,
		}

		value, err := g.Generate(propSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to generate property %s: %w", key, err)
		}
		result[key] = value
	}

	return result, nil
}

func (g *Generator) generateString(schema *config.SchemaDefinition) (string, error) {
	if schema.Faker != "" {
		return g.generateFakerString(schema.Faker)
	}

	
	length := 10
	if schema.Min != nil && schema.Max != nil {
		length = g.rand.Intn(*schema.Max-*schema.Min+1) + *schema.Min
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[g.rand.Intn(len(charset))]
	}

	return string(result), nil
}

func (g *Generator) generateInteger(schema *config.SchemaDefinition) (int, error) {
	min := 0
	max := 100

	if schema.Min != nil {
		min = *schema.Min
	}
	if schema.Max != nil {
		max = *schema.Max
	}

	if schema.Faker == "random_int" {
		return g.rand.Intn(max-min+1) + min, nil
	}

	return g.rand.Intn(max-min+1) + min, nil
}

func (g *Generator) generateFloat(schema *config.SchemaDefinition) (float64, error) {
	min := 0.0
	max := 100.0

	if schema.Min != nil {
		min = float64(*schema.Min)
	}
	if schema.Max != nil {
		max = float64(*schema.Max)
	}

	return min + g.rand.Float64()*(max-min), nil
}

func (g *Generator) generateBoolean() (bool, error) {
	return g.rand.Intn(2) == 1, nil
}

func (g *Generator) generateFakerString(fakerType string) (string, error) {
	switch fakerType {
	case "name":
		return faker.Name(), nil
	case "email":
		return faker.Email(), nil
	case "phone":
		return faker.Phonenumber(), nil
	case "address":
		return faker.GetRealAddress().Address, nil
	case "company":
		return faker.Word(), nil 
	case "word":
		return faker.Word(), nil
	case "sentence":
		return faker.Sentence(), nil
	case "paragraph":
		return faker.Paragraph(), nil
	case "url":
		return faker.URL(), nil
	case "uuid":
		return faker.UUIDHyphenated(), nil
	case "date":
		return faker.Date(), nil
	case "timestamp":
		return faker.Timestamp(), nil
	default:
		return "", fmt.Errorf("unsupported faker type: %s", fakerType)
	}
}

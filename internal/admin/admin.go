package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ShiftyX1/Facade/internal/config"
	"github.com/ShiftyX1/Facade/internal/state"
	"github.com/gorilla/mux"
)

type Admin struct {
	config       *config.Config
	stateManager *state.Manager
	stats        *Statistics
}

type Statistics struct {
	StartTime    time.Time              `json:"start_time"`
	Requests     map[string]int         `json:"requests"`
	Responses    map[string]int         `json:"responses"`
	TotalRequests int                   `json:"total_requests"`
	Errors       int                    `json:"errors"`
	LastActivity time.Time              `json:"last_activity"`
}

type RouteInfo struct {
	Path        string                 `json:"path"`
	Method      string                 `json:"method"`
	Status      int                    `json:"status"`
	Description string                 `json:"description"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Schema      *config.SchemaDefinition `json:"schema,omitempty"`
	Example     string                 `json:"example,omitempty"`
}

func New(cfg *config.Config, sm *state.Manager) *Admin {
	return &Admin{
		config:       cfg,
		stateManager: sm,
		stats: &Statistics{
			StartTime: time.Now(),
			Requests:  make(map[string]int),
			Responses: make(map[string]int),
		},
	}
}

func (a *Admin) SetupRoutes(router *mux.Router) {
	// Admin dashboard
	router.HandleFunc("/_admin", a.dashboardHandler).Methods("GET")
	router.HandleFunc("/_admin/", a.dashboardHandler).Methods("GET")

	// API for admin
	router.HandleFunc("/_admin/api/routes", a.getRoutesHandler).Methods("GET")
	router.HandleFunc("/_admin/api/stats", a.getStatsHandler).Methods("GET")
	router.HandleFunc("/_admin/api/state", a.getStateHandler).Methods("GET")
	router.HandleFunc("/_admin/api/test", a.testEndpointHandler).Methods("POST")
	router.HandleFunc("/_admin/api/openapi", a.getOpenAPIHandler).Methods("GET")

	// Static files (CSS, JS)
	router.HandleFunc("/_admin/static/{file}", a.staticHandler).Methods("GET")
}

func (a *Admin) TrackRequest(method, path string, status int) {
	key := fmt.Sprintf("%s %s", method, path)
	a.stats.Requests[key]++
	a.stats.Responses[fmt.Sprintf("%d", status)]++
	a.stats.TotalRequests++
	a.stats.LastActivity = time.Now()
	
	if status >= 400 {
		a.stats.Errors++
	}
}
func (a *Admin) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := a.generateDashboardHTML()
	if _, err := w.Write([]byte(html)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (a *Admin) getRoutesHandler(w http.ResponseWriter, r *http.Request) {
	routes := make([]RouteInfo, 0, len(a.config.Routes))
	
	for _, route := range a.config.Routes {
		routeInfo := RouteInfo{
			Path:   route.Path,
			Method: route.Method,
			Status: route.Response.Status,
			Headers: route.Response.Headers,
			Schema: route.Response.Schema,
		}

		// Generate description
		if route.Response.Body != "" {
			routeInfo.Description = "Returns static response"
			routeInfo.Example = route.Response.Body
		} else if route.Response.Schema != nil {
			routeInfo.Description = "Returns generated data based on schema"
		}
		
		routes = append(routes, routeInfo)
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(routes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *Admin) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(a.stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *Admin) getStateHandler(w http.ResponseWriter, r *http.Request) {
	state := a.stateManager.GetAll()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"state": state,
		"keys":  a.stateManager.Keys(),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *Admin) testEndpointHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Method string            `json:"method"`
		Path   string            `json:"path"`
		Headers map[string]string `json:"headers"`
		Body   string            `json:"body"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// TODO: Here we can add logic to test endpoints
	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Test request to %s %s would be processed", request.Method, request.Path),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *Admin) staticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]

	// Basic static files
	switch file {
	case "style.css":
		w.Header().Set("Content-Type", "text/css")
		if _, err := w.Write([]byte(a.getCSS())); err != nil {
			http.Error(w, "Failed to write CSS", http.StatusInternalServerError)
		}
	case "script.js":
		w.Header().Set("Content-Type", "application/javascript")
		if _, err := w.Write([]byte(a.getJavaScript())); err != nil {
			http.Error(w, "Failed to write JavaScript", http.StatusInternalServerError)
		}
	default:
		http.NotFound(w, r)
	}
}
func (a *Admin) getOpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	openapi := a.generateOpenAPISpec()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openapi); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (a *Admin) generateOpenAPISpec() map[string]interface{} {
	paths := make(map[string]map[string]interface{})
	
	for _, route := range a.config.Routes {
		if paths[route.Path] == nil {
			paths[route.Path] = make(map[string]interface{})
		}
		
		operation := map[string]interface{}{
			"summary": fmt.Sprintf("%s %s", route.Method, route.Path),
			"responses": map[string]interface{}{
				fmt.Sprintf("%d", route.Response.Status): map[string]interface{}{
					"description": "Success response",
				},
			},
		}
		
		if route.Response.Schema != nil {
			operation["responses"].(map[string]interface{})[fmt.Sprintf("%d", route.Response.Status)].(map[string]interface{})["content"] = map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": a.convertSchemaToOpenAPI(route.Response.Schema),
				},
			}
		}
		
		paths[route.Path][strings.ToLower(route.Method)] = operation
	}
	
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":   "Facade Mock API",
			"version": "1.0.0",
			"description": "Generated OpenAPI specification for Facade mock server",
		},
		"servers": []map[string]interface{}{
			{
				"url": fmt.Sprintf("http://%s:%d", a.config.Server.Host, a.config.Server.Port),
				"description": "Mock server",
			},
		},
		"paths": paths,
	}
}

func (a *Admin) convertSchemaToOpenAPI(schema *config.SchemaDefinition) map[string]interface{} {
	result := map[string]interface{}{
		"type": schema.Type,
	}
	
	if schema.Properties != nil {
		properties := make(map[string]interface{})
		for name, prop := range schema.Properties {
			properties[name] = map[string]interface{}{
				"type": prop.Type,
			}
			if prop.Faker != "" {
				properties[name].(map[string]interface{})["description"] = fmt.Sprintf("Generated using faker: %s", prop.Faker)
			}
		}
		result["properties"] = properties
	}
	
	if schema.Items != nil {
		result["items"] = a.convertSchemaToOpenAPI(schema.Items)
	}
	
	return result
}
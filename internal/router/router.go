package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/ShiftyX1/Facade/internal/conditions"
	"github.com/ShiftyX1/Facade/internal/config"
	"github.com/ShiftyX1/Facade/internal/schema"
	"github.com/ShiftyX1/Facade/internal/state"
	"github.com/ShiftyX1/Facade/internal/template"
	"github.com/gorilla/mux"
)

type Router struct {
	mux             *mux.Router
	config          *config.Config
	templateEngine  *template.Engine
	schemaGenerator *schema.Generator
	stateManager    *state.Manager
	conditionsEval  *conditions.Evaluator
}

func New(cfg *config.Config) *Router {
	templateEngine := template.New()

	return &Router{
		mux:             mux.NewRouter(),
		config:          cfg,
		templateEngine:  templateEngine,
		schemaGenerator: schema.New(),
		stateManager:    state.New(cfg.Global.LogRequests),
		conditionsEval:  conditions.New(templateEngine),
	}
}

func (r *Router) Setup() error {
	for _, route := range r.config.Routes {
		if err := r.addRoute(route); err != nil {
			return fmt.Errorf("failed to add route %s %s: %w", route.Method, route.Path, err)
		}
	}

	return nil
}

func (r *Router) Handler() http.Handler {
	return r.mux
}

func (r *Router) addRoute(route config.Route) error {
	handler := r.createHandler(route)

	r.mux.HandleFunc(route.Path, handler).Methods(route.Method)

	slog.Debug("Route registered", "method", route.Method, "path", route.Path)
	return nil
}

func (r *Router) createHandler(route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx, err := r.extractRequestContext(req)
		if err != nil {
			slog.Error("Failed to extract request context", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var responseConfig *config.Response
		var selectedCondition *config.RouteCondition

		if len(route.Response.Conditions) > 0 {
			selectedCondition, err = r.conditionsEval.EvaluateConditions(route.Response.Conditions, ctx)
			if err != nil {
				slog.Error("Failed to evaluate conditions", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if selectedCondition != nil {
			responseConfig = &config.Response{
				Status:  selectedCondition.Status,
				Body:    selectedCondition.Body,
				Schema:  selectedCondition.Schema,
				Headers: route.Response.Headers,
			}
			if responseConfig.Status == 0 {
				responseConfig.Status = route.Response.Status
			}
		} else {
			responseConfig = &route.Response
		}

		responseBody, err := r.generateResponse(responseConfig, ctx)
		if err != nil {
			slog.Error("Failed to generate response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		r.setResponseHeaders(w, responseConfig.Headers)

		if req.Method == "POST" || req.Method == "PUT" {
			r.saveToState(route.Path, ctx.Body)
		}

		w.WriteHeader(responseConfig.Status)
		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}
}

func (r *Router) extractRequestContext(req *http.Request) (*template.Context, error) {
	ctx := &template.Context{
		PathParams:  make(map[string]string),
		QueryParams: make(map[string]string),
		Body:        make(map[string]interface{}),
	}

	vars := mux.Vars(req)
	for key, value := range vars {
		ctx.PathParams[key] = value
	}

	for key, values := range req.URL.Query() {
		if len(values) > 0 {
			ctx.QueryParams[key] = values[0]
		}
	}

	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}

		if len(bodyBytes) > 0 {
			var bodyData interface{}
			if err := json.Unmarshal(bodyBytes, &bodyData); err == nil {
				if bodyMap, ok := bodyData.(map[string]interface{}); ok {
					ctx.Body = bodyMap
				}
			}
		}
	}

	return ctx, nil
}

func (r *Router) generateResponse(response *config.Response, ctx *template.Context) (string, error) {

	if response.Schema != nil {
		data, err := r.schemaGenerator.Generate(response.Schema)
		if err != nil {
			return "", fmt.Errorf("failed to generate schema data: %w", err)
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("failed to marshal schema data: %w", err)
		}

		return string(jsonData), nil
	}

	if response.Body != "" {
		rendered, err := r.templateEngine.Render(response.Body, ctx)
		if err != nil {
			return "", fmt.Errorf("failed to render template: %w", err)
		}
		return rendered, nil
	}

	return "", nil
}

func (r *Router) setResponseHeaders(w http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
}

func (r *Router) saveToState(path string, body map[string]interface{}) {
	if len(body) == 0 {
		return
	}

	key := r.generateStateKey(path, body)
	r.stateManager.Set(key, body)
}

func (r *Router) generateStateKey(path string, body map[string]interface{}) string {

	if id, exists := body["id"]; exists {
		return fmt.Sprintf("%s:%v", path, id)
	}

	return path
}

func (r *Router) GetStateManager() *state.Manager {
	return r.stateManager
}

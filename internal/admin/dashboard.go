package admin

import (
    "fmt"
)

func (a *Admin) generateDashboardHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Facade Admin Dashboard</title>
    <link rel="stylesheet" href="/_admin/static/style.css">
</head>
<body>
    <div class="container">
        <header class="header">
            <h1>🎭 Facade Admin Dashboard</h1>
            <div class="server-info">
                <span>Server: ` + a.config.Server.Host + `:` + fmt.Sprintf("%d", a.config.Server.Port) + `</span>
                <span class="status-indicator online">● Online</span>
            </div>
        </header>

        <nav class="nav-tabs">
            <button class="tab-button active" onclick="showTab('overview')">Overview</button>
            <button class="tab-button" onclick="showTab('routes')">API Routes</button>
            <button class="tab-button" onclick="showTab('swagger')">Swagger UI</button>
            <button class="tab-button" onclick="showTab('testing')">Testing</button>
            <button class="tab-button" onclick="showTab('state')">State</button>
        </nav>

        <!-- Overview Tab -->
        <div id="overview" class="tab-content active">
            <div class="stats-grid">
                <div class="stat-card">
                    <h3>Total Requests</h3>
                    <div class="stat-value" id="total-requests">Loading...</div>
                </div>
                <div class="stat-card">
                    <h3>Active Routes</h3>
                    <div class="stat-value">` + fmt.Sprintf("%d", len(a.config.Routes)) + `</div>
                </div>
                <div class="stat-card">
                    <h3>Errors</h3>
                    <div class="stat-value error" id="error-count">Loading...</div>
                </div>
                <div class="stat-card">
                    <h3>Uptime</h3>
                    <div class="stat-value" id="uptime">Loading...</div>
                </div>
            </div>
            
            <div class="charts-section">
                <h3>Request Statistics</h3>
                <div id="request-chart" class="chart-container">
                    <p>Request statistics will be displayed here</p>
                </div>
            </div>
        </div>

        <!-- Routes Tab -->
        <div id="routes" class="tab-content">
            <div class="routes-header">
                <h3>API Routes</h3>
                <button onclick="refreshRoutes()" class="btn btn-primary">Refresh</button>
            </div>
            <div id="routes-list" class="routes-container">
                Loading routes...
            </div>
        </div>

        <!-- Swagger Tab -->
        <div id="swagger" class="tab-content">
            <div class="swagger-container">
                <h3>API Documentation</h3>
                <div class="swagger-actions">
                    <button onclick="downloadOpenAPI()" class="btn btn-secondary">Download OpenAPI Spec</button>
                    <button onclick="copyOpenAPIUrl()" class="btn btn-secondary">Copy OpenAPI URL</button>
                </div>
                <div id="swagger-ui" class="swagger-ui-container">
                    <div class="openapi-preview" id="openapi-preview">
                        Loading OpenAPI specification...
                    </div>
                </div>
            </div>
        </div>

        <!-- Testing Tab -->
        <div id="testing" class="tab-content">
            <div class="testing-container">
                <h3>API Testing</h3>
                <form id="test-form" class="test-form">
                    <div class="form-group">
                        <label>Method:</label>
                        <select id="test-method">
                            <option>GET</option>
                            <option>POST</option>
                            <option>PUT</option>
                            <option>DELETE</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>URL:</label>
                        <input type="text" id="test-url" placeholder="/api/users" value="/">
                    </div>
                    <div class="form-group">
                        <label>Headers (JSON):</label>
                        <textarea id="test-headers" placeholder='{"Content-Type": "application/json"}'>{}</textarea>
                    </div>
                    <div class="form-group">
                        <label>Body:</label>
                        <textarea id="test-body" placeholder='{"name": "test"}'></textarea>
                    </div>
                    <button type="submit" class="btn btn-primary">Send Request</button>
                </form>
                <div id="test-response" class="response-container" style="display: none;">
                    <h4>Response:</h4>
                    <pre id="response-content"></pre>
                </div>
            </div>
        </div>

        <!-- State Tab -->
        <div id="state" class="tab-content">
            <div class="state-container">
                <h3>Server State</h3>
                <div class="state-actions">
                    <button onclick="refreshState()" class="btn btn-primary">Refresh</button>
                    <button onclick="clearState()" class="btn btn-danger">Clear State</button>
                </div>
                <div id="state-content" class="state-content">
                    Loading state...
                </div>
            </div>
        </div>
    </div>

    <script src="/_admin/static/script.js"></script>
</body>
</html>`
}
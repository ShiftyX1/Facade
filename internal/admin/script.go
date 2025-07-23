package admin

func (a *Admin) getJavaScript() string {
	return `
// Tab management
function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Remove active from all buttons
    document.querySelectorAll('.tab-button').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(tabName).classList.add('active');
    
    // Add active to clicked button
    event.target.classList.add('active');
    
    // Load data for specific tabs
    if (tabName === 'routes') {
        loadRoutes();
    } else if (tabName === 'swagger') {
        loadOpenAPI();
    } else if (tabName === 'state') {
        loadState();
    } else if (tabName === 'overview') {
        loadStats();
    }
}

// Load statistics
function loadStats() {
    fetch('/_admin/api/stats')
        .then(response => response.json())
        .then(data => {
            document.getElementById('total-requests').textContent = data.total_requests || 0;
            document.getElementById('error-count').textContent = data.errors || 0;
            
            // Calculate uptime
            const startTime = new Date(data.start_time);
            const now = new Date();
            const uptime = Math.floor((now - startTime) / 1000);
            document.getElementById('uptime').textContent = formatUptime(uptime);
        })
        .catch(error => {
            console.error('Error loading stats:', error);
        });
}

function formatUptime(seconds) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    if (hours > 0) {
        return hours + 'h ' + minutes + 'm ' + secs + 's';
    } else if (minutes > 0) {
        return minutes + 'm ' + secs + 's';
    } else {
        return secs + 's';
    }
}

// Load routes
function loadRoutes() {
    fetch('/_admin/api/routes')
        .then(response => response.json())
        .then(routes => {
            const container = document.getElementById('routes-list');
            container.innerHTML = '';
            
            routes.forEach(route => {
                const routeDiv = document.createElement('div');
                routeDiv.className = 'route-item';
                
                const methodClass = 'method-' + route.method.toLowerCase();
                
                routeDiv.innerHTML = '<div class="route-info"><span class="method-badge ' + methodClass + '">' + route.method + '</span><code class="route-path">' + route.path + '</code><span>' + (route.description || 'No description') + '</span></div><div><button class="btn btn-secondary" onclick="testRoute(\\''+route.method+'\\', \\''+route.path+'\\')">Test</button></div>';
                
                container.appendChild(routeDiv);
            });
        })
        .catch(error => {
            console.error('Error loading routes:', error);
            document.getElementById('routes-list').innerHTML = '<p>Error loading routes</p>';
        });
}

function refreshRoutes() {
    loadRoutes();
}

// Test a specific route
function testRoute(method, path) {
    showTab('testing');
    document.getElementById('test-method').value = method;
    document.getElementById('test-url').value = path;
}

// Load OpenAPI specification
function loadOpenAPI() {
    fetch('/_admin/api/openapi')
        .then(response => response.json())
        .then(spec => {
            document.getElementById('openapi-preview').innerHTML = '<pre>' + JSON.stringify(spec, null, 2) + '</pre>';
        })
        .catch(error => {
            console.error('Error loading OpenAPI spec:', error);
            document.getElementById('openapi-preview').innerHTML = '<p>Error loading OpenAPI specification</p>';
        });
}

function downloadOpenAPI() {
    fetch('/_admin/api/openapi')
        .then(response => response.json())
        .then(spec => {
            const blob = new Blob([JSON.stringify(spec, null, 2)], { type: 'application/json' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'facade-openapi.json';
            a.click();
            URL.revokeObjectURL(url);
        });
}

function copyOpenAPIUrl() {
    const url = window.location.origin + '/_admin/api/openapi';
    navigator.clipboard.writeText(url).then(() => {
        alert('OpenAPI URL copied to clipboard: ' + url);
    });
}

// API Testing
document.addEventListener('DOMContentLoaded', function() {
    const testForm = document.getElementById('test-form');
    if (testForm) {
        testForm.addEventListener('submit', function(e) {
            e.preventDefault();
            sendTestRequest();
        });
    }
    
    // Load initial stats
    loadStats();
    
    // Auto-refresh stats every 5 seconds
    setInterval(loadStats, 5000);
});

function sendTestRequest() {
    const method = document.getElementById('test-method').value;
    const url = document.getElementById('test-url').value;
    const headers = document.getElementById('test-headers').value;
    const body = document.getElementById('test-body').value;
    
    let requestOptions = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        }
    };
    
    // Parse headers
    try {
        if (headers.trim()) {
            const parsedHeaders = JSON.parse(headers);
            requestOptions.headers = { ...requestOptions.headers, ...parsedHeaders };
        }
    } catch (e) {
        alert('Invalid JSON in headers');
        return;
    }
    
    // Add body for POST/PUT requests
    if ((method === 'POST' || method === 'PUT') && body.trim()) {
        requestOptions.body = body;
    }
    
    // Make the request
    const fullUrl = window.location.origin + url;
    
    fetch(fullUrl, requestOptions)
        .then(async response => {
            const responseText = await response.text();
            let responseJson;
            
            try {
                responseJson = JSON.parse(responseText);
            } catch (e) {
                responseJson = responseText;
            }
            
            const responseContainer = document.getElementById('test-response');
            const responseContent = document.getElementById('response-content');
            
            const responseData = {
                status: response.status,
                statusText: response.statusText,
                headers: Object.fromEntries(response.headers.entries()),
                body: responseJson
            };
            
            responseContent.textContent = JSON.stringify(responseData, null, 2);
            responseContainer.style.display = 'block';
        })
        .catch(error => {
            const responseContainer = document.getElementById('test-response');
            const responseContent = document.getElementById('response-content');
            
            responseContent.textContent = 'Error: ' + error.message;
            responseContainer.style.display = 'block';
        });
}

// State management
function loadState() {
    fetch('/_admin/api/state')
        .then(response => response.json())
        .then(data => {
            document.getElementById('state-content').innerHTML = '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
        })
        .catch(error => {
            console.error('Error loading state:', error);
            document.getElementById('state-content').innerHTML = '<p>Error loading state</p>';
        });
}

function refreshState() {
    loadState();
}

function clearState() {
    if (confirm('Are you sure you want to clear all state data?')) {
        fetch('/_facade/state', { method: 'DELETE' })
            .then(response => response.json())
            .then(data => {
                alert('State cleared successfully');
                loadState();
            })
            .catch(error => {
                console.error('Error clearing state:', error);
                alert('Error clearing state');
            });
    }
}`
}
package admin

func (a *Admin) getCSS() string {
	return `
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f8fafc;
    color: #1a202c;
    line-height: 1.6;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

.header {
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    margin-bottom: 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.header h1 {
    color: #2d3748;
    font-size: 28px;
    font-weight: 700;
}

.server-info {
    display: flex;
    gap: 16px;
    align-items: center;
    font-size: 14px;
    color: #718096;
}

.status-indicator {
    font-weight: 600;
}

.status-indicator.online {
    color: #48bb78;
}

.nav-tabs {
    display: flex;
    background: white;
    border-radius: 12px;
    padding: 8px;
    margin-bottom: 24px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.tab-button {
    background: none;
    border: none;
    padding: 12px 20px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s;
    color: #718096;
}

.tab-button:hover {
    background: #edf2f7;
    color: #2d3748;
}

.tab-button.active {
    background: #4299e1;
    color: white;
}

.tab-content {
    display: none;
}

.tab-content.active {
    display: block;
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
    margin-bottom: 32px;
}

.stat-card {
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    text-align: center;
}

.stat-card h3 {
    font-size: 14px;
    color: #718096;
    margin-bottom: 8px;
    text-transform: uppercase;
    font-weight: 600;
    letter-spacing: 0.05em;
}

.stat-value {
    font-size: 32px;
    font-weight: 700;
    color: #2d3748;
}

.stat-value.error {
    color: #f56565;
}

.routes-container {
    background: white;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    overflow: hidden;
}

.route-item {
    padding: 20px;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.route-item:last-child {
    border-bottom: none;
}

.route-info {
    display: flex;
    gap: 16px;
    align-items: center;
}

.method-badge {
    padding: 4px 12px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
}

.method-get { background: #c6f6d5; color: #22543d; }
.method-post { background: #bee3f8; color: #2a4365; }
.method-put { background: #fbb6ce; color: #742a2a; }
.method-delete { background: #fed7e2; color: #742a2a; }

.route-path {
    font-family: 'Monaco', 'Menlo', monospace;
    background: #f7fafc;
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 14px;
}

.btn {
    padding: 10px 20px;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s;
    text-decoration: none;
    display: inline-block;
}

.btn-primary {
    background: #4299e1;
    color: white;
}

.btn-primary:hover {
    background: #3182ce;
}

.btn-secondary {
    background: #e2e8f0;
    color: #2d3748;
}

.btn-secondary:hover {
    background: #cbd5e0;
}

.btn-danger {
    background: #f56565;
    color: white;
}

.btn-danger:hover {
    background: #e53e3e;
}

.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: #2d3748;
}

.form-group input,
.form-group select,
.form-group textarea {
    width: 100%;
    padding: 12px;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    font-size: 14px;
}

.form-group textarea {
    min-height: 100px;
    font-family: 'Monaco', 'Menlo', monospace;
}

.response-container {
    margin-top: 24px;
    background: white;
    padding: 20px;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.response-container pre {
    background: #1a202c;
    color: #e2e8f0;
    padding: 16px;
    border-radius: 8px;
    overflow-x: auto;
    font-size: 14px;
}

.charts-section,
.swagger-container,
.testing-container,
.state-container {
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.routes-header,
.swagger-actions,
.state-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.openapi-preview {
    background: #f7fafc;
    padding: 20px;
    border-radius: 8px;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 14px;
    max-height: 500px;
    overflow-y: auto;
}

.state-content {
    background: #f7fafc;
    padding: 20px;
    border-radius: 8px;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 14px;
    max-height: 400px;
    overflow-y: auto;
}
`
}
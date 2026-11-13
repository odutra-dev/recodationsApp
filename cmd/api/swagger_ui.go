package main

// GetSwaggerUIHTML returns the complete HTML page with interactive theme switcher.
func GetSwaggerUIHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Recodations API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
  <style>
    /* Theme Switcher Widget */
    .theme-controller {
      position: fixed;
      top: 12px;
      right: 20px;
      z-index: 99999;
      display: flex;
      align-items: center;
      gap: 8px;
      background: rgba(255, 255, 255, 0.85);
      backdrop-filter: blur(12px);
      padding: 6px 14px;
      border-radius: 999px;
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
      border: 1px solid rgba(0, 0, 0, 0.1);
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      font-size: 13px;
      color: #333;
      transition: all 0.3s ease;
    }
    .theme-controller label {
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 6px;
      cursor: pointer;
    }
    .theme-select {
      background: transparent;
      border: none;
      font-family: inherit;
      font-size: 13px;
      font-weight: 600;
      color: inherit;
      outline: none;
      cursor: pointer;
      padding-right: 4px;
    }

    /* Common Base Styles */
    body {
      margin: 0;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
      transition: background-color 0.3s ease, color 0.3s ease;
    }
    .swagger-ui .topbar { display: none; }
    .swagger-ui .info { margin: 24px 0; }
    .swagger-ui .info .title { font-family: 'Inter', sans-serif; font-weight: 700; }
    .swagger-ui code, .swagger-ui .microlight { font-family: 'JetBrains Mono', monospace !important; }

    /* ==========================================================================
       DARK THEME
       ========================================================================== */
    body.theme-dark {
      background-color: #121212 !important;
      color: #e0e0e0 !important;
    }
    body.theme-dark .theme-controller {
      background: rgba(30, 30, 30, 0.85);
      border-color: rgba(255, 255, 255, 0.15);
      color: #f0f0f0;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
    }
    body.theme-dark .theme-select option {
      background: #1e1e1e;
      color: #fff;
    }
    body.theme-dark .swagger-ui {
      color: #e0e0e0;
    }
    body.theme-dark .swagger-ui .info .title,
    body.theme-dark .swagger-ui .info h1,
    body.theme-dark .swagger-ui .info h2,
    body.theme-dark .swagger-ui .info h3,
    body.theme-dark .swagger-ui .info h4,
    body.theme-dark .swagger-ui .info h5,
    body.theme-dark .swagger-ui .info p,
    body.theme-dark .swagger-ui .info li,
    body.theme-dark .swagger-ui .info table,
    body.theme-dark .swagger-ui .opblock-tag,
    body.theme-dark .swagger-ui .opblock .opblock-summary-operation-id,
    body.theme-dark .swagger-ui .opblock .opblock-summary-path,
    body.theme-dark .swagger-ui .opblock .opblock-summary-path__deprecated,
    body.theme-dark .swagger-ui .opblock .opblock-summary-description,
    body.theme-dark .swagger-ui .tab li button,
    body.theme-dark .swagger-ui label,
    body.theme-dark .swagger-ui .btn,
    body.theme-dark .swagger-ui .response-col_status,
    body.theme-dark .swagger-ui .response-col_links,
    body.theme-dark .swagger-ui table thead tr td,
    body.theme-dark .swagger-ui table thead tr th,
    body.theme-dark .swagger-ui .model-title,
    body.theme-dark .swagger-ui .model,
    body.theme-dark .swagger-ui .prop-name,
    body.theme-dark .swagger-ui .parameter__name,
    body.theme-dark .swagger-ui .parameter__type {
      color: #e0e0e0 !important;
    }
    body.theme-dark .swagger-ui .opblock {
      background: #1c1c1e !important;
      border-color: #2c2c2e !important;
    }
    body.theme-dark .swagger-ui .opblock .opblock-summary {
      border-color: #2c2c2e !important;
    }
    body.theme-dark .swagger-ui .opblock-body pre.microlight,
    body.theme-dark .swagger-ui .highlight-code {
      background: #141416 !important;
      color: #79c0ff !important;
    }
    body.theme-dark .swagger-ui section.models {
      background: #1c1c1e !important;
      border-color: #2c2c2e !important;
    }
    body.theme-dark .swagger-ui section.models .model-container {
      background: #141416 !important;
    }
    body.theme-dark .swagger-ui input[type=text],
    body.theme-dark .swagger-ui input[type=password],
    body.theme-dark .swagger-ui textarea,
    body.theme-dark .swagger-ui select {
      background: #242426 !important;
      color: #fff !important;
      border-color: #3a3a3c !important;
    }
    body.theme-dark .swagger-ui .scheme-container {
      background: #18181a !important;
      box-shadow: none !important;
      border-bottom: 1px solid #2c2c2e;
    }
    body.theme-dark .swagger-ui .dialog-ux .modal-ux {
      background: #1c1c1e !important;
      border-color: #2c2c2e !important;
      color: #fff !important;
    }
    body.theme-dark .swagger-ui .dialog-ux .modal-ux-header {
      border-color: #2c2c2e !important;
    }

    /* ==========================================================================
       DRACULA THEME
       ========================================================================== */
    body.theme-dracula {
      background-color: #282a36 !important;
      color: #f8f8f2 !important;
    }
    body.theme-dracula .theme-controller {
      background: rgba(68, 71, 90, 0.9);
      border-color: #bd93f9;
      color: #f8f8f2;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.6);
    }
    body.theme-dracula .theme-select option {
      background: #282a36;
      color: #f8f8f2;
    }
    body.theme-dracula .swagger-ui {
      color: #f8f8f2;
    }
    body.theme-dracula .swagger-ui .info .title {
      color: #bd93f9 !important;
    }
    body.theme-dracula .swagger-ui .info h1,
    body.theme-dracula .swagger-ui .info h2,
    body.theme-dracula .swagger-ui .info h3,
    body.theme-dracula .swagger-ui .info p,
    body.theme-dracula .swagger-ui .opblock-tag,
    body.theme-dracula .swagger-ui .opblock .opblock-summary-path,
    body.theme-dracula .swagger-ui .model-title,
    body.theme-dracula .swagger-ui .parameter__name {
      color: #f8f8f2 !important;
    }
    body.theme-dracula .swagger-ui .opblock {
      background: #21222c !important;
      border-color: #44475a !important;
    }
    body.theme-dracula .swagger-ui .opblock-body pre.microlight {
      background: #191a21 !important;
      color: #8be9fd !important;
    }
    body.theme-dracula .swagger-ui section.models {
      background: #21222c !important;
      border-color: #44475a !important;
    }
    body.theme-dracula .swagger-ui input[type=text],
    body.theme-dracula .swagger-ui input[type=password],
    body.theme-dracula .swagger-ui textarea,
    body.theme-dracula .swagger-ui select {
      background: #44475a !important;
      color: #f8f8f2 !important;
      border-color: #6272a4 !important;
    }
    body.theme-dracula .swagger-ui .scheme-container {
      background: #282a36 !important;
      box-shadow: none !important;
      border-bottom: 1px solid #44475a;
    }
    body.theme-dracula .swagger-ui .btn.authorize {
      color: #50fa7b !important;
      border-color: #50fa7b !important;
    }

    /* ==========================================================================
       NORD THEME
       ========================================================================== */
    body.theme-nord {
      background-color: #2e3440 !important;
      color: #eceff4 !important;
    }
    body.theme-nord .theme-controller {
      background: rgba(59, 66, 82, 0.9);
      border-color: #88c0d0;
      color: #eceff4;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
    }
    body.theme-nord .theme-select option {
      background: #2e3440;
      color: #eceff4;
    }
    body.theme-nord .swagger-ui {
      color: #eceff4;
    }
    body.theme-nord .swagger-ui .info .title {
      color: #88c0d0 !important;
    }
    body.theme-nord .swagger-ui .info p,
    body.theme-nord .swagger-ui .opblock-tag,
    body.theme-nord .swagger-ui .opblock .opblock-summary-path,
    body.theme-nord .swagger-ui .parameter__name {
      color: #eceff4 !important;
    }
    body.theme-nord .swagger-ui .opblock {
      background: #3b4252 !important;
      border-color: #4c566a !important;
    }
    body.theme-nord .swagger-ui .opblock-body pre.microlight {
      background: #242933 !important;
      color: #81a1c1 !important;
    }
    body.theme-nord .swagger-ui section.models {
      background: #3b4252 !important;
      border-color: #4c566a !important;
    }
    body.theme-nord .swagger-ui input[type=text],
    body.theme-nord .swagger-ui input[type=password],
    body.theme-nord .swagger-ui textarea,
    body.theme-nord .swagger-ui select {
      background: #434c5e !important;
      color: #eceff4 !important;
      border-color: #4c566a !important;
    }
    body.theme-nord .swagger-ui .scheme-container {
      background: #2e3440 !important;
      box-shadow: none !important;
      border-bottom: 1px solid #4c566a;
    }

    /* ==========================================================================
       MONOKAI THEME
       ========================================================================== */
    body.theme-monokai {
      background-color: #272822 !important;
      color: #f8f8f2 !important;
    }
    body.theme-monokai .theme-controller {
      background: rgba(62, 61, 50, 0.9);
      border-color: #fd971f;
      color: #f8f8f2;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.6);
    }
    body.theme-monokai .theme-select option {
      background: #272822;
      color: #f8f8f2;
    }
    body.theme-monokai .swagger-ui {
      color: #f8f8f2;
    }
    body.theme-monokai .swagger-ui .info .title {
      color: #e6db74 !important;
    }
    body.theme-monokai .swagger-ui .info p,
    body.theme-monokai .swagger-ui .opblock-tag,
    body.theme-monokai .swagger-ui .opblock .opblock-summary-path,
    body.theme-monokai .swagger-ui .parameter__name {
      color: #f8f8f2 !important;
    }
    body.theme-monokai .swagger-ui .opblock {
      background: #1e1f1c !important;
      border-color: #3e3d32 !important;
    }
    body.theme-monokai .swagger-ui .opblock-body pre.microlight {
      background: #171814 !important;
      color: #a6e22e !important;
    }
    body.theme-monokai .swagger-ui section.models {
      background: #1e1f1c !important;
      border-color: #3e3d32 !important;
    }
    body.theme-monokai .swagger-ui input[type=text],
    body.theme-monokai .swagger-ui input[type=password],
    body.theme-monokai .swagger-ui textarea,
    body.theme-monokai .swagger-ui select {
      background: #3e3d32 !important;
      color: #f8f8f2 !important;
      border-color: #75715e !important;
    }
    body.theme-monokai .swagger-ui .scheme-container {
      background: #272822 !important;
      box-shadow: none !important;
      border-bottom: 1px solid #3e3d32;
    }

    /* Light Theme (explicit) */
    body.theme-light {
      background-color: #fafafa !important;
      color: #3b4151 !important;
    }
  </style>
</head>
<body class="theme-light">

<!-- Theme Selection Widget -->
<div class="theme-controller">
  <label for="themeSelector">
    <span id="themeIcon">🌓</span>
    <span>Theme:</span>
  </label>
  <select id="themeSelector" class="theme-select">
    <option value="system">System (Auto)</option>
    <option value="light">☀️ Light</option>
    <option value="dark">🌙 Dark</option>
    <option value="dracula">🧛 Dracula</option>
    <option value="nord">❄️ Nord</option>
    <option value="monokai">🍂 Monokai</option>
  </select>
</div>

<div id="swagger-ui"></div>

<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js" crossorigin></script>

<script>
  // Theme Management Logic
  const STORAGE_KEY = 'recodations_swagger_theme';
  const themeSelector = document.getElementById('themeSelector');
  const themeIcon = document.getElementById('themeIcon');
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

  function getSystemTheme() {
    return mediaQuery.matches ? 'dark' : 'light';
  }

  function applyTheme(theme) {
    let effectiveTheme = theme;
    if (theme === 'system') {
      effectiveTheme = getSystemTheme();
    }

    // Remove existing theme classes
    document.body.classList.remove('theme-light', 'theme-dark', 'theme-dracula', 'theme-nord', 'theme-monokai');
    document.body.classList.add('theme-' + effectiveTheme);

    // Update icon indicator
    switch(effectiveTheme) {
      case 'dark':
        themeIcon.textContent = '🌙';
        break;
      case 'dracula':
        themeIcon.textContent = '🧛';
        break;
      case 'nord':
        themeIcon.textContent = '❄️';
        break;
      case 'monokai':
        themeIcon.textContent = '🍂';
        break;
      default:
        themeIcon.textContent = '☀️';
    }

    if (theme === 'system') {
      themeIcon.textContent = '🖥️';
    }
  }

  // Initial Load: defaults to 'system'
  const savedTheme = localStorage.getItem(STORAGE_KEY) || 'system';
  themeSelector.value = savedTheme;
  applyTheme(savedTheme);

  // Listen to User Dropdown Selection
  themeSelector.addEventListener('change', (e) => {
    const selected = e.target.value;
    localStorage.setItem(STORAGE_KEY, selected);
    applyTheme(selected);
  });

  // Listen to OS System Theme Changes dynamically
  mediaQuery.addEventListener('change', () => {
    if (themeSelector.value === 'system') {
      applyTheme('system');
    }
  });

  // Initialize Swagger UI
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/swagger/doc.json',
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
      ],
      layout: "StandaloneLayout",
    });
  };
</script>
</body>
</html>`
}

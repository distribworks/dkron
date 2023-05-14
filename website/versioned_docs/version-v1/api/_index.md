---
title: "API"
weight: 100
---

<link
  rel="stylesheet"
  type="text/css"
  href="https://unpkg.com/swagger-ui-dist@3.51.0/swagger-ui.css"
/>
<style>
  body {
    line-height: 1.7;
  }
  .swagger-ui .info .title small pre {
    background-color: inherit;
    padding: inherit;
  }
  .swagger-ui .scheme-container { display: none !important; } 
</style>

<div id="swagger-ui"></div>

<script src="https://unpkg.com/swagger-ui-dist@3.51.0/swagger-ui-standalone-preset.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@3.51.0/swagger-ui-bundle.js"></script>
<script>
  window.onload = function () {
    // Begin Swagger UI call region
    const ui = SwaggerUIBundle({
      url: "/v1.2/swagger.yaml",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      plugins: [
        SwaggerUIBundle.plugins.DownloadUrl
      ],
      layout: "BaseLayout"
    })
    // End Swagger UI call region

    window.ui = ui
  }
</script>

---
title: "API"
weight: 100
---

<link rel="stylesheet" type="text/css" href="/css/swagger-ui.css">
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

<script src="/js/swagger-ui-bundle.js"> </script>
<script src="/js/swagger-ui-standalone-preset.js"> </script>
<script>
  window.onload = function () {
    // Begin Swagger UI call region
    const ui = SwaggerUIBundle({
      url: "/v2.0/swagger.yaml",
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

angular-rickshaw
================

An AngularJS directive for [Rickshaw](http://code.shutterstock.com/rickshaw/).

[Documentation / Demo](http://ngyewch.github.io/angular-rickshaw/).


Requirements
============

Include Rickshaw (CSS, JS), D3 (JS) and the Angular Rickshaw (JS) in your `<head>` section:
```html
<link rel="stylesheet" href="PATH_TO_CSS/rickshaw/rickshaw.min.css">

<script src="PATH_TO_JS/d3/d3.min.js"></script>
<script src="PATH_TO_JS/rickshaw/rickshaw.min.js"></script>
<script src="PATH_TO_JS/angular-rickshaw/rickshaw.min.js"></script>
```

Usage
=====

Create a div with the ```rickshaw``` attribute or a ```rickshaw``` element to instantiate a Rickshaw graph.

```html
<rickshaw
    rickshaw-options="options"
    rickshaw-features="features"
    rickshaw-series="series">
</rickshaw>
```

The ```rickshaw-options``` attribute provides access to Rickshaw configuration.
The ```rickshaw-features``` attribute provides access to Rickshaw features.
The ```rickshaw-series``` attribute provides access to Rickshaw series data.

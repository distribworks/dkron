"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[7416],{3905:(e,r,t)=>{t.d(r,{Zo:()=>d,kt:()=>b});var n=t(7294);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function a(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function s(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?a(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):a(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function i(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)t=a[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)t=a[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var c=n.createContext({}),u=function(e){var r=n.useContext(c),t=r;return e&&(t="function"==typeof e?e(r):s(s({},r),e)),t},d=function(e){var r=u(e.components);return n.createElement(c.Provider,{value:r},e.children)},l="mdxType",p={inlineCode:"code",wrapper:function(e){var r=e.children;return n.createElement(n.Fragment,{},r)}},f=n.forwardRef((function(e,r){var t=e.components,o=e.mdxType,a=e.originalType,c=e.parentName,d=i(e,["components","mdxType","originalType","parentName"]),l=u(t),f=o,b=l["".concat(c,".").concat(f)]||l[f]||p[f]||a;return t?n.createElement(b,s(s({ref:r},d),{},{components:t})):n.createElement(b,s({ref:r},d))}));function b(e,r){var t=arguments,o=r&&r.mdxType;if("string"==typeof e||o){var a=t.length,s=new Array(a);s[0]=f;var i={};for(var c in r)hasOwnProperty.call(r,c)&&(i[c]=r[c]);i.originalType=e,i[l]="string"==typeof e?e:o,s[1]=i;for(var u=2;u<a;u++)s[u]=t[u];return n.createElement.apply(null,s)}return n.createElement.apply(null,t)}f.displayName="MDXCreateElement"},4846:(e,r,t)=>{t.r(r),t.d(r,{assets:()=>c,contentTitle:()=>s,default:()=>l,frontMatter:()=>a,metadata:()=>i,toc:()=>u});var n=t(7462),o=(t(7294),t(3905));const a={title:"Embedded storage"},s=void 0,i={unversionedId:"usage/storage",id:"version-v2/usage/storage",title:"Embedded storage",description:"Dkron has an embedded distributed KV store engine based on BadgerDB. This works out of the box on each dkron server.",source:"@site/versioned_docs/version-v2/usage/storage.md",sourceDirName:"usage",slug:"/usage/storage",permalink:"/docs/v2/usage/storage",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v2/usage/storage.md",tags:[],version:"v2",frontMatter:{title:"Embedded storage"},sidebar:"tutorialSidebar",previous:{title:"Job retries",permalink:"/docs/v2/usage/retries"},next:{title:"Target nodes spec",permalink:"/docs/v2/usage/target-nodes-spec"}},c={},u=[],d={toc:u};function l(e){let{components:r,...t}=e;return(0,o.kt)("wrapper",(0,n.Z)({},d,t,{components:r,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Dkron has an embedded distributed KV store engine based on BadgerDB. This works out of the box on each dkron server."),(0,o.kt)("p",null,"This ensures a dead easy install and setup, basically run dkron and you will have a full working node."))}l.isMDXComponent=!0}}]);
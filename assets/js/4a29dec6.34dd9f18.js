"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[6236],{3905:(e,r,t)=>{t.d(r,{Zo:()=>u,kt:()=>g});var o=t(7294);function n(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function s(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);r&&(o=o.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,o)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?s(Object(t),!0).forEach((function(r){n(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):s(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function c(e,r){if(null==e)return{};var t,o,n=function(e,r){if(null==e)return{};var t,o,n={},s=Object.keys(e);for(o=0;o<s.length;o++)t=s[o],r.indexOf(t)>=0||(n[t]=e[t]);return n}(e,r);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(o=0;o<s.length;o++)t=s[o],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(n[t]=e[t])}return n}var i=o.createContext({}),l=function(e){var r=o.useContext(i),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},u=function(e){var r=l(e.components);return o.createElement(i.Provider,{value:r},e.children)},p="mdxType",d={inlineCode:"code",wrapper:function(e){var r=e.children;return o.createElement(o.Fragment,{},r)}},f=o.forwardRef((function(e,r){var t=e.components,n=e.mdxType,s=e.originalType,i=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),p=l(t),f=n,g=p["".concat(i,".").concat(f)]||p[f]||d[f]||s;return t?o.createElement(g,a(a({ref:r},u),{},{components:t})):o.createElement(g,a({ref:r},u))}));function g(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var s=t.length,a=new Array(s);a[0]=f;var c={};for(var i in r)hasOwnProperty.call(r,i)&&(c[i]=r[i]);c.originalType=e,c[p]="string"==typeof e?e:n,a[1]=c;for(var l=2;l<s;l++)a[l]=t[l];return o.createElement.apply(null,a)}return o.createElement.apply(null,t)}f.displayName="MDXCreateElement"},657:(e,r,t)=>{t.r(r),t.d(r,{assets:()=>i,contentTitle:()=>a,default:()=>p,frontMatter:()=>s,metadata:()=>c,toc:()=>l});var o=t(7462),n=(t(7294),t(3905));const s={title:"Log Processor"},a=void 0,c={unversionedId:"usage/processors/log",id:"version-v2/usage/processors/log",title:"Log Processor",description:"Log processor writes the execution output to stdout/stderr",source:"@site/versioned_docs/version-v2/usage/processors/log.md",sourceDirName:"usage/processors",slug:"/usage/processors/log",permalink:"/docs/v2/usage/processors/log",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v2/usage/processors/log.md",tags:[],version:"v2",frontMatter:{title:"Log Processor"},sidebar:"tutorialSidebar",previous:{title:"File Processor",permalink:"/docs/v2/usage/processors/file"},next:{title:"Syslog Processor",permalink:"/docs/v2/usage/processors/syslog"}},i={},l=[{value:"Configuration",id:"configuration",level:2}],u={toc:l};function p(e){let{components:r,...t}=e;return(0,n.kt)("wrapper",(0,o.Z)({},u,t,{components:r,mdxType:"MDXLayout"}),(0,n.kt)("p",null,"Log processor writes the execution output to stdout/stderr"),(0,n.kt)("h2",{id:"configuration"},"Configuration"),(0,n.kt)("p",null,"Parameters"),(0,n.kt)("p",null,(0,n.kt)("inlineCode",{parentName:"p"},"forward: Forward the output to the next processor")),(0,n.kt)("p",null,"Example"),(0,n.kt)("pre",null,(0,n.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "job_name",\n    "command": "echo \'Hello log\'",\n    "schedule": "@every 2m",\n    "tags": {\n        "role": "web"\n    },\n    "processors": {\n        "log": {\n            "forward": true\n        }\n    }\n}\n')))}p.isMDXComponent=!0}}]);
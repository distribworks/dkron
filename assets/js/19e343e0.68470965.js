"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[1490],{3905:(e,r,t)=>{t.d(r,{Zo:()=>u,kt:()=>g});var o=t(7294);function n(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function s(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);r&&(o=o.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,o)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?s(Object(t),!0).forEach((function(r){n(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):s(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function i(e,r){if(null==e)return{};var t,o,n=function(e,r){if(null==e)return{};var t,o,n={},s=Object.keys(e);for(o=0;o<s.length;o++)t=s[o],r.indexOf(t)>=0||(n[t]=e[t]);return n}(e,r);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(o=0;o<s.length;o++)t=s[o],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(n[t]=e[t])}return n}var l=o.createContext({}),c=function(e){var r=o.useContext(l),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},u=function(e){var r=c(e.components);return o.createElement(l.Provider,{value:r},e.children)},p="mdxType",y={inlineCode:"code",wrapper:function(e){var r=e.children;return o.createElement(o.Fragment,{},r)}},m=o.forwardRef((function(e,r){var t=e.components,n=e.mdxType,s=e.originalType,l=e.parentName,u=i(e,["components","mdxType","originalType","parentName"]),p=c(t),m=n,g=p["".concat(l,".").concat(m)]||p[m]||y[m]||s;return t?o.createElement(g,a(a({ref:r},u),{},{components:t})):o.createElement(g,a({ref:r},u))}));function g(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var s=t.length,a=new Array(s);a[0]=m;var i={};for(var l in r)hasOwnProperty.call(r,l)&&(i[l]=r[l]);i.originalType=e,i[p]="string"==typeof e?e:n,a[1]=i;for(var c=2;c<s;c++)a[c]=t[c];return o.createElement.apply(null,a)}return o.createElement.apply(null,t)}m.displayName="MDXCreateElement"},7349:(e,r,t)=>{t.r(r),t.d(r,{assets:()=>l,contentTitle:()=>a,default:()=>p,frontMatter:()=>s,metadata:()=>i,toc:()=>c});var o=t(7462),n=(t(7294),t(3905));const s={title:"Syslog Processor"},a=void 0,i={unversionedId:"usage/processors/syslog",id:"version-v1/usage/processors/syslog",title:"Syslog Processor",description:"Syslog processor writes the execution output to the system syslog daemon",source:"@site/versioned_docs/version-v1/usage/processors/syslog.md",sourceDirName:"usage/processors",slug:"/usage/processors/syslog",permalink:"/docs/v1/usage/processors/syslog",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v1/usage/processors/syslog.md",tags:[],version:"v1",frontMatter:{title:"Syslog Processor"},sidebar:"tutorialSidebar",previous:{title:"Log Processor",permalink:"/docs/v1/usage/processors/log"},next:{title:"Job retries",permalink:"/docs/v1/usage/retries"}},l={},c=[{value:"Configuration",id:"configuration",level:2}],u={toc:c};function p(e){let{components:r,...t}=e;return(0,n.kt)("wrapper",(0,o.Z)({},u,t,{components:r,mdxType:"MDXLayout"}),(0,n.kt)("p",null,"Syslog processor writes the execution output to the system syslog daemon"),(0,n.kt)("p",null,"Note: Only work on linux systems"),(0,n.kt)("h2",{id:"configuration"},"Configuration"),(0,n.kt)("p",null,"Parameters"),(0,n.kt)("p",null,(0,n.kt)("inlineCode",{parentName:"p"},"forward: Forward the output to the next processor")),(0,n.kt)("p",null,"Example"),(0,n.kt)("pre",null,(0,n.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "job_name",\n    "command": "echo \'Hello syslog\'",\n    "schedule": "@every 2m",\n    "tags": {\n        "role": "web"\n    },\n    "processors": {\n        "syslog": {\n            "forward": true\n        }\n    }\n}\n')))}p.isMDXComponent=!0}}]);
"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[6176],{2132:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>l,contentTitle:()=>i,default:()=>u,frontMatter:()=>o,metadata:()=>c,toc:()=>a});var n=r(17624),s=r(95788);const o={},i="Shell Executor",c={id:"usage/executors/shell",title:"Shell Executor",description:"Shell executor runs a system command",source:"@site/versioned_docs/version-v3/usage/executors/shell.md",sourceDirName:"usage/executors",slug:"/usage/executors/shell",permalink:"/docs/v3/usage/executors/shell",draft:!1,unlisted:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v3/usage/executors/shell.md",tags:[],version:"v3",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"NATS Executor",permalink:"/docs/v3/usage/executors/nats"},next:{title:"Job metadata",permalink:"/docs/v3/usage/metatags"}},l={},a=[{value:"Configuration",id:"configuration",level:2},{value:"Job execution prometheus metrics",id:"job-execution-prometheus-metrics",level:2},{value:"Exposed metrics",id:"exposed-metrics",level:3}];function d(e){const t={br:"br",code:"code",h1:"h1",h2:"h2",h3:"h3",p:"p",pre:"pre",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",...(0,s.MN)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.h1,{id:"shell-executor",children:"Shell Executor"}),"\n",(0,n.jsx)(t.p,{children:"Shell executor runs a system command"}),"\n",(0,n.jsx)(t.h2,{id:"configuration",children:"Configuration"}),"\n",(0,n.jsx)(t.p,{children:"Params"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{children:"shell: Run this command using a shell environment\ncommand: The command to run\nenv: Env vars separated by comma\ncwd: Chdir before command run\ntimeout: Force kill job after specified time. Format: https://golang.org/pkg/time/#ParseDuration.\n"})}),"\n",(0,n.jsx)(t.p,{children:"Example"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-json",children:'{\n  "executor": "shell",\n  "executor_config": {\n      "shell": "true",\n      "command": "my_command",\n      "env": "ENV_VAR=va1,ANOTHER_ENV_VAR=var2",\n      "cwd": "/app",\n      "timeout": "24h"\n  }\n}\n'})}),"\n",(0,n.jsx)(t.h2,{id:"job-execution-prometheus-metrics",children:"Job execution prometheus metrics"}),"\n",(0,n.jsxs)(t.p,{children:["Path: ",(0,n.jsx)(t.code,{children:"/metrics"}),(0,n.jsx)(t.br,{}),"\n","Port: 9422",(0,n.jsx)(t.br,{}),"\n","or configure via environment variable ",(0,n.jsx)(t.code,{children:"SHELL_EXECUTOR_PROMETHEUS_PORT"})]}),"\n",(0,n.jsx)(t.h3,{id:"exposed-metrics",children:"Exposed metrics"}),"\n",(0,n.jsxs)(t.table,{children:[(0,n.jsx)(t.thead,{children:(0,n.jsxs)(t.tr,{children:[(0,n.jsx)(t.th,{children:"Name"}),(0,n.jsx)(t.th,{style:{textAlign:"left"},children:"Type"}),(0,n.jsx)(t.th,{style:{textAlign:"right"},children:"Description"}),(0,n.jsx)(t.th,{style:{textAlign:"right"},children:"Labels"})]})}),(0,n.jsxs)(t.tbody,{children:[(0,n.jsxs)(t.tr,{children:[(0,n.jsx)(t.td,{children:"dkron_job_cpu_usage"}),(0,n.jsx)(t.td,{style:{textAlign:"left"},children:"gauge"}),(0,n.jsx)(t.td,{style:{textAlign:"right"},children:"current CPU usage by job"}),(0,n.jsx)(t.td,{style:{textAlign:"right"},children:"job_name"})]}),(0,n.jsxs)(t.tr,{children:[(0,n.jsx)(t.td,{children:"dkron_job_mem_usage_kb"}),(0,n.jsx)(t.td,{style:{textAlign:"left"},children:"gauge"}),(0,n.jsx)(t.td,{style:{textAlign:"right"},children:"current memory consumed by job"}),(0,n.jsx)(t.td,{style:{textAlign:"right"},children:"job_name"})]})]})]})]})}function u(e={}){const{wrapper:t}={...(0,s.MN)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},95788:(e,t,r)=>{r.d(t,{MN:()=>a});var n=r(11504);function s(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){s(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function c(e,t){if(null==e)return{};var r,n,s=function(e,t){if(null==e)return{};var r,n,s={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(s[r]=e[r]);return s}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(s[r]=e[r])}return s}var l=n.createContext({}),a=function(e){var t=n.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},d={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},u=n.forwardRef((function(e,t){var r=e.components,s=e.mdxType,o=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),h=a(r),m=s,x=h["".concat(l,".").concat(m)]||h[m]||d[m]||o;return r?n.createElement(x,i(i({ref:t},u),{},{components:r})):n.createElement(x,i({ref:t},u))}));u.displayName="MDXCreateElement"}}]);
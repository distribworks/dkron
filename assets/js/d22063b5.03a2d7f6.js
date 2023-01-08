"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[6213],{3905:(e,t,n)=>{n.d(t,{Zo:()=>i,kt:()=>m});var o=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function a(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);t&&(o=o.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,o)}return n}function s(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?a(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):a(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function u(e,t){if(null==e)return{};var n,o,r=function(e,t){if(null==e)return{};var n,o,r={},a=Object.keys(e);for(o=0;o<a.length;o++)n=a[o],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(o=0;o<a.length;o++)n=a[o],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var l=o.createContext({}),c=function(e){var t=o.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):s(s({},t),e)),n},i=function(e){var t=c(e.components);return o.createElement(l.Provider,{value:t},e.children)},p="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return o.createElement(o.Fragment,{},t)}},g=o.forwardRef((function(e,t){var n=e.components,r=e.mdxType,a=e.originalType,l=e.parentName,i=u(e,["components","mdxType","originalType","parentName"]),p=c(n),g=r,m=p["".concat(l,".").concat(g)]||p[g]||d[g]||a;return n?o.createElement(m,s(s({ref:t},i),{},{components:n})):o.createElement(m,s({ref:t},i))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=n.length,s=new Array(a);s[0]=g;var u={};for(var l in t)hasOwnProperty.call(t,l)&&(u[l]=t[l]);u.originalType=e,u[p]="string"==typeof e?e:r,s[1]=u;for(var c=2;c<a;c++)s[c]=n[c];return o.createElement.apply(null,s)}return o.createElement.apply(null,n)}g.displayName="MDXCreateElement"},5382:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>s,default:()=>p,frontMatter:()=>a,metadata:()=>u,toc:()=>c});var o=n(7462),r=(n(7294),n(3905));const a={title:"Target nodes spec",weight:10},s=void 0,u={unversionedId:"usage/target-nodes-spec",id:"version-v1/usage/target-nodes-spec",title:"Target nodes spec",description:"Target nodes spec",source:"@site/versioned_docs/version-v1/usage/target-nodes-spec.md",sourceDirName:"usage",slug:"/usage/target-nodes-spec",permalink:"/docs/v1/usage/target-nodes-spec",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v1/usage/target-nodes-spec.md",tags:[],version:"v1",frontMatter:{title:"Target nodes spec",weight:10},sidebar:"tutorialSidebar",previous:{title:"Job retries",permalink:"/docs/v1/usage/retries"}},l={},c=[{value:"Target nodes spec",id:"target-nodes-spec",level:2},{value:"Examples:",id:"examples",level:3}],i={toc:c};function p(e){let{components:t,...n}=e;return(0,r.kt)("wrapper",(0,o.Z)({},i,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"target-nodes-spec"},"Target nodes spec"),(0,r.kt)("p",null,"You can choose whether a job is run on a node or nodes by specifying tags and a count of target nodes having this tag do you want a job to run."),(0,r.kt)("p",null,"{{% notice note %}}\nThe target node syntax: ",(0,r.kt)("inlineCode",{parentName:"p"},"[tag-value]:[count]"),"\n{{% /notice %}}"),(0,r.kt)("h3",{id:"examples"},"Examples:"),(0,r.kt)("p",null,"Target all nodes with a tag:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "job_name",\n    "command": "/bin/true",\n    "schedule": "@every 2m",\n    "tags": {\n        "role": "web"\n    }\n}\n')),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-mermaid"},'graph LR;\n    J("Job tags: #quot;role#quot;: #quot;web#quot;") --\x3e|Run Job|N1["Node1 tags: #quot;role#quot;: #quot;web#quot;"]\n    J --\x3e|Run Job|N2["Node2 tags: #quot;role#quot;: #quot;web#quot;"]\n    J --\x3e|Run Job|N3["Node2 tags: #quot;role#quot;: #quot;web#quot;"]\n')),(0,r.kt)("p",null,"Target only one nodes of a group of nodes with a tag:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "job_name",\n    "command": "/bin/true",\n    "schedule": "@every 2m",\n    "tags": {\n        "role": "web:1"\n    }\n}\n')),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-mermaid"},'graph LR;\n    J("Job tags: #quot;role#quot;: #quot;web:1#quot;") --\x3e|Run Job|N1["Node1 tags: #quot;role#quot;: #quot;web#quot;"]\n    J -.- N2["Node2 tags: #quot;role#quot;: #quot;web#quot;"]\n    J -.- N3["Node2 tags: #quot;role#quot;: #quot;web#quot;"]\n')),(0,r.kt)("p",null,"Dkron will try to run the job in the amount of nodes indicated by that count having that tag."))}p.isMDXComponent=!0}}]);
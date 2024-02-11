"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[9088],{93436:(e,r,o)=>{o.r(r),o.d(r,{assets:()=>l,contentTitle:()=>i,default:()=>u,frontMatter:()=>s,metadata:()=>c,toc:()=>a});var t=o(17624),n=o(95788);const s={},i="Processors",c={id:"usage/processors/index",title:"Processors",description:"Execution Processors",source:"@site/versioned_docs/version-v3/usage/processors/index.md",sourceDirName:"usage/processors",slug:"/usage/processors/",permalink:"/docs/v3/usage/processors/",draft:!1,unlisted:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v3/usage/processors/index.md",tags:[],version:"v3",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"Developing plugins",permalink:"/docs/v3/usage/plugins/develop"},next:{title:"File Processor",permalink:"/docs/v3/usage/processors/file"}},l={},a=[{value:"Execution Processors",id:"execution-processors",level:2},{value:"Built-in processors",id:"built-in-processors",level:3}];function p(e){const r={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",ol:"ol",p:"p",ul:"ul",...(0,n.MN)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(r.h1,{id:"processors",children:"Processors"}),"\n",(0,t.jsx)(r.h2,{id:"execution-processors",children:"Execution Processors"}),"\n",(0,t.jsx)(r.p,{children:"Processor plugins are called when an execution response has been received. They are passed the resulting execution data and configuration parameters, this plugins can perform a variety of operations with the execution and it's very flexible and per Job, examples of operations this plugins can do:"}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Execution output storage, forwarding or redirection."}),"\n",(0,t.jsx)(r.li,{children:"Notification"}),"\n",(0,t.jsx)(r.li,{children:"Monitoring"}),"\n"]}),"\n",(0,t.jsx)(r.p,{children:"For example, Processor plugins can be used to redirect the output of a job execution to different targets."}),"\n",(0,t.jsx)(r.p,{children:"Currently Dkron provides you with some built-in plugins but the list keeps growing. Some of the features previously implemented in the application will be progessively moved to plugins."}),"\n",(0,t.jsx)(r.h3,{id:"built-in-processors",children:"Built-in processors"}),"\n",(0,t.jsx)(r.p,{children:"Dkron provides the following built-in processors:"}),"\n",(0,t.jsxs)(r.ol,{start:"0",children:["\n",(0,t.jsx)(r.li,{children:"not specified - Store the output in the key value store (Slow performance, good for testing, default method)"}),"\n",(0,t.jsx)(r.li,{children:"log - Output the execution log to Dkron stdout (Good performance, needs parsing)"}),"\n",(0,t.jsx)(r.li,{children:"syslog - Output to the syslog (Good performance, needs parsing)"}),"\n",(0,t.jsx)(r.li,{children:"files - Output to multiple files (Good performance, needs parsing)"}),"\n"]}),"\n",(0,t.jsxs)(r.p,{children:[(0,t.jsx)(r.a,{href:"/pro/",children:"Dkro Pro"})," provides you with several more processors."]}),"\n",(0,t.jsxs)(r.p,{children:["All plugins accepts one configuration option: ",(0,t.jsx)(r.code,{children:"forward"})," Indicated if the plugin must forward the original execution output. This allows for chaining plugins and sending output to different targets at the same time."]})]})}function u(e={}){const{wrapper:r}={...(0,n.MN)(),...e.components};return r?(0,t.jsx)(r,{...e,children:(0,t.jsx)(p,{...e})}):p(e)}},95788:(e,r,o)=>{o.d(r,{MN:()=>a});var t=o(11504);function n(e,r,o){return r in e?Object.defineProperty(e,r,{value:o,enumerable:!0,configurable:!0,writable:!0}):e[r]=o,e}function s(e,r){var o=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);r&&(t=t.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),o.push.apply(o,t)}return o}function i(e){for(var r=1;r<arguments.length;r++){var o=null!=arguments[r]?arguments[r]:{};r%2?s(Object(o),!0).forEach((function(r){n(e,r,o[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(o)):s(Object(o)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(o,r))}))}return e}function c(e,r){if(null==e)return{};var o,t,n=function(e,r){if(null==e)return{};var o,t,n={},s=Object.keys(e);for(t=0;t<s.length;t++)o=s[t],r.indexOf(o)>=0||(n[o]=e[o]);return n}(e,r);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(t=0;t<s.length;t++)o=s[t],r.indexOf(o)>=0||Object.prototype.propertyIsEnumerable.call(e,o)&&(n[o]=e[o])}return n}var l=t.createContext({}),a=function(e){var r=t.useContext(l),o=r;return e&&(o="function"==typeof e?e(r):i(i({},r),e)),o},p={inlineCode:"code",wrapper:function(e){var r=e.children;return t.createElement(t.Fragment,{},r)}},u=t.forwardRef((function(e,r){var o=e.components,n=e.mdxType,s=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),d=a(o),f=n,h=d["".concat(l,".").concat(f)]||d[f]||p[f]||s;return o?t.createElement(h,i(i({ref:r},u),{},{components:o})):t.createElement(h,i({ref:r},u))}));u.displayName="MDXCreateElement"}}]);
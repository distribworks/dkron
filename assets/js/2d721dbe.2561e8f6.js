"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[4599],{3905:(e,r,t)=>{t.d(r,{Zo:()=>l,kt:()=>f});var o=t(7294);function n(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function s(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);r&&(o=o.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,o)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?s(Object(t),!0).forEach((function(r){n(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):s(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function i(e,r){if(null==e)return{};var t,o,n=function(e,r){if(null==e)return{};var t,o,n={},s=Object.keys(e);for(o=0;o<s.length;o++)t=s[o],r.indexOf(t)>=0||(n[t]=e[t]);return n}(e,r);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(o=0;o<s.length;o++)t=s[o],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(n[t]=e[t])}return n}var c=o.createContext({}),p=function(e){var r=o.useContext(c),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},l=function(e){var r=p(e.components);return o.createElement(c.Provider,{value:r},e.children)},m="mdxType",u={inlineCode:"code",wrapper:function(e){var r=e.children;return o.createElement(o.Fragment,{},r)}},d=o.forwardRef((function(e,r){var t=e.components,n=e.mdxType,s=e.originalType,c=e.parentName,l=i(e,["components","mdxType","originalType","parentName"]),m=p(t),d=n,f=m["".concat(c,".").concat(d)]||m[d]||u[d]||s;return t?o.createElement(f,a(a({ref:r},l),{},{components:t})):o.createElement(f,a({ref:r},l))}));function f(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var s=t.length,a=new Array(s);a[0]=d;var i={};for(var c in r)hasOwnProperty.call(r,c)&&(i[c]=r[c]);i.originalType=e,i[m]="string"==typeof e?e:n,a[1]=i;for(var p=2;p<s;p++)a[p]=t[p];return o.createElement.apply(null,a)}return o.createElement.apply(null,t)}d.displayName="MDXCreateElement"},9797:(e,r,t)=>{t.r(r),t.d(r,{assets:()=>c,contentTitle:()=>a,default:()=>m,frontMatter:()=>s,metadata:()=>i,toc:()=>p});var o=t(7462),n=(t(7294),t(3905));const s={},a="Email processor",i={unversionedId:"pro/processors/email",id:"pro/processors/email",title:"Email processor",description:"The Email processor provides flexibility to job email notifications.",source:"@site/docs/pro/processors/email.md",sourceDirName:"pro/processors",slug:"/pro/processors/email",permalink:"/docs/pro/processors/email",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/pro/processors/email.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"Elasticsearch processor",permalink:"/docs/pro/processors/elasticsearch"},next:{title:"Slack processor",permalink:"/docs/pro/processors/slack"}},c={},p=[],l={toc:p};function m(e){let{components:r,...t}=e;return(0,n.kt)("wrapper",(0,o.Z)({},l,t,{components:r,mdxType:"MDXLayout"}),(0,n.kt)("h1",{id:"email-processor"},"Email processor"),(0,n.kt)("p",null,"The Email processor provides flexibility to job email notifications."),(0,n.kt)("p",null,"Configuration of the email processor is stored in a file named ",(0,n.kt)("inlineCode",{parentName:"p"},"dkron-processor-email.yml")," in the same locations as ",(0,n.kt)("inlineCode",{parentName:"p"},"dkron.yml"),", and should include a list of providers, it can include any number of providers."),(0,n.kt)("p",null,"Example:"),(0,n.kt)("pre",null,(0,n.kt)("code",{parentName:"pre",className:"language-yaml"},"provider1:\n  host: smtp.myprovider.com\n  port: 25\n  username: myusername\n  password: mypassword\n  from: cron@mycompany.com\n  subjectPrefix: '[Staging] '\n")),(0,n.kt)("p",null,"Then configure each job with the following options:"),(0,n.kt)("p",null,"Example:"),(0,n.kt)("pre",null,(0,n.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "processors": {\n    "email": {\n      "provider": "provider1",\n      "emails": "team@mycompany.com, owner@mycompany.com",\n      "onSuccess": "true"\n    }\n  }\n}\n')),(0,n.kt)("p",null,"By default the email procesor doesn't send emails on job success, the ",(0,n.kt)("inlineCode",{parentName:"p"},"onSuccess")," parameter, enables it, like in the previous example."))}m.isMDXComponent=!0}}]);
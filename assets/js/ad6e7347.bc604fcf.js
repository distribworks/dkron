"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[5502],{3905:(e,r,t)=>{t.d(r,{Zo:()=>p,kt:()=>f});var n=t(7294);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function s(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?s(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):s(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function c(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},s=Object.keys(e);for(n=0;n<s.length;n++)t=s[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(n=0;n<s.length;n++)t=s[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var i=n.createContext({}),l=function(e){var r=n.useContext(i),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},p=function(e){var r=l(e.components);return n.createElement(i.Provider,{value:r},e.children)},u="mdxType",m={inlineCode:"code",wrapper:function(e){var r=e.children;return n.createElement(n.Fragment,{},r)}},d=n.forwardRef((function(e,r){var t=e.components,o=e.mdxType,s=e.originalType,i=e.parentName,p=c(e,["components","mdxType","originalType","parentName"]),u=l(t),d=o,f=u["".concat(i,".").concat(d)]||u[d]||m[d]||s;return t?n.createElement(f,a(a({ref:r},p),{},{components:t})):n.createElement(f,a({ref:r},p))}));function f(e,r){var t=arguments,o=r&&r.mdxType;if("string"==typeof e||o){var s=t.length,a=new Array(s);a[0]=d;var c={};for(var i in r)hasOwnProperty.call(r,i)&&(c[i]=r[i]);c.originalType=e,c[u]="string"==typeof e?e:o,a[1]=c;for(var l=2;l<s;l++)a[l]=t[l];return n.createElement.apply(null,a)}return n.createElement.apply(null,t)}d.displayName="MDXCreateElement"},763:(e,r,t)=>{t.r(r),t.d(r,{assets:()=>i,contentTitle:()=>a,default:()=>u,frontMatter:()=>s,metadata:()=>c,toc:()=>l});var n=t(7462),o=(t(7294),t(3905));const s={},a="Slack processor",c={unversionedId:"pro/processors/slack",id:"pro/processors/slack",title:"Slack processor",description:"The Slack processor provides slack notifications with multiple configurations and rich format.",source:"@site/docs/pro/processors/slack.md",sourceDirName:"pro/processors",slug:"/pro/processors/slack",permalink:"/docs/pro/processors/slack",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/pro/processors/slack.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"Email processor",permalink:"/docs/pro/processors/email"},next:{title:"Upgrade from v1 to v2",permalink:"/docs/upgrading/from_v1_to_v2"}},i={},l=[],p={toc:l};function u(e){let{components:r,...s}=e;return(0,o.kt)("wrapper",(0,n.Z)({},p,s,{components:r,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"slack-processor"},"Slack processor"),(0,o.kt)("p",null,"The Slack processor provides slack notifications with multiple configurations and rich format."),(0,o.kt)("p",null,"Configuration of the slack processor is stored in a file named ",(0,o.kt)("inlineCode",{parentName:"p"},"dkron-processor-slack.yml")," in the same locations as ",(0,o.kt)("inlineCode",{parentName:"p"},"dkron.yml"),", and should include a list of teams, it can include any number of teams."),(0,o.kt)("p",null,(0,o.kt)("img",{src:t(9433).Z,width:"643",height:"229"})),(0,o.kt)("p",null,"Example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"team1:\n  webhook_url: https://hooks.slack.com/services/XXXXXXXXXXXXXXXXXXX\n  bot_name: Dkron Production\n")),(0,o.kt)("p",null,"Then configure each job with the following options:"),(0,o.kt)("p",null,"Example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "processors": {\n    "slack": {\n      "team": "team1",\n      "channel": "#cron-production",\n      "onSuccess": "true"\n    }\n  }\n}\n')),(0,o.kt)("p",null,"By default the slack procesor doesn't send notifications on job success, the ",(0,o.kt)("inlineCode",{parentName:"p"},"onSuccess")," parameter, enables it, like in the previous example."))}u.isMDXComponent=!0},9433:(e,r,t)=>{t.d(r,{Z:()=>n});const n=t.p+"assets/images/slack-c682ec1651a106f521d514f05ac8c26c.png"}}]);
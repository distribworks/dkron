"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[6335],{3905:(e,t,r)=>{r.d(t,{Zo:()=>u,kt:()=>g});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function a(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function c(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var l=n.createContext({}),s=function(e){var t=n.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):a(a({},t),e)),r},u=function(e){var t=s(e.components);return n.createElement(l.Provider,{value:t},e.children)},p="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},f=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,i=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),p=s(r),f=o,g=p["".concat(l,".").concat(f)]||p[f]||d[f]||i;return r?n.createElement(g,a(a({ref:t},u),{},{components:r})):n.createElement(g,a({ref:t},u))}));function g(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=r.length,a=new Array(i);a[0]=f;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c[p]="string"==typeof e?e:o,a[1]=c;for(var s=2;s<i;s++)a[s]=r[s];return n.createElement.apply(null,a)}return n.createElement.apply(null,r)}f.displayName="MDXCreateElement"},7456:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>l,contentTitle:()=>a,default:()=>p,frontMatter:()=>i,metadata:()=>c,toc:()=>s});var n=r(7462),o=(r(7294),r(3905));const i={title:"Clustering"},a=void 0,c={unversionedId:"usage/clustering",id:"version-v1/usage/clustering",title:"Clustering",description:"Configure a cluster",source:"@site/versioned_docs/version-v1/usage/clustering.md",sourceDirName:"usage",slug:"/usage/clustering",permalink:"/docs/v1/usage/clustering",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v1/usage/clustering.md",tags:[],version:"v1",frontMatter:{title:"Clustering"},sidebar:"tutorialSidebar",previous:{title:"Job chaining",permalink:"/docs/v1/usage/chaining"},next:{title:"Concurrency",permalink:"/docs/v1/usage/concurrency"}},l={},s=[{value:"Configure a cluster",id:"configure-a-cluster",level:2},{value:"Etcd",id:"etcd",level:3}],u={toc:s};function p(e){let{components:t,...r}=e;return(0,o.kt)("wrapper",(0,n.Z)({},u,r,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h2",{id:"configure-a-cluster"},"Configure a cluster"),(0,o.kt)("p",null,"Dkron can run in HA mode, avoiding SPOFs, this mode provides better scalability and better reliability for users that wants a high level of confidence in the cron jobs they need to run."),(0,o.kt)("p",null,"To form a cluster, server nodes need to know the address of its peers as in the following example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"# dkron.yml\njoin:\n- 10.19.3.9\n- 10.19.4.64\n- 10.19.7.215\n")),(0,o.kt)("h3",{id:"etcd"},"Etcd"),(0,o.kt)("p",null,"For a more in detail guide of clustering with etcd follow this guide: ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/clustering.md"},"https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/clustering.md")))}p.isMDXComponent=!0}}]);
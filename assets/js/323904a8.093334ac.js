"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[5344],{3905:(e,n,t)=>{t.d(n,{Zo:()=>d,kt:()=>y});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function a(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var c=r.createContext({}),s=function(e){var n=r.useContext(c),t=n;return e&&(t="function"==typeof e?e(n):a(a({},n),e)),t},d=function(e){var n=s(e.components);return r.createElement(c.Provider,{value:n},e.children)},p="mdxType",u={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},k=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,i=e.originalType,c=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),p=s(t),k=o,y=p["".concat(c,".").concat(k)]||p[k]||u[k]||i;return t?r.createElement(y,a(a({ref:n},d),{},{components:t})):r.createElement(y,a({ref:n},d))}));function y(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var i=t.length,a=new Array(i);a[0]=k;var l={};for(var c in n)hasOwnProperty.call(n,c)&&(l[c]=n[c]);l.originalType=e,l[p]="string"==typeof e?e:o,a[1]=l;for(var s=2;s<i;s++)a[s]=t[s];return r.createElement.apply(null,a)}return r.createElement.apply(null,t)}k.displayName="MDXCreateElement"},3091:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>a,default:()=>p,frontMatter:()=>i,metadata:()=>l,toc:()=>s});var r=t(7462),o=(t(7294),t(3905));const i={date:new Date("2019-03-22T00:00:00.000Z"),title:"dkron keygen",slug:"dkron_keygen",url:"/v1.2/cli/dkron_keygen/"},a=void 0,l={unversionedId:"cli/dkron_keygen",id:"version-v1/cli/dkron_keygen",title:"dkron keygen",description:"dkron keygen",source:"@site/versioned_docs/version-v1/cli/dkron_keygen.md",sourceDirName:"cli",slug:"/cli/dkron_keygen",permalink:"/docs/v1/cli/dkron_keygen",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v1/cli/dkron_keygen.md",tags:[],version:"v1",frontMatter:{date:"2019-03-22T00:00:00.000Z",title:"dkron keygen",slug:"dkron_keygen",url:"/v1.2/cli/dkron_keygen/"},sidebar:"tutorialSidebar",previous:{title:"dkron doc",permalink:"/docs/v1/cli/dkron_doc"},next:{title:"dkron leave",permalink:"/docs/v1/cli/dkron_leave"}},c={},s=[{value:"dkron keygen",id:"dkron-keygen",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 22-Mar-2019",id:"auto-generated-by-spf13cobra-on-22-mar-2019",level:6}],d={toc:s};function p(e){let{components:n,...t}=e;return(0,o.kt)("wrapper",(0,r.Z)({},d,t,{components:n,mdxType:"MDXLayout"}),(0,o.kt)("h2",{id:"dkron-keygen"},"dkron keygen"),(0,o.kt)("p",null,"Generates a new encryption key"),(0,o.kt)("h3",{id:"synopsis"},"Synopsis"),(0,o.kt)("p",null,"Generates a new encryption key that can be used to configure the\nagent to encrypt traffic. The output of this command is already\nin the proper format that the agent expects."),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"dkron keygen [flags]\n")),(0,o.kt)("h3",{id:"options"},"Options"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"  -h, --help   help for keygen\n")),(0,o.kt)("h3",{id:"options-inherited-from-parent-commands"},"Options inherited from parent commands"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"      --config string   config file path\n")),(0,o.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("a",{parentName:"li",href:"/docs/v1/cli/dkron/"},"dkron"),"\t - Open source distributed job scheduling system")),(0,o.kt)("h6",{id:"auto-generated-by-spf13cobra-on-22-mar-2019"},"Auto generated by spf13/cobra on 22-Mar-2019"))}p.isMDXComponent=!0}}]);
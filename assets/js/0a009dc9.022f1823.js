"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[3938],{3905:(e,n,r)=>{r.d(n,{Zo:()=>s,kt:()=>y});var t=r(7294);function o(e,n,r){return n in e?Object.defineProperty(e,n,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[n]=r,e}function i(e,n){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);n&&(t=t.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),r.push.apply(r,t)}return r}function a(e){for(var n=1;n<arguments.length;n++){var r=null!=arguments[n]?arguments[n]:{};n%2?i(Object(r),!0).forEach((function(n){o(e,n,r[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(r,n))}))}return e}function l(e,n){if(null==e)return{};var r,t,o=function(e,n){if(null==e)return{};var r,t,o={},i=Object.keys(e);for(t=0;t<i.length;t++)r=i[t],n.indexOf(r)>=0||(o[r]=e[r]);return o}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(t=0;t<i.length;t++)r=i[t],n.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var c=t.createContext({}),p=function(e){var n=t.useContext(c),r=n;return e&&(r="function"==typeof e?e(n):a(a({},n),e)),r},s=function(e){var n=p(e.components);return t.createElement(c.Provider,{value:n},e.children)},d="mdxType",u={inlineCode:"code",wrapper:function(e){var n=e.children;return t.createElement(t.Fragment,{},n)}},k=t.forwardRef((function(e,n){var r=e.components,o=e.mdxType,i=e.originalType,c=e.parentName,s=l(e,["components","mdxType","originalType","parentName"]),d=p(r),k=o,y=d["".concat(c,".").concat(k)]||d[k]||u[k]||i;return r?t.createElement(y,a(a({ref:n},s),{},{components:r})):t.createElement(y,a({ref:n},s))}));function y(e,n){var r=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var i=r.length,a=new Array(i);a[0]=k;var l={};for(var c in n)hasOwnProperty.call(n,c)&&(l[c]=n[c]);l.originalType=e,l[d]="string"==typeof e?e:o,a[1]=l;for(var p=2;p<i;p++)a[p]=r[p];return t.createElement.apply(null,a)}return t.createElement.apply(null,r)}k.displayName="MDXCreateElement"},9135:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>c,contentTitle:()=>a,default:()=>d,frontMatter:()=>i,metadata:()=>l,toc:()=>p});var t=r(7462),o=(r(7294),r(3905));const i={date:new Date("2019-08-26T00:00:00.000Z"),title:"dkron keygen",slug:"dkron_keygen",url:"/2.0/pro/cli/dkron_keygen/"},a=void 0,l={unversionedId:"pro/cli/dkron_keygen",id:"version-v2/pro/cli/dkron_keygen",title:"dkron keygen",description:"dkron keygen",source:"@site/versioned_docs/version-v2/pro/cli/dkron_keygen.md",sourceDirName:"pro/cli",slug:"/pro/cli/dkron_keygen",permalink:"/docs/v2/pro/cli/dkron_keygen",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v2/pro/cli/dkron_keygen.md",tags:[],version:"v2",frontMatter:{date:"2019-08-26T00:00:00.000Z",title:"dkron keygen",slug:"dkron_keygen",url:"/2.0/pro/cli/dkron_keygen/"},sidebar:"tutorialSidebar",previous:{title:"dkron doc",permalink:"/docs/v2/pro/cli/dkron_doc"},next:{title:"dkron leave",permalink:"/docs/v2/pro/cli/dkron_leave"}},c={},p=[{value:"dkron keygen",id:"dkron-keygen",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 26-Aug-2019",id:"auto-generated-by-spf13cobra-on-26-aug-2019",level:6}],s={toc:p};function d(e){let{components:n,...r}=e;return(0,o.kt)("wrapper",(0,t.Z)({},s,r,{components:n,mdxType:"MDXLayout"}),(0,o.kt)("h2",{id:"dkron-keygen"},"dkron keygen"),(0,o.kt)("p",null,"Generates a new encryption key"),(0,o.kt)("h3",{id:"synopsis"},"Synopsis"),(0,o.kt)("p",null,"Generates a new encryption key that can be used to configure the\nagent to encrypt traffic. The output of this command is already\nin the proper format that the agent expects."),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"dkron keygen [flags]\n")),(0,o.kt)("h3",{id:"options"},"Options"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"  -h, --help   help for keygen\n")),(0,o.kt)("h3",{id:"options-inherited-from-parent-commands"},"Options inherited from parent commands"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"      --config string   config file (default is /etc/dkron/dkron.yml)\n")),(0,o.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("a",{parentName:"li",href:"/docs/v2/cli/dkron/"},"dkron"),"\t - Professional distributed job scheduling system")),(0,o.kt)("h6",{id:"auto-generated-by-spf13cobra-on-26-aug-2019"},"Auto generated by spf13/cobra on 26-Aug-2019"))}d.isMDXComponent=!0}}]);
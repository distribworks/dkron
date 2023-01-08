"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[650],{3905:(e,t,o)=>{o.d(t,{Zo:()=>s,kt:()=>f});var n=o(7294);function r(e,t,o){return t in e?Object.defineProperty(e,t,{value:o,enumerable:!0,configurable:!0,writable:!0}):e[t]=o,e}function i(e,t){var o=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),o.push.apply(o,n)}return o}function l(e){for(var t=1;t<arguments.length;t++){var o=null!=arguments[t]?arguments[t]:{};t%2?i(Object(o),!0).forEach((function(t){r(e,t,o[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(o)):i(Object(o)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(o,t))}))}return e}function a(e,t){if(null==e)return{};var o,n,r=function(e,t){if(null==e)return{};var o,n,r={},i=Object.keys(e);for(n=0;n<i.length;n++)o=i[n],t.indexOf(o)>=0||(r[o]=e[o]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)o=i[n],t.indexOf(o)>=0||Object.prototype.propertyIsEnumerable.call(e,o)&&(r[o]=e[o])}return r}var c=n.createContext({}),p=function(e){var t=n.useContext(c),o=t;return e&&(o="function"==typeof e?e(t):l(l({},t),e)),o},s=function(e){var t=p(e.components);return n.createElement(c.Provider,{value:t},e.children)},d="mdxType",m={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},u=n.forwardRef((function(e,t){var o=e.components,r=e.mdxType,i=e.originalType,c=e.parentName,s=a(e,["components","mdxType","originalType","parentName"]),d=p(o),u=r,f=d["".concat(c,".").concat(u)]||d[u]||m[u]||i;return o?n.createElement(f,l(l({ref:t},s),{},{components:o})):n.createElement(f,l({ref:t},s))}));function f(e,t){var o=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=o.length,l=new Array(i);l[0]=u;var a={};for(var c in t)hasOwnProperty.call(t,c)&&(a[c]=t[c]);a.originalType=e,a[d]="string"==typeof e?e:r,l[1]=a;for(var p=2;p<i;p++)l[p]=o[p];return n.createElement.apply(null,l)}return n.createElement.apply(null,o)}u.displayName="MDXCreateElement"},8210:(e,t,o)=>{o.r(t),o.d(t,{assets:()=>c,contentTitle:()=>l,default:()=>d,frontMatter:()=>i,metadata:()=>a,toc:()=>p});var n=o(7462),r=(o(7294),o(3905));const i={date:new Date("2022-06-05T00:00:00.000Z"),title:"dkron completion",slug:"dkron_completion",url:"/cli/dkron_completion/"},l=void 0,a={unversionedId:"cli/dkron_completion",id:"cli/dkron_completion",title:"dkron completion",description:"dkron completion",source:"@site/docs/cli/dkron_completion.md",sourceDirName:"cli",slug:"/cli/dkron_completion",permalink:"/docs/cli/dkron_completion",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/cli/dkron_completion.md",tags:[],version:"current",frontMatter:{date:"2022-06-05T00:00:00.000Z",title:"dkron completion",slug:"dkron_completion",url:"/cli/dkron_completion/"},sidebar:"tutorialSidebar",previous:{title:"dkron agent",permalink:"/docs/cli/dkron_agent"},next:{title:"dkron completion bash",permalink:"/docs/cli/dkron_completion_bash"}},c={},p=[{value:"dkron completion",id:"dkron-completion",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 5-Jun-2022",id:"auto-generated-by-spf13cobra-on-5-jun-2022",level:6}],s={toc:p};function d(e){let{components:t,...o}=e;return(0,r.kt)("wrapper",(0,n.Z)({},s,o,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"dkron-completion"},"dkron completion"),(0,r.kt)("p",null,"Generate the autocompletion script for the specified shell"),(0,r.kt)("h3",{id:"synopsis"},"Synopsis"),(0,r.kt)("p",null,"Generate the autocompletion script for dkron for the specified shell.\nSee each sub-command's help for details on how to use the generated script."),(0,r.kt)("h3",{id:"options"},"Options"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"  -h, --help   help for completion\n")),(0,r.kt)("h3",{id:"options-inherited-from-parent-commands"},"Options inherited from parent commands"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"      --config string   config file path\n")),(0,r.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"/docs/cli/dkron/"},"dkron"),"\t - Open source distributed job scheduling system"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"/docs/cli/dkron_completion_bash/"},"dkron completion bash"),"\t - Generate the autocompletion script for bash"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"/docs/cli/dkron_completion_fish/"},"dkron completion fish"),"\t - Generate the autocompletion script for fish"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"/docs/cli/dkron_completion_powershell/"},"dkron completion powershell"),"\t - Generate the autocompletion script for powershell"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"/docs/cli/dkron_completion_zsh/"},"dkron completion zsh"),"\t - Generate the autocompletion script for zsh")),(0,r.kt)("h6",{id:"auto-generated-by-spf13cobra-on-5-jun-2022"},"Auto generated by spf13/cobra on 5-Jun-2022"))}d.isMDXComponent=!0}}]);
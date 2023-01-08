"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[1660],{3905:(e,o,t)=>{t.d(o,{Zo:()=>a,kt:()=>f});var n=t(7294);function r(e,o,t){return o in e?Object.defineProperty(e,o,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[o]=t,e}function l(e,o){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);o&&(n=n.filter((function(o){return Object.getOwnPropertyDescriptor(e,o).enumerable}))),t.push.apply(t,n)}return t}function i(e){for(var o=1;o<arguments.length;o++){var t=null!=arguments[o]?arguments[o]:{};o%2?l(Object(t),!0).forEach((function(o){r(e,o,t[o])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):l(Object(t)).forEach((function(o){Object.defineProperty(e,o,Object.getOwnPropertyDescriptor(t,o))}))}return e}function p(e,o){if(null==e)return{};var t,n,r=function(e,o){if(null==e)return{};var t,n,r={},l=Object.keys(e);for(n=0;n<l.length;n++)t=l[n],o.indexOf(t)>=0||(r[t]=e[t]);return r}(e,o);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(n=0;n<l.length;n++)t=l[n],o.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(r[t]=e[t])}return r}var c=n.createContext({}),s=function(e){var o=n.useContext(c),t=o;return e&&(t="function"==typeof e?e(o):i(i({},o),e)),t},a=function(e){var o=s(e.components);return n.createElement(c.Provider,{value:o},e.children)},d="mdxType",u={inlineCode:"code",wrapper:function(e){var o=e.children;return n.createElement(n.Fragment,{},o)}},m=n.forwardRef((function(e,o){var t=e.components,r=e.mdxType,l=e.originalType,c=e.parentName,a=p(e,["components","mdxType","originalType","parentName"]),d=s(t),m=r,f=d["".concat(c,".").concat(m)]||d[m]||u[m]||l;return t?n.createElement(f,i(i({ref:o},a),{},{components:t})):n.createElement(f,i({ref:o},a))}));function f(e,o){var t=arguments,r=o&&o.mdxType;if("string"==typeof e||r){var l=t.length,i=new Array(l);i[0]=m;var p={};for(var c in o)hasOwnProperty.call(o,c)&&(p[c]=o[c]);p.originalType=e,p[d]="string"==typeof e?e:r,i[1]=p;for(var s=2;s<l;s++)i[s]=t[s];return n.createElement.apply(null,i)}return n.createElement.apply(null,t)}m.displayName="MDXCreateElement"},1127:(e,o,t)=>{t.r(o),t.d(o,{assets:()=>c,contentTitle:()=>i,default:()=>d,frontMatter:()=>l,metadata:()=>p,toc:()=>s});var n=t(7462),r=(t(7294),t(3905));const l={date:new Date("2022-06-05T00:00:00.000Z"),title:"dkron completion powershell",slug:"dkron_completion_powershell",url:"/cli/dkron_completion_powershell/"},i=void 0,p={unversionedId:"cli/dkron_completion_powershell",id:"cli/dkron_completion_powershell",title:"dkron completion powershell",description:"dkron completion powershell",source:"@site/docs/cli/dkron_completion_powershell.md",sourceDirName:"cli",slug:"/cli/dkron_completion_powershell",permalink:"/docs/cli/dkron_completion_powershell",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/cli/dkron_completion_powershell.md",tags:[],version:"current",frontMatter:{date:"2022-06-05T00:00:00.000Z",title:"dkron completion powershell",slug:"dkron_completion_powershell",url:"/cli/dkron_completion_powershell/"},sidebar:"tutorialSidebar",previous:{title:"dkron completion fish",permalink:"/docs/cli/dkron_completion_fish"},next:{title:"dkron completion zsh",permalink:"/docs/cli/dkron_completion_zsh"}},c={},s=[{value:"dkron completion powershell",id:"dkron-completion-powershell",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 5-Jun-2022",id:"auto-generated-by-spf13cobra-on-5-jun-2022",level:6}],a={toc:s};function d(e){let{components:o,...t}=e;return(0,r.kt)("wrapper",(0,n.Z)({},a,t,{components:o,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"dkron-completion-powershell"},"dkron completion powershell"),(0,r.kt)("p",null,"Generate the autocompletion script for powershell"),(0,r.kt)("h3",{id:"synopsis"},"Synopsis"),(0,r.kt)("p",null,"Generate the autocompletion script for powershell."),(0,r.kt)("p",null,"To load completions in your current shell session:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"dkron completion powershell | Out-String | Invoke-Expression\n")),(0,r.kt)("p",null,"To load completions for every new session, add the output of the above command\nto your powershell profile."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"dkron completion powershell [flags]\n")),(0,r.kt)("h3",{id:"options"},"Options"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"  -h, --help              help for powershell\n      --no-descriptions   disable completion descriptions\n")),(0,r.kt)("h3",{id:"options-inherited-from-parent-commands"},"Options inherited from parent commands"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"      --config string   config file path\n")),(0,r.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("a",{parentName:"li",href:"/docs/cli/dkron_completion/"},"dkron completion"),"\t - Generate the autocompletion script for the specified shell")),(0,r.kt)("h6",{id:"auto-generated-by-spf13cobra-on-5-jun-2022"},"Auto generated by spf13/cobra on 5-Jun-2022"))}d.isMDXComponent=!0}}]);
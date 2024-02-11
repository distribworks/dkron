"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[512],{80616:(e,n,o)=>{o.r(n),o.d(n,{assets:()=>l,contentTitle:()=>s,default:()=>p,frontMatter:()=>i,metadata:()=>c,toc:()=>a});var r=o(17624),t=o(95788);const i={date:new Date("2023-09-02T00:00:00.000Z"),title:"dkron completion bash",slug:"dkron_completion_bash",url:"/docs/pro/cli/dkron_completion_bash/"},s=void 0,c={id:"pro/cli/dkron_completion_bash",title:"dkron completion bash",description:"dkron completion bash",source:"@site/versioned_docs/version-v3/pro/cli/dkron_completion_bash.md",sourceDirName:"pro/cli",slug:"/pro/cli/dkron_completion_bash",permalink:"/docs/v3/pro/cli/dkron_completion_bash",draft:!1,unlisted:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v3/pro/cli/dkron_completion_bash.md",tags:[],version:"v3",frontMatter:{date:"2023-09-02T00:00:00.000Z",title:"dkron completion bash",slug:"dkron_completion_bash",url:"/docs/pro/cli/dkron_completion_bash/"},sidebar:"tutorialSidebar",previous:{title:"dkron completion",permalink:"/docs/v3/pro/cli/dkron_completion"},next:{title:"dkron completion fish",permalink:"/docs/v3/pro/cli/dkron_completion_fish"}},l={},a=[{value:"dkron completion bash",id:"dkron-completion-bash",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Linux:",id:"linux",level:4},{value:"macOS:",id:"macos",level:4},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 2-Sep-2023",id:"auto-generated-by-spf13cobra-on-2-sep-2023",level:6}];function d(e){const n={a:"a",code:"code",h2:"h2",h3:"h3",h4:"h4",h6:"h6",li:"li",p:"p",pre:"pre",ul:"ul",...(0,t.MN)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(n.h2,{id:"dkron-completion-bash",children:"dkron completion bash"}),"\n",(0,r.jsx)(n.p,{children:"Generate the autocompletion script for bash"}),"\n",(0,r.jsx)(n.h3,{id:"synopsis",children:"Synopsis"}),"\n",(0,r.jsx)(n.p,{children:"Generate the autocompletion script for the bash shell."}),"\n",(0,r.jsx)(n.p,{children:"This script depends on the 'bash-completion' package.\nIf it is not installed already, you can install it via your OS's package manager."}),"\n",(0,r.jsx)(n.p,{children:"To load completions in your current shell session:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"\tsource <(dkron completion bash)\n"})}),"\n",(0,r.jsx)(n.p,{children:"To load completions for every new session, execute once:"}),"\n",(0,r.jsx)(n.h4,{id:"linux",children:"Linux:"}),"\n",(0,r.jsx)(n.p,{children:"dkron completion bash > /etc/bash_completion.d/dkron"}),"\n",(0,r.jsx)(n.h4,{id:"macos",children:"macOS:"}),"\n",(0,r.jsx)(n.p,{children:"dkron completion bash > $(brew --prefix)/etc/bash_completion.d/dkron"}),"\n",(0,r.jsx)(n.p,{children:"You will need to start a new shell for this setup to take effect."}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"dkron completion bash\n"})}),"\n",(0,r.jsx)(n.h3,{id:"options",children:"Options"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"  -h, --help              help for bash\n      --no-descriptions   disable completion descriptions\n"})}),"\n",(0,r.jsx)(n.h3,{id:"options-inherited-from-parent-commands",children:"Options inherited from parent commands"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"      --config string   config file (default is /etc/dkron/dkron.yml)\n"})}),"\n",(0,r.jsx)(n.h3,{id:"see-also",children:"SEE ALSO"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.a,{href:"/docs/pro/cli/dkron_completion/",children:"dkron completion"}),"\t - Generate the autocompletion script for the specified shell"]}),"\n"]}),"\n",(0,r.jsx)(n.h6,{id:"auto-generated-by-spf13cobra-on-2-sep-2023",children:"Auto generated by spf13/cobra on 2-Sep-2023"})]})}function p(e={}){const{wrapper:n}={...(0,t.MN)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(d,{...e})}):d(e)}},95788:(e,n,o)=>{o.d(n,{MN:()=>a});var r=o(11504);function t(e,n,o){return n in e?Object.defineProperty(e,n,{value:o,enumerable:!0,configurable:!0,writable:!0}):e[n]=o,e}function i(e,n){var o=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),o.push.apply(o,r)}return o}function s(e){for(var n=1;n<arguments.length;n++){var o=null!=arguments[n]?arguments[n]:{};n%2?i(Object(o),!0).forEach((function(n){t(e,n,o[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(o)):i(Object(o)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(o,n))}))}return e}function c(e,n){if(null==e)return{};var o,r,t=function(e,n){if(null==e)return{};var o,r,t={},i=Object.keys(e);for(r=0;r<i.length;r++)o=i[r],n.indexOf(o)>=0||(t[o]=e[o]);return t}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)o=i[r],n.indexOf(o)>=0||Object.prototype.propertyIsEnumerable.call(e,o)&&(t[o]=e[o])}return t}var l=r.createContext({}),a=function(e){var n=r.useContext(l),o=n;return e&&(o="function"==typeof e?e(n):s(s({},n),e)),o},d={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},p=r.forwardRef((function(e,n){var o=e.components,t=e.mdxType,i=e.originalType,l=e.parentName,p=c(e,["components","mdxType","originalType","parentName"]),h=a(o),m=t,u=h["".concat(l,".").concat(m)]||h[m]||d[m]||i;return o?r.createElement(u,s(s({ref:n},p),{},{components:o})):r.createElement(u,s({ref:n},p))}));p.displayName="MDXCreateElement"}}]);
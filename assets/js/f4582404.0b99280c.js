"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[6935],{3905:(e,r,t)=>{t.d(r,{Zo:()=>s,kt:()=>u});var n=t(7294);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function a(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function i(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?a(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):a(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function p(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)t=a[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)t=a[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var l=n.createContext({}),c=function(e){var r=n.useContext(l),t=r;return e&&(t="function"==typeof e?e(r):i(i({},r),e)),t},s=function(e){var r=c(e.components);return n.createElement(l.Provider,{value:r},e.children)},d="mdxType",f={inlineCode:"code",wrapper:function(e){var r=e.children;return n.createElement(n.Fragment,{},r)}},m=n.forwardRef((function(e,r){var t=e.components,o=e.mdxType,a=e.originalType,l=e.parentName,s=p(e,["components","mdxType","originalType","parentName"]),d=c(t),m=o,u=d["".concat(l,".").concat(m)]||d[m]||f[m]||a;return t?n.createElement(u,i(i({ref:r},s),{},{components:t})):n.createElement(u,i({ref:r},s))}));function u(e,r){var t=arguments,o=r&&r.mdxType;if("string"==typeof e||o){var a=t.length,i=new Array(a);i[0]=m;var p={};for(var l in r)hasOwnProperty.call(r,l)&&(p[l]=r[l]);p.originalType=e,p[d]="string"==typeof e?e:o,i[1]=p;for(var c=2;c<a;c++)i[c]=t[c];return n.createElement.apply(null,i)}return n.createElement.apply(null,t)}m.displayName="MDXCreateElement"},4805:(e,r,t)=>{t.r(r),t.d(r,{assets:()=>l,contentTitle:()=>i,default:()=>d,frontMatter:()=>a,metadata:()=>p,toc:()=>c});var n=t(7462),o=(t(7294),t(3905));const a={date:new Date("2022-10-07T00:00:00.000Z"),title:"dkron raft remove-peer",slug:"dkron_raft_remove-peer",url:"/docs/pro/cli/dkron_raft_remove-peer/"},i=void 0,p={unversionedId:"pro/cli/dkron_raft_remove-peer",id:"pro/cli/dkron_raft_remove-peer",title:"dkron raft remove-peer",description:"dkron raft remove-peer",source:"@site/docs/pro/cli/dkron_raft_remove-peer.md",sourceDirName:"pro/cli",slug:"/pro/cli/dkron_raft_remove-peer",permalink:"/docs/pro/cli/dkron_raft_remove-peer",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/pro/cli/dkron_raft_remove-peer.md",tags:[],version:"current",frontMatter:{date:"2022-10-07T00:00:00.000Z",title:"dkron raft remove-peer",slug:"dkron_raft_remove-peer",url:"/docs/pro/cli/dkron_raft_remove-peer/"},sidebar:"tutorialSidebar",previous:{title:"dkron raft list-peers",permalink:"/docs/pro/cli/dkron_raft_list-peers"},next:{title:"dkron version",permalink:"/docs/pro/cli/dkron_version"}},l={},c=[{value:"dkron raft remove-peer",id:"dkron-raft-remove-peer",level:2},{value:"Options",id:"options",level:3},{value:"Options inherited from parent commands",id:"options-inherited-from-parent-commands",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 7-Oct-2022",id:"auto-generated-by-spf13cobra-on-7-oct-2022",level:6}],s={toc:c};function d(e){let{components:r,...t}=e;return(0,o.kt)("wrapper",(0,n.Z)({},s,t,{components:r,mdxType:"MDXLayout"}),(0,o.kt)("h2",{id:"dkron-raft-remove-peer"},"dkron raft remove-peer"),(0,o.kt)("p",null,"Command to list raft peers"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"dkron raft remove-peer [flags]\n")),(0,o.kt)("h3",{id:"options"},"Options"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"  -h, --help             help for remove-peer\n      --peer-id string   Remove a Dkron server with the given ID from the Raft configuration.\n")),(0,o.kt)("h3",{id:"options-inherited-from-parent-commands"},"Options inherited from parent commands"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},'      --cert-file string         Path to the client server TLS cert file\n      --config string            config file (default is /etc/dkron/dkron.yml)\n      --key-file string          Path to the client server TLS key file\n      --rpc-addr string          gRPC address of the agent. (default "{{ GetPrivateIP }}:6868")\n      --trusted-ca-file string   Path to the client server TLS trusted CA cert file\n')),(0,o.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_raft/"},"dkron raft"),"\t - Command to perform some raft operations")),(0,o.kt)("h6",{id:"auto-generated-by-spf13cobra-on-7-oct-2022"},"Auto generated by spf13/cobra on 7-Oct-2022"))}d.isMDXComponent=!0}}]);
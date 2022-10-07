"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[8845],{3905:function(e,r,t){t.d(r,{Zo:function(){return d},kt:function(){return f}});var n=t(67294);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function i(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?i(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function l(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},i=Object.keys(e);for(n=0;n<i.length;n++)t=i[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)t=i[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var c=n.createContext({}),s=function(e){var r=n.useContext(c),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},d=function(e){var r=s(e.components);return n.createElement(c.Provider,{value:r},e.children)},u={inlineCode:"code",wrapper:function(e){var r=e.children;return n.createElement(n.Fragment,{},r)}},p=n.forwardRef((function(e,r){var t=e.components,o=e.mdxType,i=e.originalType,c=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),p=s(t),f=o,k=p["".concat(c,".").concat(f)]||p[f]||u[f]||i;return t?n.createElement(k,a(a({ref:r},d),{},{components:t})):n.createElement(k,a({ref:r},d))}));function f(e,r){var t=arguments,o=r&&r.mdxType;if("string"==typeof e||o){var i=t.length,a=new Array(i);a[0]=p;var l={};for(var c in r)hasOwnProperty.call(r,c)&&(l[c]=r[c]);l.originalType=e,l.mdxType="string"==typeof e?e:o,a[1]=l;for(var s=2;s<i;s++)a[s]=t[s];return n.createElement.apply(null,a)}return n.createElement.apply(null,t)}p.displayName="MDXCreateElement"},53706:function(e,r,t){t.r(r),t.d(r,{assets:function(){return d},contentTitle:function(){return c},default:function(){return f},frontMatter:function(){return l},metadata:function(){return s},toc:function(){return u}});var n=t(87462),o=t(63366),i=(t(67294),t(3905)),a=["components"],l={date:new Date("2022-10-07T00:00:00.000Z"),title:"dkron",slug:"dkron",url:"/docs/pro/cli/dkron/"},c=void 0,s={unversionedId:"pro/cli/dkron",id:"pro/cli/dkron",title:"dkron",description:"dkron",source:"@site/docs/pro/cli/dkron.md",sourceDirName:"pro/cli",slug:"/pro/cli/dkron",permalink:"/docs/pro/cli/dkron",editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/pro/cli/dkron.md",tags:[],version:"current",frontMatter:{date:"2022-10-07T00:00:00.000Z",title:"dkron",slug:"dkron",url:"/docs/pro/cli/dkron/"},sidebar:"tutorialSidebar",previous:{title:"Authentication",permalink:"/docs/pro/auth"},next:{title:"dkron agent",permalink:"/docs/pro/cli/dkron_agent"}},d={},u=[{value:"dkron",id:"dkron",level:2},{value:"Synopsis",id:"synopsis",level:3},{value:"Options",id:"options",level:3},{value:"SEE ALSO",id:"see-also",level:3},{value:"Auto generated by spf13/cobra on 7-Oct-2022",id:"auto-generated-by-spf13cobra-on-7-oct-2022",level:6}],p={toc:u};function f(e){var r=e.components,t=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,n.Z)({},p,t,{components:r,mdxType:"MDXLayout"}),(0,i.kt)("h2",{id:"dkron"},"dkron"),(0,i.kt)("p",null,"Professional distributed job scheduling system"),(0,i.kt)("h3",{id:"synopsis"},"Synopsis"),(0,i.kt)("p",null,"Dkron is a system service that runs scheduled jobs at given intervals or times,\njust like the cron unix service but distributed in several machines in a cluster.\nIf a machine fails (the leader), a follower will take over and keep running the scheduled jobs without human intervention."),(0,i.kt)("h3",{id:"options"},"Options"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"      --config string   config file (default is /etc/dkron/dkron.yml)\n  -h, --help            help for dkron\n")),(0,i.kt)("h3",{id:"see-also"},"SEE ALSO"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_agent/"},"dkron agent"),"\t - Start a dkron agent"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_doc/"},"dkron doc"),"\t - Generate Markdown documentation for the Dkron CLI."),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_keygen/"},"dkron keygen"),"\t - Generates a new encryption key"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_leave/"},"dkron leave"),"\t - Force an agent to leave the cluster"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_raft/"},"dkron raft"),"\t - Command to perform some raft operations"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/docs/pro/cli/dkron_version/"},"dkron version"),"\t - Show version")),(0,i.kt)("h6",{id:"auto-generated-by-spf13cobra-on-7-oct-2022"},"Auto generated by spf13/cobra on 7-Oct-2022"))}f.isMDXComponent=!0}}]);
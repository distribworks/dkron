"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[5423],{3905:(e,t,n)=>{n.d(t,{Zo:()=>u,kt:()=>m});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function a(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function s(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?a(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):a(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var c=r.createContext({}),l=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):s(s({},t),e)),n},u=function(e){var t=l(e.components);return r.createElement(c.Provider,{value:t},e.children)},p="mdxType",f={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,a=e.originalType,c=e.parentName,u=i(e,["components","mdxType","originalType","parentName"]),p=l(n),d=o,m=p["".concat(c,".").concat(d)]||p[d]||f[d]||a;return n?r.createElement(m,s(s({ref:t},u),{},{components:n})):r.createElement(m,s({ref:t},u))}));function m(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var a=n.length,s=new Array(a);s[0]=d;var i={};for(var c in t)hasOwnProperty.call(t,c)&&(i[c]=t[c]);i.originalType=e,i[p]="string"==typeof e?e:o,s[1]=i;for(var l=2;l<a;l++)s[l]=n[l];return r.createElement.apply(null,s)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},3333:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>s,default:()=>p,frontMatter:()=>a,metadata:()=>i,toc:()=>l});var r=n(7462),o=(n(7294),n(3905));const a={title:"AWS ECS Executor"},s=void 0,i={unversionedId:"pro/executors/ecs",id:"version-v1/pro/executors/ecs",title:"AWS ECS Executor",description:"The ECS exeutor is capable of launching tasks in ECS clusters, then listen to a stream of CloudWatch Logs and return the output.",source:"@site/versioned_docs/version-v1/pro/executors/ecs.md",sourceDirName:"pro/executors",slug:"/pro/executors/ecs",permalink:"/docs/v1/pro/executors/ecs",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/versioned_docs/version-v1/pro/executors/ecs.md",tags:[],version:"v1",frontMatter:{title:"AWS ECS Executor"},sidebar:"tutorialSidebar",previous:{title:"Docker executor",permalink:"/docs/v1/pro/executors/docker"},next:{title:"Elasticsearch processor",permalink:"/docs/v1/pro/processors/elasticsearch"}},c={},l=[],u={toc:l};function p(e){let{components:t,...n}=e;return(0,o.kt)("wrapper",(0,r.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"The ECS exeutor is capable of launching tasks in ECS clusters, then listen to a stream of CloudWatch Logs and return the output."),(0,o.kt)("p",null,"To configure a job to be run in ECS, the executor needs a JSON Task definition template or an already defined task in ECS."),(0,o.kt)("p",null,"To allow the ECS Task runner to run tasks, the machine running Dkron needs to have the appropriate permissions configured in AWS IAM:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "Version": "2012-10-17",\n    "Statement": [\n        {\n            "Sid": "Stmt1460720941000",\n            "Effect": "Allow",\n            "Action": [\n                "ecs:RunTask",\n                "ecs:DescribeTasks",\n                "ecs:DescribeTaskDefinition",\n                "logs:FilterLogEvents",\n                "logs:DescribeLogStreams",\n                "logs:PutLogEvents"\n            ],\n            "Resource": [\n                "*"\n            ]\n        }\n    ]\n}\n')),(0,o.kt)("p",null,"To configure a job to be run with the ECS executor:"),(0,o.kt)("p",null,"Example using an existing taskdef"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "executor": "ecs",\n  "executor_config": {\n    "taskdefName": "mytaskdef-family",\n    "region": "eu-west-1",\n    "cluster": "default",\n    "env": "ENVIRONMENT=variable",\n    "service": "mycontainer",\n    "overrides": "echo,\\"Hello from dkron\\""\n  }\n}\n')),(0,o.kt)("p",null,"Example using a provided taskdef"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "executor": "ecs",\n  "executor_config": {\n    "taskdefBody": "{\\"containerDefinitions\\": [{\\"essential\\": true,\\"image\\": \\"hello-world\\",\\"memory\\": 100,\\"name\\": \\"hello-world\\"}],\\"family\\": \\"helloworld\\"}",\n    "region": "eu-west-1",\n    "cluster": "default",\n    "fargate": "yes",\n    "env": "ENVIRONMENT=variable",\n    "maxAttempts": "5000"\n  }\n}\n')),(0,o.kt)("p",null,"This is the complete list of configuration parameters of the plugin:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"taskdefBody\ntaskdefName\nregion\ncluster\nlogGroup\nfargate\nsecurityGroup\nsubnet\nenv\nservice\noverrides\nmaxAttempts // Defaults to 2000, will perform a check every 6s * 2000 times waiting a total of 12000s or 3.3h\n")))}p.isMDXComponent=!0}}]);
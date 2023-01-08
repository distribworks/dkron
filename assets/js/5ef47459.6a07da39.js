"use strict";(self.webpackChunkmy_website=self.webpackChunkmy_website||[]).push([[3194],{3905:(e,t,n)=>{n.d(t,{Zo:()=>c,kt:()=>h});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=r.createContext({}),u=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},c=function(e){var t=u(e.components);return r.createElement(l.Provider,{value:t},e.children)},p="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,l=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),p=u(n),m=a,h=p["".concat(l,".").concat(m)]||p[m]||d[m]||i;return n?r.createElement(h,o(o({ref:t},c),{},{components:n})):r.createElement(h,o({ref:t},c))}));function h(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,o=new Array(i);o[0]=m;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s[p]="string"==typeof e?e:a,o[1]=s;for(var u=2;u<i;u++)o[u]=n[u];return r.createElement.apply(null,o)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},3882:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>o,default:()=>p,frontMatter:()=>i,metadata:()=>s,toc:()=>u});var r=n(7462),a=(n(7294),n(3905));const i={title:"Cron spec",weight:20},o=void 0,s={unversionedId:"usage/cron-spec",id:"usage/cron-spec",title:"Cron spec",description:"CRON Expression Format",source:"@site/docs/usage/cron-spec.md",sourceDirName:"usage",slug:"/usage/cron-spec",permalink:"/docs/usage/cron-spec",draft:!1,editUrl:"https://github.com/distribworks/dkron/tree/main/website/docs/docs/usage/cron-spec.md",tags:[],version:"current",frontMatter:{title:"Cron spec",weight:20},sidebar:"tutorialSidebar",previous:{title:"Concurrency",permalink:"/docs/usage/concurrency"},next:{title:"Cronitor Integration",permalink:"/docs/usage/cronitor"}},l={},u=[{value:"CRON Expression Format",id:"cron-expression-format",level:2},{value:"Predefined schedules",id:"predefined-schedules",level:3},{value:"Intervals",id:"intervals",level:3},{value:"Fixed times",id:"fixed-times",level:3},{value:"Time zones",id:"time-zones",level:3}],c={toc:u};function p(e){let{components:t,...n}=e;return(0,a.kt)("wrapper",(0,r.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"cron-expression-format"},"CRON Expression Format"),(0,a.kt)("p",null,"A cron expression represents a set of times, using 6 space-separated fields."),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"Field name   | Mandatory? | Allowed values  | Allowed special characters\n----------   | ---------- | --------------  | --------------------------\nSeconds      | Yes        | 0-59            | * / , -\nMinutes      | Yes        | 0-59            | * / , -\nHours        | Yes        | 0-23            | * / , -\nDay of month | Yes        | 1-31            | * / , - ?\nMonth        | Yes        | 1-12 or JAN-DEC | * / , -\nDay of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?\n")),(0,a.kt)("p",null,'Note: Month and Day-of-week field values are case insensitive.  "SUN", "Sun",\nand "sun" are equally accepted.'),(0,a.kt)("p",null,"Special Characters"),(0,a.kt)("p",null,"Asterisk ( * )"),(0,a.kt)("p",null,"The asterisk indicates that the cron expression will match for all values of the\nfield; e.g., using an asterisk in the 5th field (month) would indicate every\nmonth."),(0,a.kt)("p",null,"Slash ( / )"),(0,a.kt)("p",null,'Slashes are used to describe increments of ranges. For example 3-59/15 in the\n1st field (seconds) would indicate the 3rd second of the minute and every 15\nseconds thereafter. The form "*\\/..." is equivalent to the form "first-last/...",\nthat is, an increment over the largest possible range of the field.  The form\n"N/..." is accepted as meaning "N-MAX/...", that is, starting at N, use the\nincrement until the end of that specific range.  It does not wrap around.'),(0,a.kt)("p",null,"Comma ( , )"),(0,a.kt)("p",null,'Commas are used to separate items of a list. For example, using "MON,WED,FRI" in\nthe 6th field (day of week) would mean Mondays, Wednesdays and Fridays.'),(0,a.kt)("p",null,"Hyphen ( - )"),(0,a.kt)("p",null,"Hyphens are used to define ranges. For example, 9-17 would indicate every\nhour between 9am and 5pm inclusive."),(0,a.kt)("p",null,"Question mark ( ? )"),(0,a.kt)("p",null,"Question mark may be used instead of '*' for leaving either day-of-month or\nday-of-week blank."),(0,a.kt)("h3",{id:"predefined-schedules"},"Predefined schedules"),(0,a.kt)("p",null,"You may use one of several pre-defined schedules in place of a cron expression."),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"Entry                  | Description                                | Equivalent To\n-----                  | -----------                                | -------------\n@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *\n@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *\n@weekly                | Run once a week, midnight on Sunday        | 0 0 0 * * 0\n@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *\n@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *\n@minutely              | Run once a minute, beginning of minute     | 0 * * * * *\n@manually              | Never runs                                 | N/A\n")),(0,a.kt)("h3",{id:"intervals"},"Intervals"),(0,a.kt)("p",null,"You may also schedule a job to execute at fixed intervals.  This is supported by\nformatting the cron spec like this:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"@every <duration>\n")),(0,a.kt)("p",null,'where "duration" is a string accepted by time.ParseDuration\n(',(0,a.kt)("a",{parentName:"p",href:"http://golang.org/pkg/time/#ParseDuration"},"http://golang.org/pkg/time/#ParseDuration"),")."),(0,a.kt)("p",null,'For example, "@every 1h30m10s" would indicate a schedule that activates every\n1 hour, 30 minutes, 10 seconds.'),(0,a.kt)("p",null,"Note: The interval does not take the job runtime into account.  For example,\nif a job takes 3 minutes to run, and it is scheduled to run every 5 minutes,\nit will have only 2 minutes of idle time between each run."),(0,a.kt)("h3",{id:"fixed-times"},"Fixed times"),(0,a.kt)("p",null,"You may also want to schedule a job to be executed once. This is supported by\nformatting the cron spec like this:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"@at <datetime>\n")),(0,a.kt)("p",null,'Where "datetime" is a string accepted by time.Parse in RFC3339 format\n(',(0,a.kt)("a",{parentName:"p",href:"https://golang.org/pkg/time/#Parse"},"https://golang.org/pkg/time/#Parse"),")."),(0,a.kt)("p",null,'For example, "@at 2018-01-02T15:04:00Z" would run the job on the specified date and time\nassuming UTC timezone.'),(0,a.kt)("h3",{id:"time-zones"},"Time zones"),(0,a.kt)("p",null,"Dkron is able to schedule jobs in time zones, if you specify the ",(0,a.kt)("inlineCode",{parentName:"p"},"timezone")," parameter in a\njob definition."),(0,a.kt)("p",null,"If the time zone is not specified, the following rules apply:"),(0,a.kt)("p",null,"All interpretation and scheduling is done in the machine's local time zone (as\nprovided by the Go time package (",(0,a.kt)("a",{parentName:"p",href:"http://www.golang.org/pkg/time"},"http://www.golang.org/pkg/time"),")."),(0,a.kt)("p",null,"Be aware that jobs scheduled during daylight-savings leap-ahead transitions will\nnot be run!"),(0,a.kt)("p",null,"If you specify ",(0,a.kt)("inlineCode",{parentName:"p"},"timezone")," the job will be scheduled taking into account daylight-savings\nand leap-ahead transitions, running the job in the actual time in the specified time zone."))}p.isMDXComponent=!0}}]);
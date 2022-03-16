import React from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';


export default function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();

  return (
    <header>
      <div className="container relative overflow-hidden">
        <div className="absolute inset-y-0 w-full h-full" aria-hidden="true">
          <div className="relative h-full">
            <svg className="absolute mt-8 -top-16" width={400} height={620} viewBox="0 0 1200 630" xmlns="http://www.w3.org/2000/svg"><defs><pattern id="pattern" x="0" y="0" width="100" height="100" patternUnits="userSpaceOnUse" patternTransform="translate(0, 0) rotate(0) skewX(0)"><svg width="46" height="46" viewBox="0 0 100 100"><g className="text-pink-100" fill="currentColor" opacity="1"><path d="M95.2171 74.9997V25.0003L50.1086 0L5 25.0003V74.9997L50.1086 100L95.2171 74.9997Z"></path></g></svg></pattern></defs><rect x="0" y="0" width={400} height={1200} fill="url(#pattern)"></rect></svg>
            <svg className="absolute mt-2 transform rotate-180 right-8" width={400} height={320} viewBox="0 0 1200 630" xmlns="http://www.w3.org/2000/svg"><defs><pattern id="pattern" x="0" y="0" width="100" height="100" patternUnits="userSpaceOnUse" patternTransform="translate(0, 0) rotate(0) skewX(0)"><svg width="46" height="46" viewBox="0 0 100 100"><g className="text-pink-100" fill="currentColor" opacity="1"><path d="M95.2171 74.9997V25.0003L50.1086 0L5 25.0003V74.9997L50.1086 100L95.2171 74.9997Z"></path></g></svg></pattern></defs><rect x="0" y="0" width={400} height={1200} fill="url(#pattern)"></rect></svg>

          </div>
        </div>

        <div className="relative px-4 mx-auto mt-16 max-w-7xl sm:mt-24 sm:px-6">
          <div className="text-center">
            <h1 className="text-4xl font-extrabold tracking-tight text-fuchsia-900 sm:text-5xl md:text-6xl lg:text-5xl xl:text-6xl">
              <span className="block xl:inline">{siteConfig.tagline}</span>
            </h1>
            <p className="max-w-md mx-auto mt-3 text-lg text-gray-500 sm:text-xl md:mt-5 md:max-w-3xl">
              {siteConfig.customFields.description}
            </p>
            <div className="mt-10 sm:flex sm:justify-center">
              <div className="rounded-md shadow">
                <a href="https://github.com/distribworks/dkron/releases" className="flex items-center justify-center w-full px-8 py-3 text-base font-medium text-white border-2 border-transparent border-solid rounded-md bg-gradient-to-r from-fuchsia-600 to-pink-600 hover:bg-fuchsia-700 md:py-4 md:text-lg md:px-10 hover:no-underline hover:text-white/90">
                  Download
                </a>
              </div>
              <div className="mt-3 rounded-md shadow sm:mt-0 sm:ml-3">
                <a href="http://test.dkron.io:8080/ui" className="flex items-center justify-center w-full px-8 py-3 text-base font-medium bg-white border-2 border-solid rounded-md text-fuchsia-600 hover:text-fuchsia-700 border-fuchsia-600 hover:bg-white md:py-4 md:text-lg md:px-10 hover:no-underline">
                  Live demo
                </a>
              </div>
            </div>
          </div>
        </div>

        <div className="relative mt-4">
          <div className="absolute inset-0 flex flex-col" aria-hidden="true">
            <div className="flex-1" />
            <div className="flex-1 w-full" />
          </div>
          <div className="px-4 mx-auto max-w-7xl sm:px-6">
            <img
              className="relative rounded-lg"
              src="../img/job-list-new.png"
              alt="App screenshot"
            />
          </div>
        </div>
      </div>
    </header>
  );
}

import React from 'react';
import { SpeakerphoneIcon } from '@heroicons/react/outline';


export default function HomepageBanner() {
  return (
    <div className="bg-gradient-to-tr from-pink-600 to-fuchsia-700">
      <div className="px-3 py-3 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <div className="flex flex-wrap items-center justify-between">
          <div className="flex items-center flex-1">
            <span className="flex p-2 rounded-md bg-fuchsia-800">
              <SpeakerphoneIcon className="w-6 h-6 text-white" aria-hidden="true" />
            </span>
            <p className="mb-0 ml-3 font-medium text-white truncate">
              <span className="md:hidden">Dkron 3.0 is here!</span>
              <span className="hidden md:inline">Big news! We're excited to announce that Dkron 3.2 is here!</span>
            </p>
          </div>
          <div className="flex-shrink-0 order-3 w-full mt-2 sm:order-2 sm:mt-0 sm:w-auto">
            <a
              href="/blog/dkron-3-2/"
              className="flex items-center justify-center px-4 py-2 text-sm font-medium bg-white border border-transparent rounded-sm shadow-sm text-fuchsia-600 hover:bg-indigo-50"
            >
              Learn more
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

import React from 'react';
import CountUp from 'react-countup';

export default function HomepageStats() {
    return (
      <section className="py-12 mt-16 bg-zinc-100">
        <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
          <h2 className="text-4xl font-bold text-fuchsia-700">Stats</h2>
          <div className="flex flex-row">
            <div class="basis-1/2 ">
              <CountUp end={807} duration={5} useEasing={true} className="text-6xl font-bold" />
              <h3>Clusters Running</h3>
            </div>
            <div class="basis-1/2">
              <CountUp end={2536} duration={5} useEasing={true} className="text-6xl font-bold" />
              <h3>Servers Running</h3>
            </div>
          </div>
        </div>
      </section>
    );
  }

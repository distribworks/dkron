import React from 'react';

const features = [
  {
    title: 'Easy integration',
    img: '../img/integration.png',
    description: (
      <>
        Dkron is easy to setup and use. Choose your OS package and it's ready to run out-of-the-box. The administration panel and it's simple JSON API makes a breeze to integrate with you current workflow or deploy system.
      </>
    ),
  },
  {
    title: 'Always available',
    img: '../img/available.png',
    description: (
      <>
        Using the power of the Raft protocol, Dkron is designed to be always available. If the cluster leader node fails, a follower will replace it, all without human intervention.
      </>
    ),
  },
  {
    title: 'Flexible targets',
    img: '../img/targets.png',
    description: (
      <>
        Simple but powerful tag-based target node selection for jobs. Tag node count allows to run jobs in an arbitrary number of nodes in the same group or groups.
      </>
    ),
  },
];


export default function HomepageFeatures() {
  return (
    <section className="py-12 mt-16 bg-zinc-100">
      <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <h2 className="text-4xl font-bold text-fuchsia-700">Features</h2>
        <dl className="space-y-10 lg:space-y-0 lg:grid lg:grid-cols-3 lg:gap-16">
          {features.map((feature) => (
            <div key={feature.title}>
              <dt>
                <div className="">
                  <img src={feature.img} alt="" />
                </div>
                <p className="mt-4 text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-pink-400 to-fuchsia-600">{feature.title}</p>
              </dt>
              <dd className="mt-4 ml-0 text-base text-gray-500">{feature.description}</dd>
            </div>
          ))}
        </dl>
      </div>
    </section>
  );
}

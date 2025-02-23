import React from 'react';
import { Fragment } from 'react'
import { CheckIcon, MinusIcon } from '@heroicons/react/solid'

const tiers = [
  {
    name: 'Free',
    href: 'https://github.com/distribworks/dkron/',
    priceYearly: 'FREE',
  },
  {
    name: 'Pro',
    href: '/pro/',
    priceYearly: 450,
  }
]
const sections = [
  {
    name: 'Features',
    features: [
      { name: "Executor plugins", tiers: { Free: true, Pro: true } },
      { name: "Processor plugins", tiers: { Free: true, Pro: true } },
      { name: "Web UI", tiers: { Free: true, Pro: true } },
      { name: "Rest API", tiers: { Free: true, Pro: true } },
      { name: "Job retries", tiers: { Free: true, Pro: true } },
      { name: "Job chaining", tiers: { Free: true, Pro: true } },
      { name: "Concurrency control", tiers: { Free: true, Pro: true } },
      { name: "Metrics", tiers: { Free: true, Pro: true } },
      { name: "Embedded storage engine", tiers: { Free: true, Pro: true } },
      { name: "Docker executor", tiers: { Pro: true } },
      { name: "AWS ECS executor", tiers: { Pro: true } },
      { name: "Elasticsearch processor", tiers: { Pro: true } },
      { name: "Advanced Email processor", tiers: { Pro: true } },
      { name: "Slack processor", tiers: { Pro: true } },
      { name: "Encryption", tiers: { Pro: true } },
      { name: "Web UI Authentication", tiers: { Pro: true } },
      { name: "API Authentication", tiers: { Pro: true } },
      { name: "Access Control", tiers: { Pro: true } },
      { name: "Cross region failover", tiers: { Pro: true } },
      { name: "Dedicated Support", tiers: { Free: "None", Pro: "Email" }},
      { name: "License", tiers: { Free: "LGPL", Pro: "Commercial - No custom terms" }},
      { name: "Purchasing", tiers: { Pro: "Credit Card" }},
    ],
  }
]


export default function HomepagePricing() {
  return (
    <section className="py-12 bg-white">
      <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <h2 className="text-4xl font-bold text-fuchsia-700">Plans</h2>

        {/* Comparison table */}
        <div className="max-w-2xl py-4 mx-auto bg-white sm:py-8 lg:max-w-7xl">
          {/* xs to lg */}
          <div className="space-y-24 lg:hidden">
            {tiers.map((tier) => (
              <section key={tier.name}>
                <div className="px-4 mb-8">
                  <h2 className="text-lg font-medium leading-6 text-gray-900">{tier.name}</h2>
                  <p className="mt-4 mb-0">
                    {tier.name === "Pro" ? (
                      <>
                        <span className="text-4xl font-extrabold text-gray-900">${tier.priceYearly}</span>
                        <span className="text-base font-medium text-gray-500">/year</span>
                      </>
                    ) : (
                      <span className="text-4xl font-extrabold text-gray-900">{tier.priceYearly}</span>
                    )}
                  </p>
                </div>

                {sections.map((section) => (
                  <table key={section.name} className="w-full inline-table">
                    <caption className="px-4 py-3 text-sm font-medium text-left text-gray-900 border-t border-gray-200 bg-gray-50">
                      {section.name}
                    </caption>
                    <thead>
                      <tr>
                        <th className="sr-only" scope="col">
                          Feature
                        </th>
                        <th className="sr-only" scope="col">
                          Included
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {section.features.map((feature) => (
                        <tr key={feature.name} className="border-t border-gray-200">
                          <th className="px-4 py-5 text-sm font-normal text-left text-gray-500" scope="row">
                            {feature.name}
                          </th>
                          <td className="py-5 pr-4">
                            {typeof feature.tiers[tier.name] === 'string' ? (
                              <span className="block text-sm text-right text-gray-700">{feature.tiers[tier.name]}</span>
                            ) : (
                              <>
                                {feature.tiers[tier.name] === true ? (
                                  <CheckIcon className="w-5 h-5 ml-auto text-green-500" aria-hidden="true" />
                                ) : (
                                  <MinusIcon className="w-5 h-5 ml-auto text-gray-400" aria-hidden="true" />
                                )}

                                <span className="sr-only">{feature.tiers[tier.name] === true ? 'Yes' : 'No'}</span>
                              </>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                ))}

                <div className="px-4 pt-5 border-t border-gray-200">
                  {tier.name === "Pro" ? (
                    <a
                      href="https://gum.co/dkron-pro"
                      target="_blank"
                      className="block w-full py-3 text-base font-medium text-center text-white border-2 border-transparent rounded-md gumroad-button bg-gradient-to-r from-fuchsia-600 to-pink-600 hover:bg-fuchsia-700 md:py-4 md:text-lg md:px-10 hover:no-underline hover:text-white/90"
                    >
                      Buy
                    </a>
                  ) : (
                    <a
                      href="https://github.com/distribworks/dkron/releases"
                      className="block w-full py-3 text-base font-medium text-center text-white border-2 border-transparent rounded-md bg-gradient-to-r from-fuchsia-600 to-pink-600 hover:bg-fuchsia-700 md:py-4 md:text-lg md:px-10 hover:no-underline hover:text-white/90"
                    >
                      Download
                    </a>
                  )}
                </div>
              </section>
            ))}
          </div>

          {/* lg+ */}
          <div className="hidden lg:block">
            <table className="w-full table-fixed">
              <caption className="sr-only">Pricing plan comparison</caption>
              <thead>
                <tr>
                  <th className="pb-4 pl-6 pr-6 text-sm font-medium text-left text-gray-900" scope="col">
                    <span className="sr-only">Feature by</span>
                    <span>Plans</span>
                  </th>
                  {tiers.map((tier) => (
                    <th
                      key={tier.name}
                      className="w-1/4 px-6 pb-4 text-lg font-medium leading-6 text-left text-gray-900"
                      scope="col"
                    >
                      {tier.name}
                      <svg className="w-2 h-2 mx-2 text-gray-800/20" fill="currentColor" viewBox="0 0 8 8">
                        <circle cx={4} cy={4} r={3} />
                      </svg>
                      <a className="text-xs" href={tier.href}>Learn more</a>
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="border-t border-gray-200 divide-y divide-gray-200">
                <tr>
                  <th className="py-8 pl-6 pr-6 text-sm font-medium text-left text-gray-900 align-top" scope="row">
                    Pricing
                  </th>
                  {tiers.map((tier) => (
                    <td key={tier.name} className="h-full px-6 py-8 align-top">
                      <div className="flex flex-col justify-between h-full">
                        <div>
                          <p className="mb-0">
                            {tier.name === "Pro" ? (
                              <>
                                <span className="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-pink-400 to-fuchsia-600">${tier.priceYearly}</span>
                                <span className="text-base font-medium text-gray-500">/year</span>
                              </>
                            ) : (
                              <span className="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-pink-400 to-fuchsia-600">{tier.priceYearly}</span>
                            )}
                          </p>
                        </div>
                      </div>
                    </td>
                  ))}
                </tr>
                {sections.map((section) => (
                  <Fragment key={section.name}>
                    <tr>
                      <th
                        className="py-3 pl-6 text-sm font-medium text-left text-gray-900 bg-gray-50"
                        colSpan={4}
                        scope="colgroup"
                      >
                        {section.name}
                      </th>
                    </tr>
                    {section.features.map((feature) => (
                      <tr key={feature.name}>
                        <th className="py-2 pl-6 pr-6 text-sm font-normal text-left text-gray-500" scope="row">
                          {feature.name}
                        </th>
                        {tiers.map((tier) => (
                          <td key={tier.name} className="px-6 py-2">
                            {typeof feature.tiers[tier.name] === 'string' ? (
                              <span className="block text-sm text-gray-700">{feature.tiers[tier.name]}</span>
                            ) : (
                              <>
                                {feature.tiers[tier.name] === true ? (
                                  <CheckIcon className="w-5 h-5 text-green-500" aria-hidden="true" />
                                ) : (
                                  <MinusIcon className="w-5 h-5 text-gray-400" aria-hidden="true" />
                                )}

                                <span className="sr-only">
                                  {feature.tiers[tier.name] === true ? 'Included' : 'Not included'} in {tier.name}
                                </span>
                              </>
                            )}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </Fragment>
                ))}
              </tbody>
              <tfoot>
                <tr className="border-t border-gray-200">
                  <th className="sr-only" scope="row">
                    Choose your plan
                  </th>
                  {tiers.map((tier) => (
                    <td key={tier.name} className="px-6 pt-5">
                      {tier.name === "Pro" ? (
                        <>
                          <a
                            href="https://buy.stripe.com/9AQ9BP1UV8jkgowaEF"
                            target="_blank"
                            className="block w-full py-3 mt-2 text-base font-medium text-center text-white border-2 border-transparent rounded-md gumroad-button bg-gradient-to-r from-fuchsia-600 to-pink-600 hover:bg-fuchsia-700 md:py-4 md:text-lg md:px-10 hover:no-underline hover:text-white/90"
                          >
                            Subscribe
                          </a>
                          <a
                            href={tier.href}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="block w-full mt-1 text-center"
                          >
                            Learn more
                          </a>
                        </>
                      ) : (
                        <a
                          href={tier.href}
                          className="block w-full py-3 mb-5 text-base font-medium text-center text-white border-2 border-transparent rounded-md bg-gradient-to-r from-fuchsia-600 to-pink-600 hover:bg-fuchsia-700 md:py-4 md:text-lg md:px-10 hover:no-underline hover:text-white/90"
                        >
                          Download
                        </a>
                      )}
                    </td>
                  ))}
                </tr>
              </tfoot>
            </table>
          </div>
        </div>

        <p className="mb-0">All sales come with a two week, 100% money back guarantee.</p>
        <p>Looking to embed Dkron in your virtual server or appliance? <a href="docs/pro/commercial-faq">Read the Commercial FAQ.</a></p>
      </div>
    </section>
  );
}



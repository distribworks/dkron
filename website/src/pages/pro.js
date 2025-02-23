
import React from 'react';
import Layout from '@theme/Layout';
import { ShieldCheckIcon, HeartIcon, SupportIcon, CheckIcon } from '@heroicons/react/outline';

const features = [
  {
    title: 'Security',
    img: <ShieldCheckIcon className="w-12 h-12" />,
    description: (
      <>
        Pro has enhanced security using industry standard SSL encryption for communication between all components of the application, the embedded storage engine and nodes.
        You can also enable basic authentication to restrict access to the WebUI and the API.
      </>
    ),
  },
  {
    title: 'Pro plugins',
    img: <HeartIcon className="w-12 h-12" />,
    description: (
      <>
        Do you need to store job output in Elasticsearch? Do you need to run docker based jobs?
        Dkron Pro adds some commercially supported plugins ready to cover your needs.
      </>
    ),
  },
  {
    title: 'Support',
    img: <SupportIcon className="w-12 h-12" />,
    description: (
      <>
        Priority support from the author.
        Workload automation is a critical process in your business.
        Guarantee direct access to a Dkron expert. Your subscription gives you priority support for any unforeseen issues.
      </>
    ),
  },
];

const details = [
  {
    name: 'Documentation',
    description:
      <>
        <p>Detailed documentation about configuring and using each feature can be found in the Dkron docs site. Read the <a href="/pro/commercial-faq/">Commercial FAQ</a> for further details.</p>
        <p>Sales of Dkron Pro also benefit the community by ensuring that Dkron itself will remain well supported for the foreseeable future.</p>
      </>
  },
  {
    name: 'Support',
    description:
      <>
        <p>When you buy Dkron Pro, a custom URL associated with your email address will be sent to you. You use this URL to install the package corresponding to your architecture. You configure and use Dkron Pro exactly like you would Dkron.</p>
        <p><strong>Pro tip</strong>: use a mailing list for your email when purchasing to ensure you get critical email updates, even if employees leave the company.</p>
      </>
  },
  {
    name: 'Installation',
    description:
      <>
        <p>Dkron Pro will receive bug fixes and new functionality over time. All upgrades will be free to subscribers with a simple package upgrade. See the changelog for more detail.</p>
      </>
  },
  {
    name: 'Upgrades',
    description:
      <>
        <p>Dkron Pro will receive bug fixes and new functionality over time. All upgrades will be free to subscribers with a simple package upgrade. See the changelog for more detail.</p>
      </>
  },
  {
    name: 'Licensing',
    description:
      <>
        <p>Dkron is available under the terms of the GNU LGPLv3 license.</p>
        <p>In addition to its useful functionality, buying Dkron Pro grants your organization a Dkron commercial license instead of the GNU LGPL, avoiding any legal issues your lawyers might raise. Please see the <a href="/pro/commercial-faq/">Commercial FAQ</a> for further detail on licensing including options for distributing Dkron Pro with your own products.</p>
      </>
  },
];

const featuresList = [
  { name: 'Multi-region support'},
  { name: 'Full SSL encryption'},
  { name: 'Elasticsearch processor'},
  { name: 'Docker executor'},
  { name: 'AWS ECS executor'},
  { name: 'Advanced email processor'},
  { name: 'WebUI and API authorization'},
]


function Pro() {
  return (
    <Layout title="Dkron Pro">
      <div className="main-wrapper">
        <header className="bg-gradient-to-tr from-pink-600 to-fuchsia-700">
          <div className="container relative overflow-hidden">
            <div className="absolute inset-y-0 w-full h-full" aria-hidden="true">
              <div className="relative h-full">
                <svg className="absolute mt-8 -top-16 opacity-10" width={400} height={620} viewBox="0 0 1200 630" xmlns="http://www.w3.org/2000/svg"><defs><pattern id="pattern" x="0" y="0" width="100" height="100" patternUnits="userSpaceOnUse" patternTransform="translate(0, 0) rotate(0) skewX(0)"><svg width="46" height="46" viewBox="0 0 100 100"><g className="text-pink-100" fill="currentColor" opacity="1"><path d="M95.2171 74.9997V25.0003L50.1086 0L5 25.0003V74.9997L50.1086 100L95.2171 74.9997Z"></path></g></svg></pattern></defs><rect x="0" y="0" width={400} height={1200} fill="url(#pattern)"></rect></svg>
                <svg className="absolute mt-2 transform rotate-180 right-8 opacity-10" width={400} height={320} viewBox="0 0 1200 630" xmlns="http://www.w3.org/2000/svg"><defs><pattern id="pattern" x="0" y="0" width="100" height="100" patternUnits="userSpaceOnUse" patternTransform="translate(0, 0) rotate(0) skewX(0)"><svg width="46" height="46" viewBox="0 0 100 100"><g className="text-pink-100" fill="currentColor" opacity="1"><path d="M95.2171 74.9997V25.0003L50.1086 0L5 25.0003V74.9997L50.1086 100L95.2171 74.9997Z"></path></g></svg></pattern></defs><rect x="0" y="0" width={400} height={1200} fill="url(#pattern)"></rect></svg>
              </div>
            </div>

            <div className="absolute right-0 flex items-center justify-center transform -translate-x-1/2 -translate-y-1/2 opacity-10 left-1/2 top-1/2">
              <svg class="w-96 h-auto" viewBox="0 0 120 124" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M111.75 52.2198V74.4199C111.75 82.3199 107.54 89.6198 100.7 93.5698L81.4795 104.67C80.1895 97.7599 75.6896 91.9798 69.5796 88.9498C68.0796 88.1998 67.1197 86.6698 67.1197 84.9898V71.9099C67.1197 67.8399 63.8196 64.5399 59.7496 64.5399C55.6796 64.5399 52.3794 67.8399 52.3794 71.9099V84.9898C52.3794 86.6698 51.4195 88.1998 49.9095 88.9498C43.7995 91.9898 39.2996 97.7599 38.0096 104.67L18.7896 93.5698C11.9496 89.6198 7.73955 82.3199 7.73955 74.4199V52.2198C14.3695 54.5598 21.6196 53.5499 27.3096 49.7699C28.7096 48.8399 30.5095 48.7799 31.9695 49.6199L43.1495 56.0698C46.6095 58.0698 51.1095 57.0798 53.2295 53.6898C55.4595 50.1298 54.2695 45.4699 50.6695 43.3899L39.3396 36.8498C37.8796 36.0098 37.0395 34.4198 37.1395 32.7398C37.5595 25.9298 34.8095 19.1398 29.4695 14.5698L48.6895 3.46982C55.5295 -0.480176 63.9596 -0.480176 70.7996 3.46982L90.0196 14.5698C84.6796 19.1398 81.9297 25.9298 82.3497 32.7298C82.4497 34.4098 81.6095 36.0098 80.1495 36.8498L68.8196 43.3899C65.2896 45.4299 64.0897 49.9298 66.1197 53.4598C68.1597 56.9898 72.6595 58.1899 76.1895 56.1599L87.5196 49.6199C88.9796 48.7799 90.7795 48.8399 92.1795 49.7699C97.8695 53.5399 105.12 54.5598 111.75 52.2198ZM68.8895 97.1898C68.1695 96.6198 67.1097 97.1498 67.1097 98.0698V108.75C67.1097 112.82 63.8095 116.12 59.7395 116.12C55.6695 116.12 52.3697 112.82 52.3697 108.75V98.0798C52.3697 97.1498 51.2996 96.6398 50.5696 97.2098C46.9596 100.08 44.7296 104.6 45.0196 109.64C45.4496 116.97 51.3395 122.93 58.6595 123.45C67.2895 124.06 74.4695 117.24 74.4695 108.75C74.4795 104.06 72.2895 99.8898 68.8895 97.1898ZM20.4595 45.0698C21.3295 44.7298 21.4297 43.5398 20.6197 43.0798L11.3794 37.7398C7.78944 35.6698 6.59961 31.0298 8.79961 27.4798C10.9196 24.0498 15.5196 23.1099 19.0196 25.1299L29.8096 31.3598C29.7996 25.9998 26.8695 20.7899 21.6595 18.1699C15.0295 14.8399 6.75953 17.0899 2.71953 23.3099C-1.91047 30.4499 0.429507 39.9199 7.69951 44.1199C11.7395 46.4599 16.4295 46.6598 20.4595 45.0698ZM117.179 23.9898C113.109 16.9398 104.1 14.5298 97.0396 18.5998C92.9896 20.9398 90.4796 24.8998 89.8296 29.1898C89.6896 30.1098 90.6795 30.7799 91.4795 30.3099L100.72 24.9698C104.25 22.9298 108.75 24.1399 110.79 27.6699C112.83 31.1999 111.62 35.6998 108.09 37.7398L98.8497 43.0798C98.0497 43.5398 98.1396 44.7298 98.9996 45.0698C103.04 46.6598 107.73 46.4599 111.78 44.1299C118.84 40.0499 121.249 31.0398 117.179 23.9898Z" fill="white">
                </path>
              </svg>
            </div>

            <div className="relative px-4 mx-auto my-32 max-w-7xl sm:mt-24 sm:px-6">
              <div className="text-center">
                <h1 className="text-4xl font-extrabold tracking-tight text-white sm:text-5xl md:text-6xl lg:text-5xl xl:text-6xl">
                  <span className="block xl:inline">Dkron Pro</span>
                </h1>
                <p className="max-w-md mx-auto mt-3 text-lg text-gray-100 sm:text-xl md:mt-5 md:max-w-3xl">
                  Improved security, features and reliability for your scheduled jobs
                </p>
                <div className="mt-10 sm:flex sm:justify-center">
                  <div className="rounded-md shadow">
                  <a 
                    href="https://buy.stripe.com/9AQ9BP1UV8jkgowaEF"
                    target="_blank"
                    className="block w-full py-3 mt-2 text-base font-medium text-center text-white border-2 border-transparent rounded-md gumroad-button bg-gradient-to-r from-fuchsia-600 to-pink-600 hover:bg-fuchsia-700 md:py-4 md:text-lg md:px-10 hover:no-underline hover:text-white/90"
                  >
                    Subscribe
                  </a>
                  </div>
                </div>
                <p className="max-w-md mx-auto mt-3 text-lg text-gray-100 sm:text-xl md:mt-5 md:max-w-3xl">
                    Licenses are not transferable to another company. We will transfer the license from a user-specific email to a group email address (e.g. john_smith@example.com -> tech@example.com) but only for the same domain. It is strongly recommended that you buy the license using a group email address so the license is not attached to any one employee’s email address.
                </p>
              </div>
            </div>
          </div>
        </header>


        <section className="py-16 bg-zinc-100">
          <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
            <h2 className="text-4xl font-bold text-fuchsia-700">Key Features</h2>
            <dl className="mt-16 space-y-10 lg:space-y-0 lg:grid lg:grid-cols-3 lg:gap-16">
              {features.map((feature) => (
                <div key={feature.title}>
                  <dt>
                    <div className="">
                      {feature.img}
                    </div>
                    <p className="mt-4 text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-pink-400 to-fuchsia-600">{feature.title}</p>
                  </dt>
                  <dd className="mt-4 ml-0 text-base text-gray-500">{feature.description}</dd>
                </div>
              ))}
            </dl>
          </div>
        </section>

        <section class="bg-gradient-to-tr from-fuchsia-600 to-pink-600">
          <div class="py-8 bg-zinc-900/90">
            <div className="px-4 py-16 mx-auto max-w-7xl sm:px-6 lg:py-24 lg:px-8 lg:grid lg:grid-cols-3 lg:gap-x-8">
              <div>
                <h2 className="text-base font-semibold tracking-wide text-pink-200 uppercase">Everything you need</h2>
                <p className="mt-2 text-3xl font-extrabold text-pink-600">Features</p>
              </div>
              <div className="mt-12 lg:mt-0 lg:col-span-2">
                <dl className="space-y-10 sm:space-y-0 sm:grid sm:grid-cols-2 sm:grid-rows-4 sm:grid-flow-col sm:gap-x-6 sm:gap-y-6 lg:gap-x-8">
                  {featuresList.map((item) => (
                    <div key={item.name} className="relative">
                      <dt>
                        <CheckIcon className="absolute w-6 h-6 text-green-500" aria-hidden="true" />
                        <p className="mb-0 text-lg font-medium leading-6 text-gray-50 ml-9">{item.name}</p>
                      </dt>
                    </div>
                  ))}
                </dl>
              </div>
            </div>
          </div>
        </section>

        <section className="bg-white">
          <div className="grid max-w-2xl grid-cols-1 px-4 py-16 mx-auto gap-y-16 gap-x-8 sm:px-6 lg:max-w-7xl lg:px-8 lg:grid-cols-2">
            <div>
              <h2 className="text-3xl font-extrabold tracking-tight sm:text-4xl text-fuchsia-700">Product Details</h2>

              <dl className="grid grid-cols-1 mt-16 gap-x-6 gap-y-10 sm:grid-cols-2 sm:gap-y-16 lg:gap-x-8">
                {details.map((detail) => (
                  <div key={detail.name} className="pt-4 border-t border-gray-200">
                    <dt className="font-medium text-transparent bg-clip-text bg-gradient-to-r from-pink-400 to-fuchsia-600">{detail.name}</dt>
                    <dd className="mt-3 ml-0 text-gray-500">{detail.description}</dd>
                  </div>
                ))}
              </dl>
            </div>
            <div className="grid items-center justify-center gap-4 sm:gap-6">
              <img
                src="../img/dkron-black.png"
                alt=""
                className="rounded-lg"
              />
              <img
                src="../img/dkron-gradient.png"
                alt=""
                className="rounded-lg"
              />
              <img
                src="../img/dkron-gray.png"
                alt=""
                className="rounded-lg"
              />
            </div>
          </div>
        </section>
      </div>
{/*

      <div class="container">
        <div class="row">
          <div class="main-content">
            <h2>Product details</h2>
            <h3>FEATURES</h3>
            <p>Dkron Pro contains the following functionality:</p>
            <ul>
              <li>
                Multi-region support
              </li>
              <li>
                Full SSL encryption
              </li>
              <li>
                Elasticsearch processor
              </li>
              <li>
                Docker executor
              </li>
              <li>
                AWS ECS executor
              </li>
              <li>
                Advanced email processor
              </li>
              <li>
                WebUI and API authorization
              </li>
            </ul>

          </div>
        </div>
      </div> */}
    </Layout>
  );
}

export default Pro;

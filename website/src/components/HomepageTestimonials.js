import React from 'react';


export default function HomepageTestimonials() {
  return (
    <section className="bg-gradient-to-tr from-fuchsia-600 to-pink-600">
      <div className="py-16 bg-zinc-900/90">
        <div className="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
          <h2 className="text-4xl font-bold text-pink-600">What people are saying</h2>
        </div>
        <div className="mx-auto max-w-7xl md:grid md:grid-cols-2 md:px-6 lg:px-8">
          <div className="px-4 py-12 sm:px-6 md:flex md:flex-col md:py-16 md:pl-0 md:pr-10 lg:pr-16">
            <blockquote className="mt-6 border-0 md:flex-grow md:flex md:flex-col">
              <div className="relative text-lg font-medium text-white md:flex-grow">
                <svg
                  className="absolute w-16 h-16 transform -translate-x-3 -translate-y-2 text-fuchsia-300/20 -top-8 -left-8"
                  fill="currentColor"
                  viewBox="0 0 32 32"
                  aria-hidden="true"
                >
                  <path d="M9.352 4C4.456 7.456 1 13.12 1 19.36c0 5.088 3.072 8.064 6.624 8.064 3.36 0 5.856-2.688 5.856-5.856 0-3.168-2.208-5.472-5.088-5.472-.576 0-1.344.096-1.536.192.48-3.264 3.552-7.104 6.624-9.024L9.352 4zm16.512 0c-4.8 3.456-8.256 9.12-8.256 15.36 0 5.088 3.072 8.064 6.624 8.064 3.264 0 5.856-2.688 5.856-5.856 0-3.168-2.304-5.472-5.184-5.472-.576 0-1.248.096-1.44.192.48-3.264 3.456-7.104 6.528-9.024L25.864 4z" />
                </svg>
                <p className="relative">
                  Dkron is an essential part of our fault-tolerant distributed automation workflow. Its deceptively simple setup hides a well implemented architecture. Using established tools to coordinate our large number of endpoints, allows us to set-and-forget many of the tasks that would normally require a single, dedicated server.
                </p>
              </div>
              <footer className="mt-8">
                <div className="flex items-start">
                  <div className="inline-flex flex-shrink-0 border-2 border-white rounded-full">
                    <img
                      className="w-12 h-12 rounded-full"
                      src="https://www.gravatar.com/avatar/196032196bfacfc2e1e4c5891ca500a6"
                      alt=""
                    />
                  </div>
                  <div className="ml-4">
                    <div className="text-base font-medium text-white">Geoff Jukes</div>
                    <div className="text-base font-medium text-pink-200">Director of Information Technology, Blackstone Publishing</div>
                  </div>
                </div>
              </footer>
            </blockquote>
          </div>
          <div className="px-4 py-12 sm:px-6 md:flex md:flex-col md:py-16 md:pl-0 md:pr-10 lg:pr-16">
            <blockquote className="mt-6 border-0 md:flex-grow md:flex md:flex-col">
              <div className="relative text-lg font-medium text-white md:flex-grow">
                <svg
                  className="absolute w-16 h-16 transform -translate-x-3 -translate-y-2 text-fuchsia-300/20 -top-8 -left-8"
                  fill="currentColor"
                  viewBox="0 0 32 32"
                  aria-hidden="true"
                >
                  <path d="M9.352 4C4.456 7.456 1 13.12 1 19.36c0 5.088 3.072 8.064 6.624 8.064 3.36 0 5.856-2.688 5.856-5.856 0-3.168-2.208-5.472-5.088-5.472-.576 0-1.344.096-1.536.192.48-3.264 3.552-7.104 6.624-9.024L9.352 4zm16.512 0c-4.8 3.456-8.256 9.12-8.256 15.36 0 5.088 3.072 8.064 6.624 8.064 3.264 0 5.856-2.688 5.856-5.856 0-3.168-2.304-5.472-5.184-5.472-.576 0-1.248.096-1.44.192.48-3.264 3.456-7.104 6.528-9.024L25.864 4z" />
                </svg>
                <p className="relative">
                  Dkron is an essential piece of our infrastructure: you can trust that the cron jobs will execute. Also, Dkron help us to identify failed cron jobs really fast. The system and the API are very easy to use and well designed. Will definitely keep using it!
                </p>
              </div>
              <footer className="mt-8">
                <div className="flex items-start">
                  <div className="inline-flex flex-shrink-0 border-2 border-white rounded-full">
                    <img
                      className="w-12 h-12 rounded-full"
                      src="https://www.gravatar.com/avatar/37803c500666eded79f8cd23b276b77c"
                      alt=""
                    />
                  </div>
                  <div className="ml-4">
                    <div className="text-base font-medium text-white">Adrià Galín</div>
                    <div className="text-base font-medium text-pink-200">Director of SRE, Jobandtalent</div>
                  </div>
                </div>
              </footer>
            </blockquote>
          </div>
        </div>
      </div>
    </section>
  );
}


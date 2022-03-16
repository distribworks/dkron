import React from 'react';

const users = [
  {
    name: "MallGroup",
    url: "https://mallgroup.com",
    img: "/img/mallgroup-logo.svg"
  },
  {
    name: "Rubicon Project",
    url: "https://rubiconproject.com/",
    img: "/img/rubicon-logo.svg"
  },
  {
    name: "Jobandtalent",
    url: "https://www.jobandtalent.com",
    img: "/img/jt-logo.png"
  },
  {
    name: "Linkfluence",
    url: "https://linkfluence.com/",
    img: "/img/logo_linkfluence.png"
  },
  {
    name: "Voiceworks",
    url: "https://www.voiceworks.com/en",
    img: "/img/voiceworks-logo.svg"
  },
  {
    name: "benemen",
    url: "https://rubiconproject.com/",
    img: "/img/benemen-logo.png"
  },
  {
    name: "Blackstone Publishing",
    url: "https://rubiconproject.com/",
    img: "/img/Blackstone_publishing_Logo_v2.png"
  },
  {
    name: "Allianz",
    url: "https://rubiconproject.com/",
    img: "/img/Allianz.svg"
  },
]

export default function HomepageUsers() {
  return (
    <section className="bg-slate-700/90">
      <div className="px-4 py-12 mx-auto border-t border-gray-200 max-w-7xl sm:px-6 lg:py-20 lg:px-8">
        <h3 className="mb-12 text-base font-normal tracking-wider text-center uppercase text-gray-50">Trusted by</h3>
        <div className="grid grid-cols-2 gap-12 md:grid-cols-3 lg:grid-cols-4">
          {users.map((user) => (
            <div key={user.name} className="flex justify-center col-span-1 lg:col-span-1">
              <a href={user.url} target="_blank">
                <img className="h-10 transition duration-200 filter grayscale hover:grayscale-0 brightness-200 hover:brightness-100" src={user.img} alt={user.name} />
              </a>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

import React from 'react';

const users = [
  {
    name: "MallGroup",
    url: "https://mallgroup.com",
    img: "/img/mallgroup-logo.svg"
  },
  {
    name: "American Express",
    url: "https://aexp.com",
    img: "/img/AXP_BlueBoxLogo_EXTRALARGEscale_RGB_DIGITAL_1600x1600.png"
  },
  {
    name: "Flickr",
    url: "https://flickr.com",
    img: "/img/Flickr_logo.png"
  },
  {
    name: "Socal Gas",
    url: "https://www.socalgas.com/",
    img: "/img/SoCal_Gas.png"
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
    name: "enreach",
    url: "https://www.enreach.fi/",
    img: "/img/logo-enreach.svg"
  },
  {
    name: "Blackstone Publishing",
    url: "https://www.blackstonepublishing.com/",
    img: "/img/Blackstone_publishing_Logo_v2.png"
  },
  {
    name: "Allianz",
    url: "https://www.allianz.com/",
    img: "/img/Allianz.svg"
  },
  {
    name: "Kata AI",
    url: "https://kata.ai/",
    img: "/img/kata_ai.png"
  },
  {
    name: "delcampe",
    url: "https://www.delcampe.net/",
    img: "/img/logo-delcampe.svg"
  },
]

export default function HomepageUsers() {
  return (
    <section className="bg-slate-700/90">
      <div className="px-4 py-12 mx-auto border-t border-gray-200 max-w-7xl sm:px-6 lg:py-20 lg:px-8">
        <h3 className="mb-12 text-sm font-bold tracking-wider text-center uppercase text-gray-50">Trusted by</h3>
        <div className="grid grid-cols-2 gap-12 md:grid-cols-3 lg:grid-cols-4">
          {users.map((user) => (
            <div key={user.name} className="flex justify-center col-span-1 lg:col-span-1">
              <a href={user.url} target="_blank" rel="noopener noreferrer">
                <img className="h-10 transition duration-200 filter grayscale hover:grayscale-0 brightness-200 hover:brightness-100" src={user.img} alt={user.name} />
              </a>
            </div>
          ))}
        </div>
      </div>

      <div className="max-w-lg px-4 py-12 mx-auto border-t border-gray-200 sm:px-6 lg:py-20 lg:px-8">
        <h3 className="mb-12 text-sm font-bold tracking-wider text-center uppercase text-gray-50">Partners</h3>
        <div className="grid grid-cols-2 gap-12">
          <div className="flex justify-center col-span-1 lg:col-span-1">
            <a href="https://www.noraina.cloud/" target="_blank" rel="noopener noreferrer">
            <img src="//www.noraina.cloud/assets/static/logo.0f04e74.1968b830a90e93ddc4373b0d6426be8b.png" className="h-16 transition duration-200 filter brightness-100 hover:brightness-200" alt="Noraina Cloud" />
            </a>
          </div>
          <div className="flex justify-center col-span-1 lg:col-span-1">
            <a href="https://gemfury.com/f/partner" target="_blank" rel="noopener noreferrer">
              <img src="//badge.fury.io/fp/gemfury.svg" alt="" />
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}

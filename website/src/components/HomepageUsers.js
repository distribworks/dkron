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
            <a href="https://metal.equinix.com" target="_blank" rel="noopener noreferrer">
              <svg width="251" height="25" viewBox="0 0 190 20" xmlns="http://www.w3.org/2000/svg">
                <path d="m45.871 10.46h3.7322v-0.98682h-3.7322v-2.171h5.5643v-0.98684h-6.7178v7.5658h6.9214v-0.9869h-5.7679v-2.4342z"/>
                <path d="m75.457 10.855c0.0679 1.1842-0.8142 2.171-2.0357 2.2368h-0.2714c-1.1536 0.0658-2.1036-0.8553-2.1714-1.9737v-0.2631-4.4737h-1.0857v4.5395c-0.0679 1.6447 1.2214 3.0921 2.9857 3.1578h0.2714c2.5786 0 3.3929-1.7105 3.3929-3.1578v-4.5395h-1.0858v4.4737z"/>
                <path d="m83.939 6.3158h-1.0857v7.5h1.0857v-7.5z"/>
                <path d="m97.239 12.237-5.4964-5.921h-1.0857v7.5658h1.0857v-5.9211l5.4964 5.9211h1.0857v-7.5658h-1.0857v5.921z"/>
                <path d="m106.2 6.3158h-1.085v7.5h1.085v-7.5z"/>
                <path d="m117.05 9.8026 3.257-3.4868h-1.425l-2.511 2.7632-2.511-2.7632h-1.493l3.122 3.4868-3.868 4.079h1.425l3.189-3.3553 3.054 3.3553h1.561l-3.8-4.079z"/>
                <path d="m60.121 6.0526c-2.2393 0-4.0714 1.7763-4.0714 3.9474v0.1316c-0.0679 2.171 1.6964 3.9473 3.8679 4.0131h0.2714c0.6107 0 1.2214-0.1315 1.8321-0.3947l0.6107 0.7237h1.2893l-1.0857-1.3158c0.95-0.7895 1.5607-1.9079 1.4929-3.1579 0.0678-2.171-1.6965-3.9474-3.9357-4.0132-0.0679 0.06579-0.1358 0.06579-0.2715 0.06579zm2.9857 4.079c0.0679 0.921-0.3392 1.7763-1.0178 2.3684l-0.6107-0.6579h-1.2893l1.0857 1.1842c-0.3393 0.1316-0.7464 0.1974-1.0857 0.1974-1.6286 0-2.9857-1.3158-2.9857-2.8948v-0.1315c-0.1357-1.579 1.0857-2.9606 2.7143-3.0921h0.1357c1.6964 0 2.9857 1.3158 2.9857 2.9605 0.0678-0.0658 0.0678 0 0.0678 0.0658z"/>
                <path d="m29.518 4.671v8.8816l-2.1036 0.7237v-10.329l-6.3107-2.171v14.671l-2.1036 0.7237v-16.118l-3.1893-1.0526-3.1893 1.0526v16.118l-2.1035-0.7237v-14.671l-6.3108 2.171v10.329l-2.1036-0.7237v-8.8816l-2.1036 0.6579v9.7368l6.3107 2.1053v-11.776l2.1036-0.72369v13.224l6.3107 2.1053v-17.5l1.0857-0.32895 1.0857 0.32895v17.5l6.3107-2.1053v-13.224l2.1036 0.72369v11.776l6.3107-2.1053v-9.7368l-2.1035-0.6579z" fill="#ED1C24"/>
                <path d="m35.625 14.342c1.0179 0 1.9-0.8553 1.9-1.8421s-0.8821-1.8421-1.9-1.8421-1.9 0.8553-1.9 1.8421c0 1.0526 0.8821 1.8421 1.9 1.8421zm0-0.1974c-0.95 0-1.6964-0.7236-1.6964-1.6447 0-0.921 0.7464-1.6447 1.6964-1.6447s1.6964 0.7237 1.6964 1.6447c0.0679 0.9211-0.6785 1.6447-1.6964 1.6447 0.0679 0 0.0679 0 0 0zm-0.7464-0.7236h0.475v-0.5921h0.2714l0.4071 0.5921h0.6108l-0.475-0.6579c0.2714-0.0658 0.475-0.329 0.4071-0.5921 0-0.4606-0.3393-0.6579-0.8143-0.6579h-0.95l0.0679 1.9079zm0.475-0.9869v-0.5263h0.4071c0.2036 0 0.3393 0.0658 0.3393 0.2632 0 0.1973-0.1357 0.2631-0.3393 0.2631h-0.4071z" fill="#ED1C24"/>
                <path d="m134.36 0h-0.339v20h0.339v-20z" fill="#707073"/>
                <path d="m147.25 6.3816h1.493l2.511 3.5526 2.442-3.5526h1.493v7.5h-1.425v-5.329l-2.578 3.5527h-0.068l-2.511-3.4869v5.329h-1.357v-7.5658z" fill="#707073"/>
                <path d="m158.11 6.3816h6.039v1.1842h-4.546v1.9737h4.071v1.1842h-4.071v1.9737h4.614v1.1842h-6.107v-7.5z" fill="#707073"/>
                <path d="m168.96 7.6316h-2.578v-1.25h6.582v1.25h-2.579v6.25h-1.425v-6.25z" fill="#707073"/>
                <path d="m177.04 6.3158h1.357l3.529 7.5h-1.493l-0.814-1.7763h-3.868l-0.814 1.7763h-1.425l3.528-7.5zm2.036 4.6053-1.425-2.9606-1.357 2.9606h2.782z" fill="#707073"/>
                <path d="m184.3 6.3816h1.425v6.25h4.275v1.1842h-5.7v-7.4342z" fill="#707073"/>
              </svg>
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

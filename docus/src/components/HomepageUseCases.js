import React from 'react';

const useCases = [
  {
    name: "Email delivery",
    img: "img/email-delivery.png"
  },
  {
    name: "Payroll generation",
    img:"img/payroll-generation.png"
  },
  {
    name: "Bookkeeping",
    img:"img/book-keeping.png"
  },
  {
    name: "Data consolidation for BI",
    img:"img/data-consolidation.png"
  },
  {
    name: "Recurring invoicing",
    img:"img/recurring-invoicing.png"
  },
  {
    name: "Data transfer",
    img:"img/data-transfer.png"
  },
]


export default function HomepageUseCases() {
  return (
    <section className="bg-gradient-to-tr from-pink-600 to-fuchsia-700">
      <div className="px-4 py-12 mx-auto border-t border-gray-200 max-w-7xl sm:px-6 lg:py-20 lg:px-8">
        <h3 className="text-4xl font-bold text-white">Example use cases</h3>
        <div className="grid sm:grid-cols-2 gap-0.5 md:grid-cols-3 lg:mt-8">
          {useCases.map((use) => (
            <div key={use.name} className="flex items-center col-span-1 px-8 py-8 bg-gray-50/20">
              <img className="shrink-0 max-h-16" src={use.img} alt={use.name} />
              <p className="mb-0 ml-4 text-xl text-white">{use.name}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

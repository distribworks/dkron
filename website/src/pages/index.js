import React from 'react';
import Layout from '@theme/Layout';
import HomepageBanner from '../components/HomepageBanner';
import HomepageHeader from '../components/HomepageHeader';
import HomepageDescPanel from '../components/HomepageDescPanel';
import HomepageFeatures from '../components/HomepageFeatures';
import HomepagePricing from '../components/HomepagePricing';
import HomepageTestimonials from '../components/HomepageTestimonials';
import HomepageUsers from '../components/HomepageUsers';
import HomepageUseCases from '../components/HomepageUseCases';
import HomepageStats from '../components/HomepageStats';

export default function Home() {
  return (
    <>
      <HomepageBanner />

      <Layout
        title={`Dkron - Cloud native job scheduling system`}
        description="Cloud native job scheduling system"
      >
        <HomepageHeader />

        <main>
          <HomepageDescPanel />

          <HomepageFeatures />

          <HomepageUseCases />

          <HomepagePricing />

          <HomepageStats />

          <HomepageTestimonials />

          <HomepageUsers />
        </main>
      </Layout>
    </>
  );
}

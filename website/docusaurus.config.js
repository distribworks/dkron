// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Dkron',
  tagline: 'Easy, Reliable Cron jobs',
  url: 'https://dkron.io',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'distribworks', // Usually your GitHub org/user name.
  projectName: 'dkron', // Usually your repo name.
  customFields: {
    description: 'A distributed Cron service with, API, no SPOF and an easy to use dashboard.',
    description_extended: 'Dkron is a system service for workload automation that runs scheduled jobs, just like unix cron service but distributed in several machines in a cluster. This is the only job scheduler in the market with truly no SPOF. It is open source and available for free.'
  },

  plugins: [
    async function myPlugin(context, options) {
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
          // Appends TailwindCSS and AutoPrefixer.
          postcssOptions.plugins.push(require("tailwindcss"));
          postcssOptions.plugins.push(require("autoprefixer"));
          return postcssOptions;
        },
      };
    },
  ],

  markdown: {
    mermaid: true,
  },
  themes: ['@docusaurus/theme-mermaid'],

  presets: [
    [
      'redocusaurus',
      {
        // Plugin Options for loading OpenAPI files
        specs: [
          {
            spec: 'static/openapi/openapi.yaml',
            route: '/api/',
          },
        ],
        // Theme Options for modifying how redoc renders them
        theme: {
          // Change with your site colors
          primaryColor: '#1890ff',
        },
      },
    ],
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/distribworks/dkron/tree/main/website/docs/',
          lastVersion: 'current',
          versions: {
            current: {
              label: 'v3',
              path: '',
            },
          },
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/distribworks/dkron/tree/main/website/blog/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      colorMode: {
        defaultMode: 'light',
        disableSwitch: true,
      },
      navbar: {
        title: '',
        logo: {
          alt: 'Dkron Logo',
          src: 'img/dkron-logo-black.png',
        },
        items: [
          {
            type: 'doc',
            docId: 'basics/getting-started',
            position: 'left',
            label: 'Docs',
          },
          {
            to: '/api/', label: 'API', position: 'left'
          },
          {to: '/blog', label: 'Blog', position: 'left'},
          {to: '/license', label: 'License', position: 'left'},
          {
            href: 'https://github.com/distribworks/dkron',
            label: 'GitHub',
            position: 'right',
          },
          {
            type: 'docsVersionDropdown',
          },
          {
            to: '/pro/', label: 'PRO', position: 'left', className: 'navbar-link-go-pro'
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Intro',
                to: '/docs/intro',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/dkron',
              },
              {
                label: 'Gitter',
                href: 'https://gitter.im/distribworks/dkron',
              },
              {
                label: 'Twitter',
                href: 'https://twitter.com/distribworks',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Blog',
                to: '/blog',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/distribworks/dkron',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Distributed Works. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;

module.exports = {
  title: 'StackHead',
  tagline: 'Open Source Web Server management',
  url: 'https://docs.stackhead.io',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'getstackhead', // Usually your GitHub org/user name.
  projectName: 'deployment', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'StackHead',
      logo: {
        alt: 'StackHead Logo',
        src: 'img/logo.svg',
      },
      links: [
        {
          to: 'introduction/getting-started',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {
          href: 'https://github.com/getstackhead/deployment',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [],
      copyright: `Copyright © ${new Date().getFullYear()} Mario Lubenka and StackHead contributors. Built with Docusaurus.`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          path: '../',
          routeBasePath: '',
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/getstackhead/deployment/edit/master/Documentation/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};

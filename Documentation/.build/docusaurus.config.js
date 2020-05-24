module.exports = {
  title: 'StackHead',
  tagline: 'Open Source Web Server management',
  url: 'https://docs.stackhead.io',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'getstackhead', // Usually your GitHub org/user name.
  projectName: 'stackhead', // Usually your repo name.
  themeConfig: {
    announcementBar: {
      id: 'wip_message', // Any value that will identify this message.
      content:
        'This project is highly <strong>WORK-IN-PROGRESS</strong>. Stuff may break and change without notice!',
      backgroundColor: 'darkorange'
    },
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
          href: 'https://github.com/getstackhead/stackhead',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [],
      copyright: `Copyright © ${new Date().getFullYear()} Mario Lubenka and StackHead contributors. Documentation built with Docusaurus.`,
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
          editUrl: 'https://github.com/getstackhead/stackhead/edit/master/Documentation/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};

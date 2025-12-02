// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  tutorialSidebar: [
    'quickstart',
    'installation',
    {
      type: 'category',
      label: 'Usage',
      items: ['cli', 'api'],
    },
    'configuration',
    'architecture',
    'monitoring',
    {
      type: 'category',
      label: 'Reference',
      items: ['troubleshooting', 'examples'],
    },
    'contributing',
    'todo',
  ],
};

module.exports = sidebars;

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
    'troubleshooting',
    {
      type: 'category',
      label: 'Developers',
      items: ['contributing', 'todo'],
    },
  ],
};

module.exports = sidebars;

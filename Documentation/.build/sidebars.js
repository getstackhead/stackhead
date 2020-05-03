module.exports = {
  someSidebar: {
    Introduction: [
      'introduction/getting-started',
    ],
    Concepts: [
      'concepts/projects',
      'concepts/provision-and-deploy',
      'concepts/capabilities'
    ],
    Configuration: [
      {
        type: 'category',
        label: 'Project definition',
        items: [
          'configuration/project-definition/introduction',
          'configuration/project-definition/container',
          'configuration/project-definition/native'
        ],
      },
      'configuration/capabilities',
      'configuration/plays'
    ],
  },
};

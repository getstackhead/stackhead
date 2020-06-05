module.exports = {
  someSidebar: {
    Introduction: [
      'introduction/concepts',
      'introduction/getting-started',
    ],
    // Concepts: [
    //   'concepts/projects',
    //   'concepts/provision-and-deploy',
    //   'concepts/capabilities'
    // ],
     Configuration: [
       {
         type: 'category',
         label: 'Project definition',
         items: [
           'configuration/project-definition/introduction',
           {
             type: 'category',
             label: 'Project types',
             items: [
               'configuration/project-definition/container',
               'configuration/project-definition/native'
             ],
           }
         ],
       },
    //   'configuration/capabilities',
    //   'configuration/security',
    //   'configuration/plays'
     ],
    'Technical Documentation': [
      'technical-documentation/workflow'
    ]
  },
};

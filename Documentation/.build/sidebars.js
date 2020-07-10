module.exports = {
  someSidebar: {
    Introduction: [
      'introduction/workflow',
      'introduction/installation',
      'introduction/getting-started',
      'introduction/plays'
    ],
    Configuration: [
      'configuration/project-definition',
      {
        type: 'category',
        label: 'Project types',
        items: [
          'configuration/project-definition/container',
          'configuration/project-definition/native'
        ]
      },
      'configuration/capabilities',
      'configuration/security'
    ],
    'Technical Documentation': [
      'technical-documentation/terraform',
      'technical-documentation/ssl-certificates'
    ]
  }
}

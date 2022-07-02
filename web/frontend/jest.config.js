module.exports = {
  preset: '@vue/cli-plugin-unit-jest/presets/no-babel',
  reporters: [
    'default',
    [
      'jest-junit',
      {
        outputName: 'unit-tests-frontend.xml',
        outputDirectory: '../../output'
      }
    ]
  ]
}

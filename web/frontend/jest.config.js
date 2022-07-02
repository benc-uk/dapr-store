module.exports = {
  preset: '@vue/cli-plugin-unit-jest/presets/no-babel'
}

// module.exports = {
//   preset: '@vue/cli-plugin-unit-jest',
//   moduleNameMapper: {
//     '\\.(css)$': 'identity-obj-proxy'
//   },
//   setupFiles: ['<rootDir>/tests/unit/setup.js'],
//   verbose: true,
//   collectCoverage: true,
//   coverageDirectory: '../../output/web-coverage',
//   coverageReporters: ['html', 'text-summary'],
//   reporters: [
//     'default',
//     [
//       'jest-junit',
//       {
//         outputName: 'unit-tests-frontend.xml',
//         outputDirectory: '../../output'
//       }
//     ]
//   ]
// }

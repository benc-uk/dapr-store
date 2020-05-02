module.exports = {
  preset: '@vue/cli-plugin-unit-jest',
  moduleNameMapper: {
    '\\.(css)$': 'identity-obj-proxy'
  },
  setupFiles: ['<rootDir>/tests/no-console.js'],
  verbose: true
}

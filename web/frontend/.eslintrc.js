module.exports = {
  root: true,
  env: {
    node: true,
    es6: true,
    jest: true
  },
  extends: ['plugin:vue/vue3-essential', 'eslint:recommended'],
  parserOptions: {
    ecmaVersion: 2020
  },
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'vue/multi-word-component-names': 'off'
  }
}

/*

{
  "env": {
    "node": true,
    "es6": true,
    "jest": true
  },
  "extends": ["eslint:recommended", "plugin:vue/recommended"],
  "rules": {
    // Errors & best practices
    "no-var": "error",
    "no-console": "off",
    "no-unused-vars": ["error", { "argsIgnorePattern": "next|res|req" }],
    "curly": "error",

    // ES6
    "arrow-spacing": "error",
    "arrow-parens": "error",

    // Vue
    "vue/html-closing-bracket-newline": "off",
    "vue/max-attributes-per-line": "off",
    "vue/singleline-html-element-content-newline": "off",
    "vue/html-self-closing": "off",
    "vue/multi-word-component-names": "off"
  }
}
*/

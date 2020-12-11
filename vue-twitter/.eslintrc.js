module.exports = {
  parser: 'vue-eslint-parser',

  extends: [
    'plugin:vue/vue3-recommended',
  ],
  rules: {
    'vue/custom-event-name-casing': 'off',
    'no-use-before-define': 'off',
    // 'no-setting-before-define': [
    //   'error',
    //   {
    //     functions: false,
    //     classes: true,
    //   },
    // ],

    'no-unused-vars': [
      'error',
      {
        argsIgnorePattern: '^h$',
        varsIgnorePattern: '^h$',
      },
    ],
    'space-before-function-paren': 'off',
    'selector-pseudo-element-no-unknown': ['off', {
      'ignorePseudoElements': ['v-deep']
    }]
  },
};

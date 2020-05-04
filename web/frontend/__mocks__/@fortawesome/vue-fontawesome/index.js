// Stubs out FontAwesome icons and prevents SVGs leaking into tests
import Vue from 'vue'

export default Vue.component('fa', {
  render(createElement) {
    return createElement('i')
  }
})
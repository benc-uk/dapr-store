import { createLocalVue, shallowMount } from '@vue/test-utils'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import ErrorBox from '@/components/ErrorBox.vue'

describe('ErrorBox.vue', () => {
  it('renders error message', () => {
    const errorMsg = 'This is an error'
    const wrapper = shallowMount(ErrorBox, { localVue,
      propsData: { error: errorMsg }
    })
    expect(wrapper).toMatchSnapshot()
  })
})
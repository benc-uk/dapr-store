import { mount } from '@vue/test-utils'
import router from '@/router'
import ErrorBox from '@/components/ErrorBox.vue'

describe('ErrorBox.vue', () => {
  it('renders error message', () => {
    const errorMsg = 'This is an error'
    const wrapper = mount(ErrorBox, {
      propsData: { error: errorMsg },
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.html()).toMatch('This is an error')
  })
})

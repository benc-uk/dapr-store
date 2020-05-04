import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import Cart from '@/views/Cart.vue'

jest.mock('@/mixins/api')

import { userProfile } from '@/main.js'

describe('Cart.vue', () => {
  it('renders cart contents', async () => {
    userProfile.userName = 'demo@example.net'
    const wrapper = mount(Cart, {
      localVue
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})
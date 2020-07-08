import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import Cart from '@/views/Cart.vue'

jest.mock('@/mixins/api')
jest.mock('@/mixins/auth')

describe('Cart.vue', () => {
  it('renders cart contents', async () => {
    const wrapper = mount(Cart, {
      localVue
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})
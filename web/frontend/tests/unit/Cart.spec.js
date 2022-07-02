import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import router from '@/router'
import Cart from '@/views/Cart.vue'

jest.mock('@/services/api')
jest.mock('@/services/auth')

describe('Cart.vue', () => {
  it('renders cart contents', async () => {
    const wrapper = mount(Cart, {
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Top&nbsp;Hat')
    expect(wrapper.html()).toMatch('Cravat')
  })
})

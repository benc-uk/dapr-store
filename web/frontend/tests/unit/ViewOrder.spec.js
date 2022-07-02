import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import ProductSearch from '@/views/ViewOrder.vue'
import router from '@/router'

jest.mock('@/services/api')
jest.mock('@/services/auth')

describe('ViewOrder.vue', () => {
  it('shows order details', async () => {
    router.push('/order/ord-mock')
    await router.isReady()

    const wrapper = mount(ProductSearch, {
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('ord-mock')
    expect(wrapper.html()).toMatch('Ascot')
  })
})

import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import ProductSingle from '@/views/ProductSingle.vue'
import router from '@/router'

const productId = 'prd1'
jest.mock('@/services/api')

describe('ProductSingle.vue', () => {
  it('renders product details', async () => {
    router.push('/product/' + productId)
    await router.isReady()

    const wrapper = mount(ProductSingle, {
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Top&nbsp;Hat')
  })
})

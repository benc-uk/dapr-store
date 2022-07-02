import { mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import ProductSearch from '@/views/ProductSearch.vue'
import router from '@/router'

jest.mock('@/services/api')

describe('ProductSearch.vue', () => {
  it('renders search for Ascot', async () => {
    router.push('/search/Ascot')
    await router.isReady()

    const wrapper = mount(ProductSearch, {
      global: {
        plugins: [router]
      }
    })

    await flushPromises()
    expect(wrapper.html()).toMatch('Ascot')
  })
})

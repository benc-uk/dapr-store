import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import VueRouter from 'vue-router'
localVue.use(VueRouter)
const router = new VueRouter()

import ProductSingle from '@/views/ProductSingle.vue'
import { userProfile } from '@/main.js'

jest.mock('@/mixins/api')

const productId = '2'

describe('ProductSingle.vue', () => {
  it('renders product details', async () => {
    router.push({ name: 'single-product', params: { id: productId } })

    const wrapper = mount(ProductSingle, {
      localVue,
      router
    })

    await flushPromises()
    wrapper.vm.addToCart()
    await localVue.nextTick()
    await flushPromises()
    expect(userProfile.cart).toHaveLength(1)
    expect(userProfile.cart[0].id).toEqual(productId)

    expect(wrapper).toMatchSnapshot()
  })
})
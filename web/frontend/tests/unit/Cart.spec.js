import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import Cart from '@/views/Cart.vue'

jest.mock('@/mixins/api')

// import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
// jest.mock('FontAwesomeIcon')

import { userProfile } from '@/main.js'

describe('Cart.vue', () => {
  it('renders cart contents', async () => {
    // Fake cart in userProfile
    userProfile.cart.push({
      id: '1',
      name: 'Top Hat (6″)',
      cost: '39.95',
      description: 'Made from 100% Wool and nice',
      image: '/img/catalog/1.jpg',
      onOffer: false
    })

    const wrapper = mount(Cart, {
      localVue
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})
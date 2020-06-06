import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import Account from '@/views/Account.vue'

jest.mock('@/mixins/api')
jest.mock('@/mixins/auth')

describe('Account.vue', () => {
  it('renders user profile', async () => {
    const wrapper = mount(Account, {
      localVue
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})
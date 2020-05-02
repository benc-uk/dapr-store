import { createLocalVue, mount } from '@vue/test-utils'
import flushPromises from 'flush-promises'
import BootstrapVue from 'bootstrap-vue'
const localVue = createLocalVue()
localVue.use(BootstrapVue)

import Account from '@/views/Account.vue'
import { userProfile } from '@/main.js'

jest.mock('@/mixins/api')

describe('Account.vue', () => {
  it('renders user profile', async () => {
    // Fake user
    userProfile.userName = 'demo@example.net'

    const wrapper = mount(Account, {
      localVue
    })

    await flushPromises()
    expect(wrapper).toMatchSnapshot()
  })
})
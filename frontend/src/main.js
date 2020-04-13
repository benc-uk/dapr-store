import Vue from 'vue'
import VueRouter from 'vue-router'
import * as msal from 'msal'

import App from './App.vue'
import router from './router'
import User from './user'

// Use Vue Bootstrap and theme
import BootstrapVue from 'bootstrap-vue'
Vue.use(BootstrapVue)
import 'bootswatch/dist/materia/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { library as faIcons } from '@fortawesome/fontawesome-svg-core'
import { faUser, faUserPlus, faShoppingBasket, faTrophy, faIdCard, faShoppingCart, faSignOutAlt } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
faIcons.add(faUser, faUserPlus, faShoppingBasket, faTrophy, faIdCard, faShoppingCart, faSignOutAlt)
Vue.component('fa', FontAwesomeIcon)

Vue.config.productionTip = false
Vue.use(VueRouter)

// Global user & auth details
let userProfile = new User()
export { userProfile }

// export { user }
export const msalApp = new msal.UserAgentApplication({
  auth: {
    clientId: '69972365-c1b6-494d-9579-5b9de2790fc3',
    redirectUri: window.location.origin
  }
})
export const accessTokenRequest = {
  scopes: [ 'api://69972365-c1b6-494d-9579-5b9de2790fc3/store-api' ]
}

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app')

// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - main app initialization and startup
// ----------------------------------------------------------------------------

import Vue from 'vue'
//import VueRouter from 'vue-router'
import App from './App.vue'

// Global services
import auth from './services/auth'
import api from './services/api'

// Use Vue Bootstrap and theme
import BootstrapVue from 'bootstrap-vue'
Vue.use(BootstrapVue)
import 'bootswatch/dist/materia/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

// Set up FontAwesome
import { library as faIcons } from '@fortawesome/fontawesome-svg-core'
import { faUser, faUserPlus, faShoppingBasket, faTrophy, faIdCard, faShoppingCart, faSignOutAlt, faTrashAlt, faRedoAlt, faSearch, faPlusCircle, faMinusCircle } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
faIcons.add(faUser, faUserPlus, faShoppingBasket, faTrophy, faIdCard, faShoppingCart, faSignOutAlt, faTrashAlt, faRedoAlt, faSearch, faPlusCircle, faMinusCircle)
Vue.component('fa', FontAwesomeIcon)

// And client side routes (held in router.js)
import router from './router'

// Let's go!
appStartup()

//
// App start up synchronized using await with the config API call
//
async function appStartup() {
  // Take local defaults from .env.development or .env.development.local
  // Or fall back to internal defaults
  let API_ENDPOINT = process.env.VUE_APP_API_ENDPOINT || '/'
  let AUTH_CLIENT_ID = process.env.VUE_APP_AUTH_CLIENT_ID || ''

  // Load config at runtime from special `/config` endpoint on frontend-host
  try {
    let configResp = await fetch('/config')
    if (configResp.ok) {
      const config = await configResp.json()
      API_ENDPOINT = config.API_ENDPOINT
      AUTH_CLIENT_ID = config.AUTH_CLIENT_ID
      console.log('### Config loaded:', config)
    }
  } catch (err) {
    console.warn('### Failed to fetch \'/config\' endpoint. Defaults will be used')
  }

  auth.configure(AUTH_CLIENT_ID)
  api.configure(API_ENDPOINT, AUTH_CLIENT_ID, 'store-api')

  // Actually mount & start the Vue app, kinda important
  new Vue({
    router,
    render: (h) => h(App),
  }).$mount('#app')
}
// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr Store frontend - main app initialization and startup
// ----------------------------------------------------------------------------

import Vue from 'vue'
import VueRouter from 'vue-router'
import * as msal from 'msal'
import Axios from 'axios'
import App from './App.vue'
import { User } from './user'

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
Vue.use(VueRouter)

// Global user & auth details
let userProfile = new User()
export { userProfile }

// Global stuff populated by call to appStartup
export let config = {}
export let msalApp = {}
export let accessTokenRequest = {}

appStartup()

//
// Most of the app config & initialization moved here
// To be synchronized with the config API call
//
async function appStartup() {
  // Load config at runtime from special `/config` endpoint on frontend-host
  try {
    let resp = await Axios.get('/config')
    config = resp.data
  } catch (err) {
    console.warn('### Failed to fetch \'/config\' endpoint. Defaults will be used')
    config = {
      API_ENDPOINT: '/',
      AUTH_CLIENT_ID: ''
    }
  }

  console.log('### ==== Config ====')
  console.log(config)

  // MSAL config used for signing in users with MS identity platform
  if (config.AUTH_CLIENT_ID) {
    console.log(`### USER SIGN-IN ENABLED. Using clientId: ${config.AUTH_CLIENT_ID}`)
    msalApp = new msal.UserAgentApplication({
      auth: {
        clientId: config.AUTH_CLIENT_ID,
        redirectUri: window.location.origin
      }
    })
    accessTokenRequest = {
      scopes: [ `api://${config.AUTH_CLIENT_ID}/store-api` ]
    }
  } else {
    console.log('### USER SIGN-IN DISABLED. Will run in demo mode, with dummy users')
  }

  if (config.API_ENDPOINT != '/') {
    console.log(`### API_ENDPOINT. Overridden with '${config.API_ENDPOINT}'`)
  }

  new Vue({
    router,
    render: (h) => h(App),
  }).$mount('#app')
}
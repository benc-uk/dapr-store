import Vue from 'vue'
import VueRouter from 'vue-router'

import App from './App.vue'
import router from './router'

// Use Vue Bootstrap and theme
import BootstrapVue from 'bootstrap-vue'
Vue.use(BootstrapVue)
import 'bootswatch/dist/materia/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { library as faIcons } from '@fortawesome/fontawesome-svg-core'
import { faUserCircle, faShoppingBasket, faTrophy } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
faIcons.add(faUserCircle, faShoppingBasket, faTrophy)
Vue.component('fa', FontAwesomeIcon)

Vue.config.productionTip = false
Vue.use(VueRouter)

new Vue({
  router,
  render: (h) => h(App),
}).$mount('#app')

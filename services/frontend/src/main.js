import Vue from 'vue'
import App from './App.vue'
import store from './store'
import router from './router'

import moment from 'vue-moment'

import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'

// Install BootstrapVue
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)

Vue.use(moment);

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.config.productionTip = false

import axios from 'axios'

const keyword = "Bearer"
const token = localStorage.getItem("token")
const baseURL = process.env.VUE_APP_BASE_URI ? process.env.VUE_APP_BASE_URI : "http://localhost:8080/api/"
const requester = axios.create({ baseURL: baseURL })

if (token) {
  requester.defaults.headers.common['Authorization'] = `${keyword} ${token}`
}

Vue.prototype.$http = requester

new Vue({
  store,
  router,
  render: h => h(App)
}).$mount('#app')

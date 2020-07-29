import Vue from 'vue'
import VueRouter from 'vue-router'
// import Home from '../views/Home.vue'

import News from '../views/News.vue'
import Login from '@/views/Login.vue' 
import Register from '@/views/Register.vue'
import SingleNews from '../views/SingleNews.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: News
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
  },
  {
    path: '/news/:uid',
    name: 'SingleNews',
    component: SingleNews
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router

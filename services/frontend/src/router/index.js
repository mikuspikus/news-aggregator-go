import Vue from 'vue'
import VueRouter from 'vue-router'

import store from '../store'

import News from '@/views/News.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import SingleNews from '@/views/SingleNews.vue'
import AdminPanel from '@/views/AdminPanel.vue'


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
    path: '/admin',
    name: "AdminPanel",
    component: AdminPanel,
    meta: {
      IsAdminOnly: true,
    }
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

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.IsAdminOnly)) {
    if (!store.getters.isAdmin) {
      next({ name: 'Home' })
      return
    }
  }
  next()
})

export default router

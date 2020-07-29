import Vue from 'vue'
import Vuex from 'vuex'

import axios from 'axios'

const keyword = 'Bearer'

const loginuri = process.env.VUE_APP_LOGINURI ? process.env.VUE_APP_LOGINURI : 'http://localhost:8080/api/auth/token'
const registeruri = process.env.VUE_APP_REGISTERURI ? process.env.VUE_APP_REGISTERURI : 'http://localhost:8080/api/user'
const refreshuri = process.env.VUE_APP_REFRESHURI ? process.env.VUE_APP_REFRESHURI : 'http://localhost:8080/api/auth/refresh'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    status: '',
    isAdmin: localStorage.getItem('isAdmin') || '',
    token: localStorage.getItem('token') || '',
    refresh_token: localStorage.getItem('refresh_token') || '',
    username: localStorage.getItem('username') || '',
    uid: localStorage.getItem('uid') || '',
  },
  mutations: {
    request(state) {
      state.status = 'loading'
    },

    login(state, payload) {
      state.status = 'success'

      state.token = payload.token
      state.refresh_token = payload.refresh_token
      state.username = payload.username
      state.uid = payload.uid
      state.isAdmin = payload.isAdmin
    },

    refresh(state, payload) {
      state.status = 'success'

      state.token = payload.token
      state.refresh_token = payload.refresh_token
    },

    error(state) {
      state.status = 'error'
    },

    logout(state) {
      state.status = ''

      state.token = ''
      state.refresh_token = ''
      state.username = ''
      state.uid = ''
      state.isAdmin = ''
    }
  },
  actions: {
    relogin({ commit }, data) {
      return new Promise(
        (resolve, reject) => {
          commit('request')

          axios({ url: refreshuri, data: {token: this.state.token, refresh_token: this.state.refresh_token}, method: 'POST' })
            .then(
              response => {
                const payload = {
                  token: response.data.token,
                  refresh_token: response.data.refresh_token
                }
                commit('refresh', payload)

                localStorage.setItem('token', payload.token)
                localStorage.setItem('refresh_token', payload.refresh_token)

                const requester = data.axios
                delete requester.defaults.headers.common['Authorization']
                requester.defaults.headers.common['Authorization'] = `${keyword} ${payload.token}`

                resolve(response)
              }
            )
            .catch(error => {
              commit('error')

              localStorage.removeItem('token')
              localStorage.removeItem('refresh_token')
              localStorage.removeItem('username')
              localStorage.removeItem('uid')
              localStorage.removeItem('isAdmin')

              reject(error)
            })
        }
      )
    },

    login({ commit }, data) {
      return new Promise(
        (resolve, reject) => {
          commit('request')

          axios({ url: loginuri, data: data.credentials, method: 'POST' })
            .then(
              response => {
                const payload = {
                  token: response.data.token,
                  refresh_token: response.data.refresh_token,
                  username: response.data.username,
                  uid: response.data.uid,
                  isAdmin: response.data.is_admin
                }

                commit('login', payload)

                localStorage.setItem('token', payload.token)
                localStorage.setItem('refresh_token', payload.refresh_token)
                localStorage.setItem('username', payload.username)
                localStorage.setItem('uid', payload.uid)
                localStorage.setItem('isAdmin', payload.is_admin)

                const requester = data.axios
                requester.defaults.headers.common['Authorization'] = `${keyword} ${payload.token}`

                resolve(response)
              })
            .catch(error => {
              commit('error')

              localStorage.removeItem('token')
              localStorage.removeItem('refresh_token')
              localStorage.removeItem('username')
              localStorage.removeItem('uid')
              localStorage.removeItem('isAdmin')

              reject(error)
            })
        }
      )
    },

    register({ commit }, data) {
      return new Promise(
        (resolve, reject) => {
          commit('request')

          axios({ url: registeruri, data: data.credentials, method: 'POST' })
            .then(response => {
              console.log(response.data)
              const payload = {
                token: response.data.token,
                refresh_token: response.data.refresh_token,
                username: response.data.username,
                uid: response.data.uid
              }

              commit('login', payload)

              localStorage.setItem('token', payload.token)
              localStorage.setItem('refresh_token', payload.refresh_token)
              localStorage.setItem('username', payload.username)
              localStorage.setItem('uid', payload.uid)

              const requester = data.axios
              requester.defaults.headers.common['Authorization'] = `${keyword} ${payload.token}`

              resolve(response)
            })
            .catch(error => {
              commit('error')

              localStorage.removeItem('token')
              localStorage.removeItem('refresh_token')
              localStorage.removeItem('username')
              localStorage.removeItem('uid')

              reject(error)
            })
        }
      )
    },

    logout({ commit }, data) {
      return new Promise(
        (resolve) => {
          commit('logout')

          localStorage.removeItem('token')
          localStorage.removeItem('refresh_token')
          localStorage.removeItem('username')
          localStorage.removeItem('uid')
          localStorage.removeItem('isAdmin')

          const requester = data.axios
          delete requester.defaults.headers.common['Authorization']

          resolve()
        }
      )
    }
  },

  getters: {
    isAdmin: state => !!state.isAdmin,
    isLogged: state => !!state.token,
    status: state => state.status,
    uid: state => state.uid,
    username: state => state.username
  },

  modules: {
  }
})

import Vue from 'vue'
import Vuex from 'vuex'

import axios from 'axios'

const keyword = 'Bearer'

const loginuri = process.env.VUE_APP_LOGIN_URI ? process.env.VUE_APP_LOGIN_URI : 'http://localhost:8080/api/auth/token'
const registeruri = process.env.VUE_APP_REGISTER_URI ? process.env.VUE_APP_REGISTER_URI : 'http://localhost:8080/api/user'
const refreshuri = process.env.VUE_APP_REFRESH_URI ? process.env.VUE_APP_REFRESH_URI : 'http://localhost:8080/api/auth/refresh'

Vue.use(Vuex)

function _makePayload(responseData) {
  return {
    token: responseData.token,
    refresh_token: responseData.refresh_token,
    username: responseData.username,
    uid: responseData.uid,
    is_admin: responseData.is_admin
  }
}

export default new Vuex.Store({
  state: {
    status: '',
    is_admin: 'true' === localStorage.getItem('is_admin') || false,
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
      state.is_admin = payload.is_admin
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
      state.is_admin = false
    }
  },
  actions: {
    relogin({ commit }, data) {
      return new Promise(
        (resolve, reject) => {
          commit('request')

          axios({ url: refreshuri, data: { token: this.state.token, refresh_token: this.state.refresh_token }, method: 'POST' })
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
              localStorage.removeItem('is_admin')

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
                const payload = _makePayload(response.data)

                commit('login', payload)

                localStorage.setItem('token', payload.token)
                localStorage.setItem('refresh_token', payload.refresh_token)
                localStorage.setItem('username', payload.username)
                localStorage.setItem('uid', payload.uid)
                localStorage.setItem('is_admin', payload.is_admin)

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
              localStorage.removeItem('is_admin')

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
              const payload = _makePayload(response.data)

              commit('login', payload)

              localStorage.setItem('token', payload.token)
              localStorage.setItem('refresh_token', payload.refresh_token)
              localStorage.setItem('username', payload.username)
              localStorage.setItem('uid', payload.uid)
              localStorage.setItem("is_admin", payload.is_admin)

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
              localStorage.removeItem('is_admin')

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
          localStorage.removeItem('is_admin')

          const requester = data.axios
          delete requester.defaults.headers.common['Authorization']

          resolve()
        }
      )
    }
  },

  getters: {
    isAdmin: state => state.is_admin,
    isLogged: state => !!state.token,
    status: state => state.status,
    uid: state => state.uid,
    username: state => state.username
  },

  modules: {
  }
})

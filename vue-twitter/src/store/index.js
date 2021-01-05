import { createStore } from 'vuex'
import { Base64 } from 'js-base64'
import { getGHToken, setGHToken } from '@/utils/local-storage'

const state = {
  gh_token: getGHToken()
}

const mutations = {
  SET_GH_TOKEN: (state, token) => {
    state.gh_token = token
  }
}

const actions = {
  setGHToken ({ commit }, token) {
    const _token = Base64.encode(token)
    setGHToken(_token)
    commit('SET_GH_TOKEN', _token)
  }
}

const getters = {
  gh_token: state => state.gh_token ? Base64.decode(state.gh_token) : ''
}

const store = createStore({
  actions,
  state,
  mutations,
  getters
})

export default store

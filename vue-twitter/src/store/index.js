import { createStore } from 'vuex'
import { Base64 } from 'js-base64'
import { getGHToken, setGHToken } from '@/utils/local-storage'

const state = {
  gh_token: getGHToken(),
  is_mobile: false
}

const mutations = {
  SET_GH_TOKEN: (state, token) => {
    state.gh_token = token
  },
  SET_IS_MOBILE: (state, isMobile) => {
    state.is_mobile = isMobile
  }
}

const actions = {
  setGHToken ({ commit }, token) {
    const _token = Base64.encode(token)
    setGHToken(_token)
    commit('SET_GH_TOKEN', _token)
  },
  setIsMobile ({ commit }, isMobile) {
    commit('SET_IS_MOBILE', isMobile)
  }
}

const getters = {
  ghToken: state => state.gh_token ? Base64.decode(state.gh_token) : '',
  isMobile: state => state.is_mobile
}

const store = createStore({
  actions,
  state,
  mutations,
  getters
})

export default store

// vuex 参考使用地址： https: //github.com/vuejs/vuex/tree/4.0

import { createStore } from 'vuex'

export const store = createStore({
  state () {
    return {
      count: 1
    }
  }
})
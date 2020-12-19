// 参考链接： https://github.com/vuejs/vue-router-next

import {
  createRouter,
  createWebHashHistory
} from 'vue-router'

// todo 目前不支持 () => import( /* webpackChunkName: 'Home' */ ) 引入
import Home from '@/page/home/index.vue'
import Login from '@/page/login/index.vue'
const routerHistory = createWebHashHistory(process.env.BASE_URL)

const router = createRouter({
  history: routerHistory,
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/',
      name: 'Home',
      component: Home
    }
  ]
})

export default router

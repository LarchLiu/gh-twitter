// 参考链接： https://github.com/vuejs/vue-router-next

import {
  createRouter,
  createWebHashHistory
} from 'vue-router'

// todo 目前不支持 () => import( /* webpackChunkName: 'Home' */ ) 引入
const routerHistory = createWebHashHistory(process.env.BASE_URL)

const router = createRouter({
  history: routerHistory,
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/page/login/index')
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('@/page/home/index')
    }
  ]
})

export default router

import { createApp } from 'vue'
import Vue3VideoPlayer from '@cloudgeek/vue3-video-player'
import store from './store'
import App from './App.vue'
import router from './router'
import './style/common.less'
import '@cloudgeek/vue3-video-player/dist/vue3-video-player.css'

// axios
import axios from './http'

// UI框架
import {
  Avatar,
  Badge,
  Button,
  Checkbox,
  Drawer,
  Dropdown,
  Input,
  Form,
  Menu,
  Modal,
  message
} from 'ant-design-vue'

import './registerServiceWorker'
// import 'ant-design-vue/dist/antd.css'

const app = createApp(App)

app.config.globalProperties.$axios = axios
app.config.globalProperties.$message = message

// ui注入
app
  .use(Avatar)
  .use(Badge)
  .use(Button)
  .use(Checkbox)
  .use(Drawer)
  .use(Dropdown)
  .use(Menu)
  .use(Input)
  .use(Form)
  .use(Modal)

// 插件注入
app.use(router)
  .use(store)
  .use(Vue3VideoPlayer)
  .mount('#app')

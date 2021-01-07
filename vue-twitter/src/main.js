import { createApp } from 'vue'
import store from './store'
import App from './App.vue'
import router from './router'
import './style/common.less'

// axios
import axios from './http'

// UI框架
import {
  Avatar,
  Badge,
  Button,
  Checkbox,
  Input,
  Form,
  Modal,
  Pagination,
  message
} from 'ant-design-vue'
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
  .use(Input)
  .use(Form)
  .use(Modal)
  .use(Pagination)

// 插件注入
app.use(router)
  .use(store)
  .mount('#app')

import { createApp } from 'vue'
import { store } from './store';
import App from './App.vue'
import router from './router';
import  './style/common.less';


// axios
import axios from './http';


// UI框架
import {
  Avatar,
  Button,
  Input,
  Calendar,
  Badge,
  Form
} from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

const app = createApp(App)

app.config.globalProperties.$axios = axios;

//ui注入
app
.use(Avatar)
.use(Button)
.use(Input)
.use(Form)
.use(Calendar)
.use(Badge)


// 插件注入
app.use(router)
.use(store)
.mount('#app')

import axios from 'axios'
import { message } from 'ant-design-vue'
// import router from '@/router'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_GLOB_API_URL, // url = base url + request url
  withCredentials: true, // send cookies when cross-domain requests
  timeout: 60000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent

    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */
  // response => response.data,
  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    const res = response.data
    if (!(res.err === null || res.err === undefined)) {
      let msg = ((res.err.message) ? res.err.message : res.err.code)
      if (msg === 'ER_DUP_ENTRY') {
        msg = '数据库已有重复条目'
      } else if (msg === 'ER_PARSE_ERROR') {
        msg = '参数有误'
      } else if (msg === 'ER_CRM_SUPPLIER_TYPE_ERR') {
        msg = '供应商传输数据类型有误'
      } else if (msg === 'ER_CRM_RETAIL_TYPE_ERR') {
        msg = '零售店传输数据类型有误'
      } else if (msg === 'ER_BEIKE_IN_RANGE_NEED_CONFIRM') {
        msg = '已有此周期数据'
      } else if (msg === 'ER_BEIKE_IN_RANGE_PERIOD_ERROR') {
        msg = '周期选择有误'
      } else if (msg === 'ER_BEIKE_OUT_RANGE_PERIOD_ERROR') {
        msg = '周期选择超出范围'
      } else if (msg === 'ER_BEIKE_SAME_EXCEL') {
        msg = '该Excel已存在'
      } else if (msg === 'ER_BEIKE_PERIOD_SOME_DAY_LOST') {
        msg = '所选计算周期有日期遗漏'
      }
      message.error({
        content: msg || 'Error',
        duration: 5
      })
      // 50008:非法的token; 50012:其他客户端登录了;  50014:Token 过期了;
      // if (res.code === 50008 || res.code === 50012 || res.code === 50014) {
      //   // 请自行在引入 MessageBox
      //   // import { Message, MessageBox } from 'element-ui'
      //   MessageBox.confirm('你已被登出，可以取消继续留在该页面，或者重新登录', '确定登出', {
      //     confirmButtonText: '重新登录',
      //     cancelButtonText: '取消',
      //     type: 'warning'
      //   }).then(() => {
      //     store.dispatch('user/resetToken').then(() => {
      //       location.reload() // 为了重新实例化vue-router对象 避免bug
      //     })
      //   })
      // }
      return Promise.reject(res.err || 'Error')
    } else {
      return res
    }
  },
  error => {
    console.log('err' + error) // for debug
    if (error.response) {
      const data = error.response.data
      if (data.err) {
        error.message = ((data.err.message) ? data.err.message : data.err.code)
        console.log('error message ' + error.message)
      }
      if (error.response.status === 401) {
        if (error.message === 'ER_AUTH_USER_OR_PWD_ERR') {
          error.message = '用户名或密码错误！'
        } else if (error.message === 'jwt expired') {
          error.message = '认证过期，请重新登录！'
          // permission 会 catch error 跳转至 login
          // store.dispatch('user/resetToken').then(() => {
          //   // location.reload() // 为了重新实例化vue-router对象 避免bug
          //   router.push(`/login?redirect=${router.currentRoute.path}`)
          // })
        } else {
          error.message = '权限认证有误，请重新登录！'
          // permission 会 catch error 跳转至 login
          // store.dispatch('user/resetToken').then(() => {
          //   router.push(`/login?redirect=${router.currentRoute.path}`)
          // })
        }
      } else {
        if (error.message === 'ER_API_PATH_NOT_MATCH') {
          error.message = 'API路径有误，请确认！'
        } else if (error.message === 'ER_REGISTER_USER_NAME_EXIST') {
          error.message = '账号已存在，请重新输入！'
        } else if (error.message === 'ER_DOMAIN_NOT_PERMIT') {
          error.message = '该域名尚无权限，请联系管理员！'
        } else if (error.message === 'ER_QUALITY_CODE_ERR') {
          error.message = '追溯码有误，请确认追溯码及产品信息！'
        } else if (error.message === 'ER_QUALITY_VER_ERR') {
          error.message = '验证码有误，非认证产品，请谨慎！'
        } else if (error.message === 'ER_USER_EXIST') {
          error.message = '该用户为其他公司管理员，不可添加！'
        } else if (error.message === 'ER_CART_TRADE_PARAM_ERROR') {
          error.message = '结算参数有误，请返回购物车重新选购'
        } else if (error.message === 'ER_GOODS_NOT_FOUND') {
          error.message = '无此商品信息'
        } else if (error.message === 'ER_BEIKE_EXCEL_FORMAT_ERROR') {
          error.message = '上传Excel文件格式有误'
        } else if (error.message === 'ER_BEIKE_EXCEL_DATE_ERROR') {
          error.message = '上传Excel文件数据日期与所选日期不符'
        } else if (error.message === 'ER_BEIKE_EXCEL_NO_DATA') {
          error.message = '上传Excel文件无数据'
        } else if (error.message === 'ER_BEIKE_EXCEL_SANFANG_DATE_ERROR') {
          error.message = '上传Excel文件自营数据日期与三方日期不符'
        }
      }
    }
    message.error({
      content: error.message,
      duration: 5
    })
    return Promise.reject(error)
  }
)

export default service

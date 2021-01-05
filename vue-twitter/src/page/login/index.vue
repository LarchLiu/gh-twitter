<template>
    <div class="login">
        <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
            <a-form-item
            class="form-tem"
            required
            label="用户名"
            v-bind="validateInfos.userName">
            <!-- eslint-disable-next-line -->
                <a-input v-model:value="form.userName"  :placeholder="namePlacholder" />
            </a-form-item>

            <a-form-item
            class="form-tem"
            required
            label="密码"
            v-bind="validateInfos.password">
            <!-- eslint-disable-next-line -->
                <a-input v-model:value="form.password" :placeholder="passWordPlacholder" />
            </a-form-item>

            <a-form-item :wrapper-col="{ span: 14, offset: 4 }"
            >
                <a-button type="primary" @click="onSubmit">
                    确认
                </a-button>
                <a-button class="cancel" @click="resetFields">
                    取消
                </a-button>
            </a-form-item>
        </a-form>
    </div>
</template>
<script>
import { getCurrentInstance, reactive } from 'vue'
import { useForm } from '@ant-design-vue/use'
import { useStore } from 'vuex'

export default {
  setup () {
    // 路由实例
    const { ctx } = getCurrentInstance()
    const labelCol = { span: 4 }
    const wrapperCol = { span: 14 }
    const store = useStore()

    // form表单内容
    const form = reactive({
      userName: '',
      password: ''
    })
    // placholder
    const namePlacholder = 'Please enter your name'
    const passWordPlacholder = 'Please enter your passWord'

    // 校验规则
    const rulesRef = reactive({
      userName: [
        {
          required: true,
          message: namePlacholder
        }
      ],
      password: [
        {
          required: true,
          message: passWordPlacholder
        }
      ]
    })

    const { resetFields, validate, validateInfos } = useForm(form, rulesRef)

    // submit
    const onSubmit = (e) => {
      e.preventDefault()

      store.dispatch('setGHToken', form.password)

      validate()
        .then(() => {
          ctx.$router.push({
            path: '/'
          })
        })
    }

    return {
      labelCol,
      wrapperCol,
      form,
      namePlacholder,
      passWordPlacholder,
      onSubmit,
      validateInfos,
      resetFields,
      name
    }
  }
}
</script>
<style lang="less" scoped>
    .login{
        width: 500px;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
    .cancel{
        margin-left: 10px;
    }
    .form-tem{
        display: flex;
    }
</style>

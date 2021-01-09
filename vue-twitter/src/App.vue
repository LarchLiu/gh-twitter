<template>
    <router-view/>
</template>

<script>
import { onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'App',
  setup () {
    const store = useStore()

    const isMobile = () => {
      const flag = navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i)
      return !!flag
    }

    onMounted(() => {
      const mobile = isMobile()
      store.dispatch('setIsMobile', mobile)
    })

    return {
      isMobile
    }
  }
}
</script>

<template>
  <div class="img-view" @click="bigImg">
    <!-- 遮罩层 -->
    <div class="img-mask"></div>
    <div class="img-preview">
      <img :src="imgSrc" @click="imgClick">
    </div>
  </div>
</template>
<script>
export default {
  props: ['imgSrc'],
  data () {
    return {
      imgClicked: false
    }
  },
  methods: {
    bigImg () {
      // 发送事件
      if (this.imgClicked) {
        this.imgClicked = false
      } else {
        this.$emit('clickit')
      }
    },
    imgClick () {
      this.imgClicked = true
    }
  }
}
</script>
<style lang="less" scoped>
.img-view {
  width: 100%;
  height: 100%;

  .img-mask {
    position: fixed;
    z-index: 9999;
    top: 0;
    left: 0;
    background: rgba(0, 0, 0, 0.7);
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .img-preview {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    overflow: auto;
    outline: 0;
    width: 100%;
    height: 100%;
    z-index: 10000;

    img {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      max-width: 100%;
      max-height: 100%;
      margin: auto;
    }
  }
}
</style>

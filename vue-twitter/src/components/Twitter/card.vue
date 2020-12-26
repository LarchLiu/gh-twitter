<template>
  <div
    :class="[ 'card', mouseEnter ? 'enter' : '']"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="head">
      <a-avatar
        class="avatar"
        :size="50"
        :src="tweet.Avatar"
      />
      <div class="detail">
        <div id="name">
          <span style="font-weight: 700; border: 0 solid black; margin-right: 10px">{{ tweet.Name }}</span>
          <span v-if="tweet.IsRetweet">转推 </span>
          <span style="font-weight: 400; color: rgb(91, 112, 131);">
            {{ '@' + tweet.Username + ' · ' + getTime(tweet.Timestamp) }}
          </span>
        </div>
        <div class="text">
          <span v-html="tweet.HTML" />
        </div>
        <div
          v-if="tweet.Photos && tweet.Photos.length > 0"
          class="image"
        >
          <a
            v-for="(img, i) in tweet.Photos"
            :key="i"
            :class="checkImgRadiusClass(i, tweet.Photos.length)"
          >
            <img
              v-if="imgHeight"
              :src="img"
              :width="checkImgWidth(i, tweet.Photos.length)"
              :height="imgHeight"
              :class="checkImgClass(i, tweet.Photos.length)"
            >
            <img
              v-else
              :ref="i === 0 && tweet.Photos.length > 1 ? 'first-img' : 'other-img'"
              :src="img + '?'"
              :width="checkImgWidth(i, tweet.Photos.length)"
              :class="checkImgClass(i, tweet.Photos.length)"
              :onload="imgOnload"
            >
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, getCurrentInstance, watch } from 'vue'
import { formatTime } from '@/utils/index.js'

export default {
  name: 'Card',
  props: {
    tweet: {
      type: Object,
      default () {
        return {}
      }
    }
  },
  setup (props) {
    const { ctx } = getCurrentInstance()
    const mouseEnter = ref(false)
    const imgHeight = ref(0)

    const getTime = (timestamp) => {
      return formatTime(timestamp, null)
    }

    const checkImgRadiusClass = (idx, len) => {
      idx++
      const className = []
      if (len === 1) {
        className.push('lt-radius')
        className.push('rt-radius')
        className.push('lb-radius')
        className.push('rb-radius')
      } else {
        if (idx === 1) {
          className.push('lt-radius')
        }
        if (idx === 2) {
          className.push('rt-radius')
        }
        if (len % 2) {
          if (idx === len) {
            className.push('lb-radius')
            className.push('rb-radius')
          }
        } else {
          if (idx === len - 1) {
            className.push('lb-radius')
          }
          if (idx === len) {
            className.push('rb-radius')
          }
        }
      }
      return className
    }

    const checkImgClass = (idx, len) => {
      const className = []
      if (len > 1) {
        className.push('fit-contain')
        if (idx === 0) {
          className.push('first-img')
        }
      }
      idx++
      if (idx % 2) {
        className.push('mg-right')
      }
      if (len > 2) {
        className.push('mg-bottom')
      }
      return className
    }

    const checkImgWidth = (idx, len) => {
      if (len % 2 && idx === len - 1) {
        return 500
      } else {
        return (500 - (len > 1 ? 2 : 0)) / (len > 2 ? 2 : len)
      }
    }

    const imgOnload = () => {
      const e = ctx.$refs['first-img']
      if (e && e.width && e.naturalHeight && e.naturalWidth) {
        imgHeight.value = e.naturalHeight * e.width / e.naturalWidth
      }
    }

    const handleMouseEnter = () => {
      mouseEnter.value = true
    }
    const handleMouseLeave = () => {
      mouseEnter.value = false
    }

    onMounted(() => {
      // ctx.$nextTick(() => {
      //   if (props.tweet.Photos && props.tweet.Photos.length > 1) {
      //     const e = document.querySelector('.first-img')
      //     if (e) {
      //       e.addEventListener('load', () => {
      //         imgOnload(e)
      //       })
      //     }
      //   }
      // })
    })

    watch(() => props.tweet, () => {
      // 内容改变图片也改变重新计算图片应该显示的高度
      imgHeight.value = 0
    })

    return {
      getTime,
      mouseEnter,
      checkImgRadiusClass,
      checkImgClass,
      handleMouseEnter,
      handleMouseLeave,
      checkImgWidth,
      imgOnload,
      imgHeight
    }
  }
}
</script>

<style lang="less" scoped>
  .card {
    padding: 10px 10px;
    border-bottom: 1px solid #eee;
    overflow: hidden;
    text-align: left;
    font-size: 15px;

    .head {
      overflow: hidden;

      .avatar {
        float: left;
        width: 50px;
      }
      .detail {
        float: right;
        width: 500px;

        .image {
          margin-top: 10px;
          .lt-radius {
            img {
              border-top-left-radius: 10px;
            }
          }
          .lb-radius {
            img {
              border-bottom-left-radius: 10px;
            }
          }
          .rt-radius {
            img {
              border-top-right-radius: 10px;
            }
          }
          .rb-radius {
            img {
              border-bottom-right-radius: 10px;
            }
          }
          .mg-right {
            margin-right: 2px;
          }
          .mg-bottom {
            margin-bottom: 2px;
          }
          .fit-contain {
            object-fit: cover;
          }
        }
        .text {
          :deep(.Emoji) {
            width: 18px;
            height: 18px;
          }
          :deep(.PrettyLink-value) {
            color: rgb(27, 149, 224);
          }
          :deep(.PrettyLink-prefix) {
            color: rgb(27, 149, 224);
          }
        }
      }
    }
  }
  .enter {
    background: rgb(247, 247, 247);
  }
</style>

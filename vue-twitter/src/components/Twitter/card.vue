<template>
  <div
    :class="[ 'card', mouseEnter ? 'enter' : '']"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div v-if="isMobile" class="body-mobile">
      <div id="header">
        <a-avatar
          class="avatar"
          :size="46"
          :src="tweet.Avatar"
        />
        <span id="name">
          <span style="font-weight: 700; border: 0 solid black; margin-right: 10px">{{ tweet.Name }}</span>
          <span v-if="tweet.IsRetweet">转推 <span style="font-weight: 400; color: rgb(91, 112, 131);">{{ getRtName(tweet.Text) }}</span></span>
          <span v-else style="font-weight: 400; color: rgb(91, 112, 131);">
            {{ '@' + tweet.Username }}
          </span>
        </span>
        <div>
          <span style="font-weight: 400; color: rgb(91, 112, 131);">
            {{ getTime(tweet.Timestamp) }}
          </span>
        </div>
      </div>
      <div ref="detail" class="detail">
        <div class="text">
          <span v-html="tweet.HTML" />
        </div>
        <div v-if="tweet.Videos" class="video">
          <div v-for="(video, i) in tweet.Videos" :key="i" class="all-radius">
            <vue3-video-player :src="video.URL" :cover="video.Preview"></vue3-video-player>
          </div>
        </div>
        <div
          v-else-if="tweet.Photos && tweet.Photos.length > 0"
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
              @click="clickImg($event)"
            >
            <img
              v-else
              :ref="i === 0 && tweet.Photos.length > 1 ? 'firstImg' : (tweet.Photos.length === 1 ? 'onlyOneImg' : 'otherImg')"
              :src="img"
              :width="checkImgWidth(i, tweet.Photos.length)"
              :class="checkImgClass(i, tweet.Photos.length)"
              :onload="imgOnload"
              @click="clickImg($event)"
            >
          </a>
        </div>
      </div>
    </div>
    <div v-else class="body">
      <a-avatar
        class="avatar"
        :size="50"
        :src="tweet.Avatar"
      />
      <div ref="detail" class="detail">
        <div id="name">
          <span style="font-weight: 700; border: 0 solid black; margin-right: 10px">{{ tweet.Name }}</span>
          <span v-if="tweet.IsRetweet">转推 <span style="font-weight: 400; color: rgb(91, 112, 131);">{{ getRtName(tweet.Text) + ' · ' + getTime(tweet.Timestamp) }}</span></span>
          <span v-else style="font-weight: 400; color: rgb(91, 112, 131);">
            {{ '@' + tweet.Username + ' · ' + getTime(tweet.Timestamp) }}
          </span>
        </div>
        <div class="text">
          <span v-html="tweet.HTML" />
        </div>
        <div v-if="tweet.Videos" class="video">
          <div v-for="(video, i) in tweet.Videos" :key="i" class="all-radius">
            <div v-if="detailWidth"  :style="{ width: detailWidth }">
              <vue3-video-player :src="video.URL" :cover="video.Preview"></vue3-video-player>
            </div>
          </div>
        </div>
        <div
          v-else-if="tweet.Photos && tweet.Photos.length > 0"
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
              @click="clickImg($event)"
            >
            <img
              v-else
              :ref="i === 0 && tweet.Photos.length > 1 ? 'firstImg' : (tweet.Photos.length === 1 ? 'onlyOneImg' : 'otherImg')"
              :src="img"
              :width="checkImgWidth(i, tweet.Photos.length)"
              :class="checkImgClass(i, tweet.Photos.length)"
              :onload="imgOnload"
              @click="clickImg($event)"
            >
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { formatTime } from '@/utils/index.js'

export default {
  name: 'Card',
  props: {
    tweet: {
      type: Object,
      default () {
        return {}
      }
    },
    isMobile: {
      type: Boolean,
      default: false
    }
  },
  setup (props, context) {
    const mouseEnter = ref(false)
    const imgHeight = ref(0)
    const firstImg = ref(null)
    const onlyOneImg = ref(null)
    const detail = ref(null)
    const detailWidth = ref(0)
    const imgMaxHeight = ref(0)

    const getTime = (timestamp) => {
      return formatTime(timestamp, null)
    }

    const getRtName = (text) => {
      var regRt = /^RT (@.*?):/
      return regRt.test(text) ? regRt.exec(text)[1] : ''
    }

    const clickImg = (e) => {
      context.emit('imgClick', e.currentTarget.src)
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
      className.push('fit-contain')
      idx++
      if (idx % 2 && len > 1) {
        className.push('mg-right')
      }
      if (len > 2) {
        className.push('mg-bottom')
      }
      return className
    }

    const checkImgWidth = (idx, len) => {
      if (len % 2 && idx === len - 1) {
        return detailWidth.value
      } else {
        return (detailWidth.value - (len > 1 ? 2 : 0)) / (len > 2 ? 2 : len)
      }
    }

    const imgOnload = () => {
      const e = firstImg.value
      if (e && e.width && e.naturalHeight && e.naturalWidth) {
        imgHeight.value = e.naturalHeight * e.width / e.naturalWidth
      }
      const ooi = onlyOneImg.value
      if (ooi && ooi.width && ooi.naturalHeight && ooi.naturalWidth) {
        const height = ooi.naturalHeight * ooi.width / ooi.naturalWidth
        if (height > imgMaxHeight.value) {
          imgHeight.value = imgMaxHeight.value
        }
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

    // watch(() => props.tweet, () => {
    //   // 内容改变图片也改变重新计算图片应该显示的高度
    //   imgHeight.value = 0
    // })

    watch(detail, () => {
      if (detail.value) {
        detailWidth.value = detail.value.clientWidth
        imgMaxHeight.value = detailWidth.value * 9 / 16
      }
    })

    return {
      getTime,
      getRtName,
      clickImg,
      mouseEnter,
      checkImgRadiusClass,
      checkImgClass,
      handleMouseEnter,
      handleMouseLeave,
      checkImgWidth,
      imgOnload,
      imgHeight,
      firstImg,
      onlyOneImg,
      detail,
      detailWidth,
      imgMaxHeight
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

    .detail-info() {
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
      .video {
        .all-radius {
          border-radius: 10px;
          overflow: hidden;
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

    .body {
      overflow: hidden;

      .avatar {
        float: left;
        width: 50px;
      }
      .detail {
        float: right;
        width: 510px;

        .detail-info()
      }
    }

    .body-mobile {
      overflow: hidden;

      .avatar {
        width: 46px;
        float: left;
        margin-right: 10px;
      }
      .detail {
        margin-top: 10px;

        .detail-info()
      }
    }
  }
  .enter {
    background: rgb(247, 247, 247);
  }
</style>

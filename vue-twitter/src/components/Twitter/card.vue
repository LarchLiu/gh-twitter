<template>
  <div class="card">
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
          class="img"
        >
          <img
            v-for="(img, i) in tweet.Photos"
            :key="i"
            :src="img"
            :width="500/(tweet.Photos.length > 2 ? 2 : tweet.Photos.length)"
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { formatTime } from "/@/utils/index.js"

export default {
  name: 'Card',
  props: {
    tweet: {
      type: Object,
      default() {
        return {}
      }
    }
  },
  setup (props) {
    const getTime = (timestamp) => {
      return formatTime(timestamp, null)
    }

    return {
      getTime
    }
  }
}
</script>

<style lang="less" scoped>
  .card {
    padding: 10px 10px;
    border: 1px solid #eee;
    border-top: none;
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

        .img {
          margin-top: 10px;
          img {
            border-radius: 5px;
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
</style>


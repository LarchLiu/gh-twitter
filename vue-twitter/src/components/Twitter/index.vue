<template>
  <div class="twitter">
    <fixed-header :id-name="idName">
      <div
        v-if="Object.keys(detail).length > 0 && !isAll"
        class="header"
      >
        <span style="font-size: 15px; font-weight: 800; border: 0 solid black; margin-right: 20px">{{ detail.Profile.Name }}</span>
        <span style="font-weight: 700; color: rgb(15, 20, 25);">
          {{ detail.Profile.TweetsCount }}
        </span>
        <span style="color: rgb(91, 112, 131); margin-right: 20px"> 推文</span>
        <span style="color: rgb(91, 112, 131); margin-right: 20px;float:right">{{ "最后更新时间: " + getTime(detail.LastUpdateTime) }}</span>
      </div>
      <div
        v-else
        class="header"
      >
        <span style="font-size: 15px; font-weight: 800; border: 0 solid black; margin-right: 20px">推文</span>
        <span v-if="Object.keys(detail).length > 0 && isAll" style="color: rgb(91, 112, 131); margin-right: 20px;float:right">
          {{ "最后更新时间: " + getTime(detail.LastUpdateTime) }}
        </span>
      </div>
    </fixed-header>
    <div
      v-if="Object.keys(detail).length > 0"
      class="content"
    >
      <div v-if="!isAll" class="profile">
        <div class="name">
          <div style="font-size: 19px; font-weight: 800; border: 0 solid black;">
            {{ detail.Profile.Name }}
          </div>
          <div style="color: rgb(91, 112, 131) ; font-size: 15px; font-weight: 400; border: 0 solid black;">
            {{ '@' + detail.Profile.Username }}
          </div>
        </div>
        <div class="biography">
          {{ detail.Profile.Biography }}
        </div>
        <div class="follow">
          <span style="font-weight: 700; color: rgb(15, 20, 25);">
            {{ detail.Profile.FollowingCount }}
          </span>
          <span style="color: rgb(91, 112, 131); margin-right: 20px"> 正在关注</span>
          <span style="font-weight: 700; color: rgb(15, 20, 25);">
            {{ detail.Profile.FollowersCount }}
          </span>
          <span style="color: rgb(91, 112, 131); margin-right: 20px"> 关注者</span>
        </div>
        <div class="link-top" />
      </div>
      <div class="tweets">
        <div
          v-for="(tweet, i) in detail.Tweets"
          :key="i"
        >
          <card v-if="isAll && usersObj" :tweet="margeDetail(tweet, { Avator: usersObj[tweet.Username].Avator, Name: usersObj[tweet.Username].Name})" />
          <card v-else :tweet="margeDetail(tweet, detail.Profile)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import FixedHeader from '@/components/FixedHeader/index.vue'
import { parseTime } from '@/utils/index.js'
import Card from './card.vue'

export default {
  name: 'Twitter',
  components: { Card, FixedHeader },
  props: {
    idName: {
      type: String,
      default: 'twitter'
    },
    detail: {
      type: Object,
      default () {
        return {}
      }
    },
    isAll: {
      type: Boolean,
      default: false
    },
    usersObj: {
      type: Object,
      default () {
        return {}
      }
    }
  },
  methods: {
    margeDetail (tweet, profile) {
      return { Avatar: profile.Avatar, Name: profile.Name, ...tweet }
    },
    getTime (unix) {
      return parseTime(unix)
    }
  }
}
</script>

<style lang="less" scoped>
  .twitter {
    margin-bottom: 15px;
    color: rgba(0, 0, 0, 1.00);
    ::v-deep(.fixed-header) {
      position:fixed;
      top:0;
      width: 600px;
      z-index:999;
    }
    .header {
      padding: 10px;
      background-color: #f7f7f7;
      border: 1px solid #eee;
      h3 {
        max-width: 160px;
        height: 18px;
        overflow: hidden;
        margin: 0;
      }
    }
    .content {
      border: 1px solid #eee;
      border-top: none;
      overflow: hidden;
      text-align: left;
      font-size: 15px;

      .profile {
        .name {
          padding: 10px 10px;
        }
        .biography {
          padding: 0 10px;
        }
        .follow {
          padding: 10px 10px;
        }
        .link-top {
          width: 100%;
          height: 1px;
          border-top: solid #eee 1px;
        }

      }
    }
  }
</style>

<template>
  <div class="twitter">
    <fixed-header :id-name="idName">
      <div
        v-if="Object.keys(detail).length > 0"
        class="header"
      >
        <span style="font-size: 15px; font-weight: 800; border: 0 solid black; margin-right: 20px">{{ detail.Profile.Name }}</span>
        <span style="font-weight: 700; color: rgb(15, 20, 25);">
          {{ detail.Profile.TweetsCount }}
        </span>
        <span style="color: rgb(91, 112, 131); margin-right: 20px"> 推文</span>
      </div>
      <div
        v-else
        class="header"
      >
        <span style="font-size: 15px; font-weight: 800; border: 0 solid black; margin-right: 20px">推文</span>
      </div>
    </fixed-header>
    <div
      v-if="Object.keys(detail).length > 0"
      class="content"
    >
      <div class="profile">
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
          <card :tweet="margeDetail(tweet, detail.Profile)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import FixedHeader from '/@/components/FixedHeader/index.vue';
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
      default() {
        return {}
      }
    }
  },
  methods: {
    margeDetail (tweet, profile) {
      return { Avatar: profile.Avatar, Name: profile.Name, ...tweet }
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
      padding: 10px 10px;
      border: 1px solid #eee;
      border-top: none;
      overflow: hidden;
      text-align: left;
      font-size: 15px;

      .profile {
        .name {
          margin-bottom: 10px;
        }
        .biography {
          margin-bottom: 10px;
        }
        .follow {
          margin-bottom: 10px;
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

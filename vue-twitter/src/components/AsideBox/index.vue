<template>
  <div class="aside">
    <fixed-header v-if="needFixed" :id-name="idName" style-class-name="fixed-header-aside" @change="handleFixedChange">
      <div class="aside-header" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
        <span style="font-size: 15px; font-weight: 800; border: 0 solid black;">{{ title }}</span>
        <slot name="btn" />
      </div>
    </fixed-header>
    <div v-else class="aside-header">
      <span style="font-size: 15px; font-weight: 800; border: 0 solid black;">{{ title }}</span>
      <slot name="btn" />
    </div>
    <div
      class="aside-content"
      @mouseenter="handleMouseEnter"
      @mouseleave="handleMouseLeave"
    >
      <slot />
    </div>
    <div
      v-if="pop_display && fixed"
      class="aside-content pop-display"
      @mouseenter="handleMouseEnter"
      @mouseleave="handleMouseLeave"
    >
      <slot />
    </div>
  </div>
</template>

<script>
import FixedHeader from '@/components/FixedHeader/index.vue'

export default {
  name: 'AsideBox',
  components: { FixedHeader },
  props: {
    title: {
      type: String,
      default: ''
    },
    needFixed: {
      type: Boolean,
      default: false
    },
    idName: {
      type: String,
      default: 'aside'
    }
  },
  data () {
    return {
      fixed: false,
      pop_display: true
    }
  },
  methods: {
    handleFixedChange (val) {
      this.fixed = val
      this.pop_display = true
    },
    handleMouseEnter () {
      if (this.fixed) {
        this.pop_display = true
      }
    },
    handleMouseLeave () {
      if (this.fixed) {
        this.pop_display = true
      }
    }
  }
}
</script>

<style lang="less" scoped>
  .aside {
    margin-bottom: 15px;
    ::v-deep(.fixed-header-aside) {
      position:fixed;
      top:0;
      width: 210px;
      z-index:999;
    }
    .aside-header {
      padding: 10px;
      background-color: #f7f7f7;
      border: 1px solid #eee;
      h3 {
        max-width: 160px;
        height: 18px;
        overflow: hidden;
        margin: 0;
      }
      i {
        float: right;
        margin-top: -14px;
      }
    }
    .aside-content {
      padding: 10px 10px;
      border: 1px solid #eee;
      border-top: none;
      overflow: hidden;
      text-align: left;
      &.pop-display {
        position:fixed;
        top:45px;
        width: 210px;
        z-index:999;
        background-color: white;
      }
    }
  }
</style>
